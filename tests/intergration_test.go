package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/mustthink/go-storage-like-redis/config"
	"github.com/mustthink/go-storage-like-redis/internal"
	"github.com/mustthink/go-storage-like-redis/internal/handlers"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

type (
	TestClient struct {
		client *http.Client
		url    string
	}

	TestRequest struct {
		Type       string                            `json:"type"`
		Collection string                            `json:"collection"`
		Keys       []string                          `json:"keys"`
		Objects    map[string]object.RequestSettings `json:"objects"`
	}
)

func init() {
	path := fmt.Sprintf("../%s", config.DefaultConfig)
	app := internal.NewApplication(path)
	// running test server
	go app.Run()
}

func (c TestClient) doRequest(t *testing.T, requestMethod string, request TestRequest) handlers.DataCode {
	var response handlers.DataCode
	switch request.Type {
	case handlers.TypeCollection:
		response = &handlers.Response{}
	case handlers.TypeObject:
		response = &handlers.Responses{}
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		t.Errorf("couldn't marshal w err: %s", err.Error())
		return response
	}

	req, err := http.NewRequest(requestMethod, c.url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Errorf("couldn't create new request w err: %s", err.Error())
		return response
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		t.Errorf("couldn't do request w err: %s", err.Error())
		return response
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("couldn't read response.Body w err: %s", err.Error())
		return response
	}

	if err := json.Unmarshal(body, response); err != nil {
		t.Errorf("couldn't unmarshal response w err: %s", err.Error())
	}

	fmt.Println(string(body), " ", response)
	return response
}

func Test_BasicRequests(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "http://localhost:8081/"

	testClient := TestClient{
		client: client,
		url:    url,
	}

	tests := []struct {
		name          string
		request       TestRequest
		requestMethod string
		wantErr       bool
	}{
		{
			name: "POST collection",
			request: TestRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodPost,
		},
		{
			name: "GET collection",
			request: TestRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodGet,
		},
		{
			name: "POST object",
			request: TestRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Objects: map[string]object.RequestSettings{
					"1": {Data: []byte("1")},
				},
			},
			requestMethod: http.MethodPost,
		},
		{
			name: "GET object",
			request: TestRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Keys:       []string{"1"},
			},
			requestMethod: http.MethodGet,
		},
		{
			name: "DELETE object",
			request: TestRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Keys:       []string{"1"},
			},
			requestMethod: http.MethodDelete,
		},
		{
			name: "GET object w error",
			request: TestRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Keys:       []string{"1"},
			},
			requestMethod: http.MethodGet,
		},
		{
			name: "DELETE collection",
			request: TestRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodDelete,
		},
		{
			name: "GET collection w error",
			request: TestRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := testClient.doRequest(t, test.requestMethod, test.request)
			fmt.Printf("response: %v\n", response)
			data, code := response.DataAndCode()
			fmt.Printf("response data: %v, code: %d\n", string(data), code)
			require.Equal(t, test.wantErr, code != http.StatusOK)
		})
	}
}

func TestTTL(t *testing.T) {
	testClient := TestClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: "http://localhost:8081/",
	}

	tests := []struct {
		name          string
		request       TestRequest
		requestMethod string
		sleepTime     time.Duration
		wantErr       bool
	}{
		{
			name: "POST with timeout",
			request: TestRequest{
				Type: handlers.TypeObject,
				Objects: map[string]object.RequestSettings{
					"1": {
						Data:    []byte("1"),
						Timeout: 5},
				},
			},
			sleepTime:     time.Second * 6,
			requestMethod: http.MethodPost,
		},
		{
			name: "GET with error",
			request: TestRequest{
				Type: handlers.TypeObject,
				Keys: []string{"1"},
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
		{
			name: "POST with deadline",
			request: TestRequest{
				Type: handlers.TypeObject,
				Objects: map[string]object.RequestSettings{
					"1": {
						Data:     []byte("1"),
						Deadline: time.Now().Add(time.Second * 5)},
				},
			},
			sleepTime:     time.Second * 6,
			requestMethod: http.MethodPost,
		},
		{
			name: "GET with error",
			request: TestRequest{
				Type: handlers.TypeObject,
				Keys: []string{"1"},
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := testClient.doRequest(t, test.requestMethod, test.request)
			fmt.Printf("response: %v\n", response)
			responses := response.(*handlers.Responses)
			fmt.Printf("response data: %v\n", (*responses)[0].Data)
			require.NotEqual(t, test.wantErr, (*responses)[0].Success)
			time.Sleep(test.sleepTime)
		})
	}
}
