package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/handlers"
)

func (c *Client) do(method string, request handlers.RequestProcessor) (handlers.Response, error) {
	var response handlers.Response
	requestBody, err := json.Marshal(request)
	if err != nil {
		err = fmt.Errorf("couldn't marshal w err: %s", err.Error())
		return response, err
	}

	req, err := http.NewRequest(method, c.url, bytes.NewBuffer(requestBody))
	if err != nil {
		err = fmt.Errorf("couldn't create new request w err: %s", err.Error())
		return response, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		err = fmt.Errorf("couldn't do request w err: %s", err.Error())
		return response, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("couldn't read response.Body w err: %s", err.Error())
		return response, err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		err = fmt.Errorf("couldn't unmarshal response w err: %s", err.Error())
		return response, err
	}

	if !response.Success {
		return response, fmt.Errorf("got an error in reponse: %s", response.Error.Message)
	}

	return response, nil
}
