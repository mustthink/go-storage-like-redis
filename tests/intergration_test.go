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

	"github.com/mustthink/go-storage-like-redis/internal/handlers"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

type TestClient struct {
	client *http.Client
	url    string
}

func (c TestClient) doRequest(t *testing.T, requestMethod string, request handlers.RequestProcessor) handlers.Response {
	var response handlers.Response
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

	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("couldn't unmarshal response w err: %s", err.Error())
	}
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
		request       handlers.RequestProcessor
		requestMethod string
		wantErr       bool
	}{
		{
			name: "POST collection",
			request: handlers.PostRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodPost,
		},
		{
			name: "GET collection",
			request: handlers.GetRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodGet,
		},
		{
			name: "POST object",
			request: handlers.PostRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Key:        "1",
				Object: object.RequestSettings{
					Data: []byte("1"),
				},
			},
			requestMethod: http.MethodPost,
		},
		{
			name: "GET object",
			request: handlers.GetRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Key:        "1",
			},
			requestMethod: http.MethodGet,
		},
		{
			name: "DELETE object",
			request: handlers.DeleteRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Key:        "1",
			},
			requestMethod: http.MethodDelete,
		},
		{
			name: "GET object w error",
			request: handlers.GetRequest{
				Type:       handlers.TypeObject,
				Collection: "test",
				Key:        "1",
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
		{
			name: "DELETE collection",
			request: handlers.DeleteRequest{
				Type:       handlers.TypeCollection,
				Collection: "test",
			},
			requestMethod: http.MethodDelete,
		},
		{
			name: "GET collection w error",
			request: handlers.DeleteRequest{
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
			fmt.Printf("response data: %v\n", string(response.Data))
			require.NotEqual(t, response.Success, test.wantErr)
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
		request       handlers.RequestProcessor
		requestMethod string
		sleepTime     time.Duration
		wantErr       bool
	}{
		{
			name: "POST with timeout",
			request: handlers.PostRequest{
				Type: handlers.TypeObject,
				Key:  "1",
				Object: object.RequestSettings{
					Data:    []byte("1"),
					Timeout: 5,
				},
			},
			sleepTime:     time.Second * 6,
			requestMethod: http.MethodPost,
		},
		{
			name: "GET with error",
			request: handlers.GetRequest{
				Type: handlers.TypeObject,
				Key:  "1",
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
		{
			name: "POST with deadline",
			request: handlers.PostRequest{
				Type: handlers.TypeObject,
				Key:  "1",
				Object: object.RequestSettings{
					Data:     []byte("1"),
					Deadline: time.Now().Add(time.Second * 5),
				},
			},
			sleepTime:     time.Second * 6,
			requestMethod: http.MethodPost,
		},
		{
			name: "GET with error",
			request: handlers.GetRequest{
				Type: handlers.TypeObject,
				Key:  "1",
			},
			requestMethod: http.MethodGet,
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := testClient.doRequest(t, test.requestMethod, test.request)
			fmt.Printf("response: %v\n", response)
			fmt.Printf("response data: %v\n", string(response.Data))
			require.NotEqual(t, response.Success, test.wantErr)
			time.Sleep(test.sleepTime)
		})
	}
}
