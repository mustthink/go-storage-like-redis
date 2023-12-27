package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mustthink/go-storage-like-redis/internal/handlers"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

// SetWithTimeout - set object with timeout in seconds
func (c *Client) SetWithTimeout(collection, key string, obj any, timeout time.Duration) error {
	objectData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("couldn't prepare object w err: %s", err.Error())
	}

	objSettings := object.RequestSettings{
		Data:    objectData,
		Timeout: timeout,
	}

	return c.set(collection, key, objSettings)
}

// SetWithDeadline - set object with deadline
func (c *Client) SetWithDeadline(collection, key string, obj any, deadline time.Time) error {
	objectData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("couldn't prepare object w err: %s", err.Error())
	}

	objSettings := object.RequestSettings{
		Data:     objectData,
		Deadline: deadline,
	}

	return c.set(collection, key, objSettings)
}

// SetTimeless - set object without any timeout
func (c *Client) SetTimeless(collection, key string, obj any) error {
	objectData, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("couldn't prepare object w err: %s", err.Error())
	}

	objSettings := object.RequestSettings{
		Data:     objectData,
		Timeless: true,
	}

	return c.set(collection, key, objSettings)
}

func (c *Client) Get(collection, key string, obj any) error {
	request := handlers.GetRequest{
		Type:       handlers.TypeObject,
		Collection: collection,
		Key:        key,
	}

	response, err := c.do(http.MethodGet, request)
	if err != nil {
		return err
	}

	err = json.Unmarshal(response.Data, obj)
	return err
}

func (c *Client) Delete(collection, key string) error {
	request := handlers.DeleteRequest{
		Type:       handlers.TypeObject,
		Collection: collection,
		Key:        key,
	}

	_, err := c.do(http.MethodDelete, request)
	return err
}

func (c *Client) set(collection, key string, objSettings object.RequestSettings) error {
	request := handlers.PostRequest{
		Type:       handlers.TypeObject,
		Collection: collection,
		Key:        key,
		Object:     objSettings,
	}

	_, err := c.do(http.MethodPost, request)
	return err
}
