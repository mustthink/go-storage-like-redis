package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type (
	// GetRequest - GET collection/objects request
	GetRequest struct {
		Collection string   `json:"collection"`
		Keys       []string `json:"keys"`
	}
)

func (r GetRequest) ProcessCollection(s storage.Storage) Response {
	return getCollectionResponse(r.Collection, s)
}

func (r GetRequest) ProcessObjects(s storage.Storage) []Response {
	var responses = make([]Response, 0, len(r.Keys))

	for _, key := range r.Keys {
		responsePart := getObjectResponse(r.Collection, key, s)
		responses = append(responses, responsePart)
	}

	return responses
}

func getCollectionResponse(name string, s storage.Storage) Response {
	collection, err := s.GetCollection(name)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}

	data, err := json.Marshal(collection)
	if err != nil {
		errMsg := errors.ErrMsgMarshalCollection(err)
		return ResponseByError(errMsg)
	}

	return Response{
		Data:    data,
		Success: true,
	}
}

func getObjectResponse(collectionName, key string, s storage.Storage) Response {
	object, err := storage.GetObject(s, collectionName, key)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}

	return Response{
		Data:    object.Binary(),
		Success: true,
	}
}
