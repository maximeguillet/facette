package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/facette/facette/pkg/library"
	"github.com/facette/facette/pkg/utils"
)

func (server *Server) handleGroup(writer http.ResponseWriter, request *http.Request) {
	var (
		groupID   string
		groupType int
	)

	if strings.HasPrefix(request.URL.Path, urlLibraryPath+"sourcegroups") {
		groupID = strings.TrimPrefix(request.URL.Path, urlLibraryPath+"sourcegroups/")
		groupType = library.LibraryItemSourceGroup
	} else if strings.HasPrefix(request.URL.Path, urlLibraryPath+"metricgroups") {
		groupID = strings.TrimPrefix(request.URL.Path, urlLibraryPath+"metricgroups/")
		groupType = library.LibraryItemMetricGroup
	}

	switch request.Method {
	case "DELETE":
		if groupID == "" {
			server.handleResponse(writer, serverResponse{mesgMethodNotAllowed}, http.StatusMethodNotAllowed)
			return
		} else if !server.handleAuth(writer, request) {
			server.handleResponse(writer, serverResponse{mesgAuthenticationRequired}, http.StatusUnauthorized)
			return
		}

		err := server.Library.DeleteItem(groupID, groupType)
		if os.IsNotExist(err) {
			server.handleResponse(writer, serverResponse{mesgResourceNotFound}, http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("ERROR: " + err.Error())
			server.handleResponse(writer, serverResponse{mesgUnhandledError}, http.StatusInternalServerError)
			return
		}

		server.handleResponse(writer, nil, http.StatusOK)

	case "GET", "HEAD":
		if groupID == "" {
			server.handleGroupList(writer, request)
			return
		}

		item, err := server.Library.GetItem(groupID, groupType)
		if os.IsNotExist(err) {
			server.handleResponse(writer, serverResponse{mesgResourceNotFound}, http.StatusNotFound)
			return
		} else if err != nil {
			log.Println("ERROR: " + err.Error())
			server.handleResponse(writer, serverResponse{mesgUnhandledError}, http.StatusInternalServerError)
			return
		}

		server.handleResponse(writer, item, http.StatusOK)

	case "POST", "PUT":
		var group *library.Group

		if response, status := server.parseStoreRequest(writer, request, groupID); status != http.StatusOK {
			server.handleResponse(writer, response, status)
			return
		}

		if request.Method == "POST" && request.FormValue("inherit") != "" {
			// Get group from library
			item, err := server.Library.GetItem(request.FormValue("inherit"), groupType)
			if os.IsNotExist(err) {
				server.handleResponse(writer, serverResponse{mesgResourceNotFound}, http.StatusNotFound)
				return
			} else if err != nil {
				log.Println("ERROR: " + err.Error())
				server.handleResponse(writer, serverResponse{mesgUnhandledError}, http.StatusInternalServerError)
				return
			}

			group = &library.Group{}
			*group = *item.(*library.Group)

			group.ID = ""
		} else {
			// Create a new group instance
			group = &library.Group{Item: library.Item{ID: groupID}, Type: groupType}
		}

		group.Modified = time.Now()

		// Parse input JSON for group data
		body, _ := ioutil.ReadAll(request.Body)

		if err := json.Unmarshal(body, group); err != nil {
			log.Println("ERROR: " + err.Error())
			server.handleResponse(writer, serverResponse{mesgResourceInvalid}, http.StatusBadRequest)
			return
		}

		// Store group data
		err := server.Library.StoreItem(group, groupType)
		if response, status := server.parseError(writer, request, err); status != http.StatusOK {
			log.Println("ERROR: " + err.Error())
			server.handleResponse(writer, response, status)
			return
		}

		if request.Method == "POST" {
			writer.Header().Add("Location", strings.TrimRight(request.URL.Path, "/")+"/"+group.ID)
			server.handleResponse(writer, nil, http.StatusCreated)
		} else {
			server.handleResponse(writer, nil, http.StatusOK)
		}

	default:
		server.handleResponse(writer, serverResponse{mesgMethodNotAllowed}, http.StatusMethodNotAllowed)
	}
}

func (server *Server) handleGroupList(writer http.ResponseWriter, request *http.Request) {
	var offset, limit int

	if response, status := server.parseListRequest(writer, request, &offset, &limit); status != http.StatusOK {
		server.handleResponse(writer, response, status)
		return
	}

	// Fill groups list
	items := make(ItemListResponse, 0)

	isSource := strings.HasPrefix(request.URL.Path, urlLibraryPath+"sourcegroups/")

	for _, group := range server.Library.Groups {
		if isSource && group.Type != library.LibraryItemSourceGroup ||
			!isSource && group.Type != library.LibraryItemMetricGroup {
			continue
		}

		if request.FormValue("filter") != "" && !utils.FilterMatch(request.FormValue("filter"), group.Name) {
			continue
		}

		items = append(items, &ItemResponse{
			ID:          group.ID,
			Name:        group.Name,
			Description: group.Description,
			Modified:    group.Modified.Format(time.RFC3339),
		})
	}

	response := &listResponse{
		list:   items,
		offset: offset,
		limit:  limit,
	}

	server.applyResponseLimit(writer, request, response)

	server.handleResponse(writer, response.list, http.StatusOK)
}

func (server *Server) handleGroupExpand(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		server.handleResponse(writer, serverResponse{mesgMethodNotAllowed}, http.StatusMethodNotAllowed)
		return
	}

	body, _ := ioutil.ReadAll(request.Body)

	query := ExpandRequest{}
	if err := json.Unmarshal(body, &query); err != nil {
		log.Println("ERROR: " + err.Error())
		server.handleResponse(writer, serverResponse{mesgResourceInvalid}, http.StatusBadRequest)
		return
	}

	response := make([]ExpandRequest, 0)

	for _, entry := range query {
		item := ExpandRequest{}

		if strings.HasPrefix(entry[1], library.LibraryGroupPrefix) {
			for _, sourceName := range server.Library.ExpandGroup(
				strings.TrimPrefix(entry[1], library.LibraryGroupPrefix),
				library.LibraryItemSourceGroup,
			) {
				if strings.HasPrefix(entry[2], library.LibraryGroupPrefix) {
					for _, metricName := range server.Library.ExpandGroup(
						strings.TrimPrefix(entry[2], library.LibraryGroupPrefix),
						library.LibraryItemMetricGroup,
					) {
						item = append(item, [3]string{entry[0], sourceName, metricName})
					}
				} else {
					item = append(item, [3]string{entry[0], sourceName, entry[2]})
				}
			}
		} else if strings.HasPrefix(entry[2], library.LibraryGroupPrefix) {
			for _, metricName := range server.Library.ExpandGroup(
				strings.TrimPrefix(entry[2], library.LibraryGroupPrefix),
				library.LibraryItemMetricGroup,
			) {
				item = append(item, [3]string{entry[0], entry[1], metricName})
			}
		} else {
			item = append(item, entry)
		}

		sort.Sort(item)
		response = append(response, item)
	}

	server.handleResponse(writer, response, http.StatusOK)
}
