package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/facette/facette/pkg/config"
	"github.com/facette/facette/pkg/connector"
	"github.com/facette/facette/pkg/library"
	"github.com/facette/facette/pkg/server"
	"github.com/facette/facette/pkg/utils"
)

var (
	serverConfig *config.Config
)

func Test_CatalogOriginList(test *testing.T) {
	base := []string{
		"test1",
		"test2",
	}

	result := make([]string, 0)

	// Test GET on source list
	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/origins/", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on source list (offset and limit)
	result = make([]string, 0)

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/origins/?limit=1", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[:1], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[:1], result)
		test.Fail()
	}

	result = make([]string, 0)

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/origins/?offset=1&limit=1",
		serverConfig.BindAddr), nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[1:2], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[1:2], result)
		test.Fail()
	}
}

func Test_CatalogOriginGet(test *testing.T) {
	base := &server.SourceResponse{Name: "source1", Origins: []string{"test1", "test2"}}
	result := &server.SourceResponse{}

	// Test GET on source1 item
	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/source1", serverConfig.BindAddr),
		nil, false, &result)
	result.Updated = ""

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on source2 item (with filter settings)
	base = &server.SourceResponse{Name: "source2", Origins: []string{"test1"}}
	result = &server.SourceResponse{}

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/source2", serverConfig.BindAddr),
		nil, false, &result)
	result.Updated = ""

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on unknown item
	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/unknown", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}
}

func Test_CatalogSourceList(test *testing.T) {
	base := []string{
		"source1",
		"source2",
	}

	result := make([]string, 0)

	// Test GET on source list
	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/", serverConfig.BindAddr), nil,
		false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on source list (limit)
	result = make([]string, 0)

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/?limit=1", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[:1], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[:1], result)
		test.Fail()
	}

	// Test GET on source list (offset and limit)
	result = make([]string, 0)

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/?offset=1&limit=1",
		serverConfig.BindAddr), nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[1:2], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[1:2], result)
		test.Fail()
	}
}

func Test_CatalogSourceGet(test *testing.T) {
	base := &server.SourceResponse{Name: "source1", Origins: []string{"test1", "test2"}}
	result := &server.SourceResponse{}

	// Test GET on source1 item
	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/source1", serverConfig.BindAddr),
		nil, false, &result)
	result.Updated = ""

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on source2 item
	base = &server.SourceResponse{Name: "source2", Origins: []string{"test1"}}
	result = &server.SourceResponse{}

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/source2", serverConfig.BindAddr),
		nil, false, &result)
	result.Updated = ""

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on unknown item
	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/sources/unknown", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}
}

func Test_CatalogMetricList(test *testing.T) {
	// Test GET on metrics list
	base := []string{
		"database1.test",
		"database1/test",
		"database2.test",
		"database2/test",
		"database3/test",
	}

	result := make([]string, 0)

	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/", serverConfig.BindAddr), nil,
		false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on metrics list (offset and limit)
	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/?limit=2", serverConfig.BindAddr),
		nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[:2], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[:2], result)
		test.Fail()
	}

	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/?offset=2&limit=2",
		serverConfig.BindAddr), nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[2:4], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[2:4], result)
		test.Fail()
	}

	// Test GET on metrics list (source-specific)
	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/?source=source1",
		serverConfig.BindAddr), nil, false, &result)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base[:4], result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base[:4], result)
		test.Fail()
	}
}

func Test_CatalogMetricGet(test *testing.T) {
	base := &server.MetricResponse{Name: "database2/test", Sources: []string{"source1", "source2"},
		Origins: []string{"test1"}}
	result := &server.MetricResponse{}

	// Test GET on metric item
	response := execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/database2/test",
		serverConfig.BindAddr), nil, false, &result)
	result.Updated = ""

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(base, result) {
		test.Logf("\nExpected %#v\nbut got  %#v", base, result)
		test.Fail()
	}

	// Test GET on unknown metric item
	response = execTestRequest(test, "GET", fmt.Sprintf("http://%s/catalog/metrics/unknown/test",
		serverConfig.BindAddr), nil, false, &result)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}
}

func Test_LibrarySourceGroupHandle(test *testing.T) {
	// Define a sample source group
	group := &library.Group{Item: library.Item{Name: "group1", Description: "A great group description."}}
	group.Entries = append(group.Entries, &library.GroupEntry{Pattern: "glob:source*", Origin: "test1"})

	expandData := server.ExpandRequest{[3]string{"test1", "group:group1-updated", "database1/test"}}

	expandBase := server.ExpandRequest{}
	expandBase = append(expandBase, [3]string{"test1", "source1", "database1/test"})
	expandBase = append(expandBase, [3]string{"test1", "source2", "database1/test"})

	execGroupHandle(test, "sourcegroups", group, expandData, expandBase)
}

func Test_LibraryMetricGroupHandle(test *testing.T) {
	// Define a sample metric group
	group := &library.Group{Item: library.Item{Name: "group1", Description: "A great group description."}}
	group.Entries = append(group.Entries, &library.GroupEntry{Pattern: "database1/test", Origin: "test1"})
	group.Entries = append(group.Entries, &library.GroupEntry{Pattern: "regexp:database[23]/test", Origin: "test1"})

	expandData := server.ExpandRequest{[3]string{"test1", "source1", "group:group1-updated"}}

	expandBase := server.ExpandRequest{}
	expandBase = append(expandBase, [3]string{"test1", "source1", "database1/test"})
	expandBase = append(expandBase, [3]string{"test1", "source1", "database2/test"})
	expandBase = append(expandBase, [3]string{"test1", "source1", "database3/test"})

	execGroupHandle(test, "metricgroups", group, expandData, expandBase)
}

func Test_LibraryGraphHandle(test *testing.T) {
	baseURL := fmt.Sprintf("http://%s/library/graphs/", serverConfig.BindAddr)

	// Define a sample graph
	stack := &library.Stack{Name: "stack0"}

	group := &library.OperGroup{Name: "group0", Type: connector.OperGroupTypeAvg}
	group.Series = append(group.Series, &library.Serie{Name: "serie0", Origin: "test", Source: "source1",
		Metric: "database1/test"})
	group.Series = append(group.Series, &library.Serie{Name: "serie1", Origin: "test", Source: "source2",
		Metric: "group:group1"})

	stack.Groups = append(stack.Groups, group)

	group = &library.OperGroup{Name: "serie2", Type: connector.OperGroupTypeNone}
	group.Series = append(group.Series, &library.Serie{Name: "serie2", Origin: "test", Source: "group:group1",
		Metric: "database2/test"})

	stack.Groups = append(stack.Groups, group)

	graphBase := &library.Graph{Item: library.Item{Name: "graph1", Description: "A great graph description."},
		StackMode: library.StackModeNormal}
	graphBase.Stacks = append(graphBase.Stacks, stack)

	// Test GET on graphs list
	listBase := server.ItemListResponse{}
	listResult := server.ItemListResponse{}

	response := execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test GET on a unknown graph item
	response = execTestRequest(test, "GET", baseURL+"/00000000-0000-0000-0000-000000000000", nil, false, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test POST into graph
	data, _ := json.Marshal(graphBase)

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusCreated {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
		test.Fail()
	}

	if response.Header.Get("Location") == "" {
		test.Logf("\nExpected `Location' header")
		test.Fail()
	}

	graphBase.ID = response.Header.Get("Location")[strings.LastIndex(response.Header.Get("Location"), "/")+1:]

	// Test GET on graph item
	graphResult := &library.Graph{}

	response = execTestRequest(test, "GET", baseURL+graphBase.ID, nil, false, &graphResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(graphBase, graphResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", graphBase, graphResult)
		test.Fail()
	}

	// Test GET on graphs list
	listBase = server.ItemListResponse{&server.ItemResponse{
		ID:          graphBase.ID,
		Name:        graphBase.Name,
		Description: graphBase.Description,
	}}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test PUT on graph item
	graphBase.Name = "graph1-updated"

	data, _ = json.Marshal(graphBase)

	response = execTestRequest(test, "PUT", baseURL+graphBase.ID, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "PUT", baseURL+graphBase.ID, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	// Test GET on graph item
	graphResult = &library.Graph{}

	response = execTestRequest(test, "GET", baseURL+graphBase.ID, nil, false, &graphResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(graphBase, graphResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", graphBase, graphResult)
		test.Fail()
	}

	// Test DELETE on graph item
	response = execTestRequest(test, "DELETE", baseURL+graphBase.ID, nil, false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+graphBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+graphBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test volatile POST into graph
	graphBase.ID = ""
	data, _ = json.Marshal(graphBase)

	response = execTestRequest(test, "POST", baseURL+"?volatile=1", strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "POST", baseURL+"?volatile=1", strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusCreated {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
		test.Fail()
	}

	if response.Header.Get("Location") == "" {
		test.Logf("\nExpected `Location' header")
		test.Fail()
	}

	graphBase.ID = response.Header.Get("Location")[strings.LastIndex(response.Header.Get("Location"), "/")+1:]

	// Test GET on volatile graph item
	graphResult = &library.Graph{}

	response = execTestRequest(test, "GET", baseURL+graphBase.ID, nil, false, &graphResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(graphBase, graphResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", graphBase, graphResult)
		test.Fail()
	}

	// Test GET on volatile graph item
	graphResult = &library.Graph{}

	response = execTestRequest(test, "GET", baseURL+graphBase.ID, nil, false, &graphResult)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test GET on graphs list (offset and limit)
	listBase = server.ItemListResponse{}

	for i := 0; i < 3; i += 1 {
		graphTemp := &library.Graph{}
		utils.Clone(graphBase, graphTemp)

		graphTemp.ID = ""
		graphTemp.Name = fmt.Sprintf("graph1-%d", i)

		data, _ = json.Marshal(graphTemp)

		response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

		if response.StatusCode != http.StatusCreated {
			test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
			test.Fail()
		}

		location := response.Header.Get("Location")

		if location == "" {
			test.Logf("\nExpected `Location' header")
			test.Fail()
		}

		graphTemp.ID = location[strings.LastIndex(location, "/")+1:]

		listBase = append(listBase, &server.ItemResponse{
			ID:          graphTemp.ID,
			Name:        graphTemp.Name,
			Description: graphTemp.Description,
		})
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?limit=1", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[:1], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[:1], listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?offset=1&limit=2", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[1:3], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[1:3], listResult)
		test.Fail()
	}
}

func Test_LibraryCollectionHandle(test *testing.T) {
	var collectionBase struct {
		*library.Collection
		Parent string `json:"parent"`
	}

	baseURL := fmt.Sprintf("http://%s/library/collections/", serverConfig.BindAddr)

	// Define a sample collection
	collectionBase.Collection = &library.Collection{Item: library.Item{Name: "collection0",
		Description: "A great collection description."}}

	collectionBase.Entries = append(collectionBase.Entries,
		&library.CollectionEntry{ID: "00000000-0000-0000-0000-000000000000",
			Options: map[string]string{"range": "-1h"}})
	collectionBase.Entries = append(collectionBase.Entries,
		&library.CollectionEntry{ID: "00000000-0000-0000-0000-000000000000",
			Options: map[string]string{"range": "-1d"}})
	collectionBase.Entries = append(collectionBase.Entries,
		&library.CollectionEntry{ID: "00000000-0000-0000-0000-000000000000",
			Options: map[string]string{"range": "-1w"}})

	// Test GET on collections list
	listBase := server.ItemListResponse{}
	listResult := server.ItemListResponse{}

	response := execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test GET on a unknown collection item
	response = execTestRequest(test, "GET", baseURL+"/00000000-0000-0000-0000-000000000000", nil, false, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test POST into collection
	data, _ := json.Marshal(collectionBase)

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusCreated {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
		test.Fail()
	}

	if response.Header.Get("Location") == "" {
		test.Logf("\nExpected `Location' header")
		test.Fail()
	}

	collectionBase.ID = response.Header.Get("Location")[strings.LastIndex(response.Header.Get("Location"), "/")+1:]

	// Test GET on collection item
	collectionResult := &library.Collection{}

	response = execTestRequest(test, "GET", baseURL+collectionBase.ID, nil, false, &collectionResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(collectionBase.Collection, collectionResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", collectionBase.Collection, collectionResult)
		test.Fail()
	}

	// Test GET on collections list
	listBase = server.ItemListResponse{&server.ItemResponse{
		ID:          collectionBase.ID,
		Name:        collectionBase.Name,
		Description: collectionBase.Description,
	}}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test PUT on collection item
	collectionBase.Name = "collection1-updated"

	data, _ = json.Marshal(collectionBase.Collection)

	response = execTestRequest(test, "PUT", baseURL+collectionBase.ID, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "PUT", baseURL+collectionBase.ID, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	// Test GET on collection item
	collectionResult = &library.Collection{}

	response = execTestRequest(test, "GET", baseURL+collectionBase.ID, nil, false, &collectionResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(collectionBase.Collection, collectionResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", collectionBase, collectionResult)
		test.Fail()
	}

	// Test DELETE on collection item
	response = execTestRequest(test, "DELETE", baseURL+collectionBase.ID, nil, false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+collectionBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+collectionBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test GET on collections list (offset and limit)
	listBase = server.ItemListResponse{}

	for i := 0; i < 3; i += 1 {
		collectionTemp := &library.Collection{}
		utils.Clone(collectionBase, collectionTemp)

		collectionTemp.ID = ""
		collectionTemp.Name = fmt.Sprintf("collection1-%d", i)

		data, _ = json.Marshal(collectionTemp)

		response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

		if response.StatusCode != http.StatusCreated {
			test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
			test.Fail()
		}

		location := response.Header.Get("Location")

		if location == "" {
			test.Logf("\nExpected `Location' header")
			test.Fail()
		}

		collectionTemp.ID = location[strings.LastIndex(location, "/")+1:]

		listBase = append(listBase, &server.ItemResponse{
			ID:          collectionTemp.ID,
			Name:        collectionTemp.Name,
			Description: collectionTemp.Description,
		})
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?limit=1", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[:1], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[:1], listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?offset=1&limit=2", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[1:3], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[1:3], listResult)
		test.Fail()
	}
}

func execGroupHandle(test *testing.T, urlPrefix string, groupBase *library.Group, expandData,
	expandBase server.ExpandRequest) {

	baseURL := fmt.Sprintf("http://%s/library/%s/", serverConfig.BindAddr, urlPrefix)

	// Test GET on groups list
	listBase := server.ItemListResponse{}
	listResult := server.ItemListResponse{}

	response := execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test GET on a unknown group item
	response = execTestRequest(test, "GET", baseURL+"/00000000-0000-0000-0000-000000000000", nil, false, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test POST into group
	data, _ := json.Marshal(groupBase)

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusCreated {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
		test.Fail()
	}

	if response.Header.Get("Location") == "" {
		test.Logf("\nExpected `Location' header")
		test.Fail()
	}

	groupBase.ID = response.Header.Get("Location")[strings.LastIndex(response.Header.Get("Location"), "/")+1:]

	// Test GET on group item
	groupResult := &library.Group{}

	response = execTestRequest(test, "GET", baseURL+groupBase.ID, nil, false, &groupResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(groupBase, groupResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", groupBase, groupResult)
		test.Fail()
	}

	// Test GET on groups list
	listBase = server.ItemListResponse{&server.ItemResponse{
		ID:          groupBase.ID,
		Name:        groupBase.Name,
		Description: groupBase.Description,
	}}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	// Test PUT on group item
	groupBase.Name = "group1-updated"

	data, _ = json.Marshal(groupBase)

	response = execTestRequest(test, "PUT", baseURL+groupBase.ID, strings.NewReader(string(data)), false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "PUT", baseURL+groupBase.ID, strings.NewReader(string(data)), true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	// Test GET on group item
	groupResult = &library.Group{}

	response = execTestRequest(test, "GET", baseURL+groupBase.ID, nil, false, &groupResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if !reflect.DeepEqual(groupBase, groupResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", groupBase, groupResult)
		test.Fail()
	}

	// Test group expansion
	data, _ = json.Marshal(expandData)

	expandResult := make([]server.ExpandRequest, 0)

	response = execTestRequest(test, "POST", fmt.Sprintf("http://%s/library/expand", serverConfig.BindAddr),
		strings.NewReader(string(data)), false, &expandResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	if len(expandResult) == 0 {
		test.Logf("\nExpected %#v\nbut got  %#v", expandBase, expandResult)
		test.Fail()
	} else if !reflect.DeepEqual(expandBase, expandResult[0]) {
		test.Logf("\nExpected %#v\nbut got  %#v", expandBase, expandResult[0])
		test.Fail()
	}

	// Test DELETE on group item
	response = execTestRequest(test, "DELETE", baseURL+groupBase.ID, nil, false, nil)

	if response.StatusCode != http.StatusUnauthorized {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusUnauthorized, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+groupBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	response = execTestRequest(test, "DELETE", baseURL+groupBase.ID, nil, true, nil)

	if response.StatusCode != http.StatusNotFound {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusNotFound, response.StatusCode)
		test.Fail()
	}

	// Test GET on groups list (offset and limit)
	listBase = server.ItemListResponse{}

	for i := 0; i < 3; i += 1 {
		groupTemp := &library.Group{}
		utils.Clone(groupBase, groupTemp)

		groupTemp.ID = ""
		groupTemp.Name = fmt.Sprintf("group1-%d", i)

		data, _ = json.Marshal(groupTemp)

		response = execTestRequest(test, "POST", baseURL, strings.NewReader(string(data)), true, nil)

		if response.StatusCode != http.StatusCreated {
			test.Logf("\nExpected %d\nbut got  %d", http.StatusCreated, response.StatusCode)
			test.Fail()
		}

		location := response.Header.Get("Location")

		if location == "" {
			test.Logf("\nExpected `Location' header")
			test.Fail()
		}

		groupTemp.ID = location[strings.LastIndex(location, "/")+1:]

		listBase = append(listBase, &server.ItemResponse{
			ID:          groupTemp.ID,
			Name:        groupTemp.Name,
			Description: groupTemp.Description,
		})
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL, nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase, listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase, listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?limit=1", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[:1], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[:1], listResult)
		test.Fail()
	}

	listResult = server.ItemListResponse{}

	response = execTestRequest(test, "GET", baseURL+"?offset=1&limit=2", nil, false, &listResult)

	if response.StatusCode != http.StatusOK {
		test.Logf("\nExpected %d\nbut got  %d", http.StatusOK, response.StatusCode)
		test.Fail()
	}

	for _, listItem := range listResult {
		listItem.Modified = ""
	}

	if !reflect.DeepEqual(listBase[1:3], listResult) {
		test.Logf("\nExpected %#v\nbut got  %#v", listBase[1:3], listResult)
		test.Fail()
	}
}

func execTestRequest(test *testing.T, method, url string, data io.Reader, auth bool,
	result interface{}) *http.Response {

	request, err := http.NewRequest(method, url, data)
	if err != nil {
		test.Fatal(err.Error())
	}

	if auth {
		// Add authentication (login: unittest, password: unittest)
		request.Header.Add("Authorization", "Basic dW5pdHRlc3Q6dW5pdHRlc3Q=")
	}

	if data != nil {
		request.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		test.Fatal(err.Error())
	}

	defer response.Body.Close()

	if result != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			test.Fatal(err.Error())
		}

		json.Unmarshal(body, result)
	}

	return response
}

func init() {
	// Load server configuration
	serverConfig = &config.Config{}
	if err := serverConfig.Load(flagConfig); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}
