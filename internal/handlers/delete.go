package handlers

import (
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type (
	// DeleteRequest - DELETE collection/objects request
	DeleteRequest struct {
		Collection string   `json:"collection"`
		Keys       []string `json:"keys"`
	}
)

func (r DeleteRequest) ProcessCollection(s storage.Storage) Response {
	return deleteCollectionResponse(r.Collection, s)
}

func (r DeleteRequest) ProcessObjects(s storage.Storage) []Response {
	var responses = make([]Response, 0, len(r.Keys))

	for _, key := range r.Keys {
		responsePart := deleteObjectResponse(r.Collection, key, s)
		responses = append(responses, responsePart)
	}

	return responses
}

func deleteCollectionResponse(name string, s storage.Storage) Response {
	err := s.DeleteCollection(name)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}
	return Response{
		Success: true,
	}
}

func deleteObjectResponse(name, key string, s storage.Storage) Response {
	err := storage.DeleteObject(s, name, key)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}
	return Response{
		Success: true,
	}
}
