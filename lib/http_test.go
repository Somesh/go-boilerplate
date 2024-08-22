package lib

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestRequestBuilder(t *testing.T) {
	ctx := context.Background()
	request := RequestBuilder(ctx, "test", nil)
	if request == nil {
		t.Errorf("Got Nil")
	}
}

func TestAddQueryParam(t *testing.T) {
	request := Request{}
	request.AddQueryParam("test_key", "test_value")
	if val, ok := request.QueryParams["test_key"]; !ok || val[0] != "test_value" {
		t.Errorf("Expected test_value , Got %s", val)
	}
}

func TestAddQueryParams(t *testing.T) {
	request := Request{}
	queryParams := map[string][]string{"test_key": []string{"1", "2"}}
	request.AddQueryParams(queryParams)
	if val, ok := request.QueryParams["test_key"]; !ok || val[0] != "1" || val[1] != "2" {
		t.Errorf("Expected [1,2] , Got %s", val)
	}
}

func TestCreateQueryString(t *testing.T) {
	request := Request{}
	queryParams := map[string][]string{"test_key": []string{"1", "2"},
		"test_key_1": []string{"3"}}
	request.AddQueryParams(queryParams)
	expected1 := "test_key=1&test_key=2&test_key_1=3"
	expected2 := "test_key_1=3&test_key=1&test_key=2"
	result := request.CreateQueryString()
	//Comparing with 2 expected values because map can be iterated in random order
	if result != expected1 && result != expected2 {
		t.Errorf("Got %s", result)
	}
}

func TestAddHeader(t *testing.T) {
	request := Request{}
	request.AddHeader("test_key", "test_value")
	if val, ok := request.Headers["test_key"]; !ok || val != "test_value" {
		t.Errorf("Expected test_value , Got %s", val)
	}
}

func TestAddheaders(t *testing.T) {
	request := Request{}
	queryParams := map[string]string{"test_key": "1", "test_key_1": "2"}
	request.AddHeaders(queryParams)
	if val, ok := request.Headers["test_key"]; !ok || val != "1" {
		t.Errorf("Expected 1 , Got %s", val)
	}
	if val, ok := request.Headers["test_key_1"]; !ok || val != "2" {
		t.Errorf("Expected 2 , Got %s", val)
	}
}

func TestAddPostBody(t *testing.T) {
	request := Request{}
	type PostBody struct {
		Name string
	}
	body := PostBody{Name: "Test"}
	request.AddPostBody(body)
	if request.RequestPostBody != body {
		t.Errorf("Got %+v", request.RequestPostBody)
	}
}

func TestExecuteRequest(t *testing.T) {
	url := "http://127.0.0.1:9011/test"
	httpClient := http.Client{}
	client := Client{Client: &httpClient,
		Timeout: 1 * time.Second}
	request := Request{RequestClient: &client,
		Url: url}
	httpmock.ActivateNonDefault(request.RequestClient.Client)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, `{"status": true}`))

	type ResponseType struct {
		Status bool `json:"status"`
	}
	var dest ResponseType
	_, err := request.executeRequest(context.Background(), "GET", false, &dest)
	if err != nil {
		t.Error(err)
	}
	if !dest.Status {
		t.Errorf("Expected Status True , Got %+v", dest)
	}
}
