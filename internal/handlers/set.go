package handlers

import (
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

type (
	// PostRequest - POST, PUT collection/objects request, set objects to storage
	PostRequest struct {
		Collection         string                            `json:"collection"`
		Objects            map[string]object.RequestSettings `json:"objects"`
		ObjectsWithoutKeys []object.RequestSettings          `json:"objects_without_keys"`
	}
)

func (r PostRequest) ProcessCollection(s storage.Storage) Response {
	return postCollectionResponse(r.Collection, s)
}

func (r PostRequest) ProcessObjects(s storage.Storage) []Response {
	responses := make([]Response, 0, len(r.Objects)+len(r.ObjectsWithoutKeys))
	for key, objSettings := range r.Objects {
		responsePart := postObjectResponse(r.Collection, key, objSettings, s)
		responses = append(responses, responsePart)
	}

	for _, objSettings := range r.ObjectsWithoutKeys {
		responsePart := postObjectResponse(r.Collection, "", objSettings, s)
		responses = append(responses, responsePart)
	}
	return responses
}

func postCollectionResponse(name string, s storage.Storage) Response {
	err := s.NewCollection(name)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}

	return Response{
		Success: true,
	}
}

func postObjectResponse(collectionName, key string, objSettings object.RequestSettings, s storage.Storage) Response {
	if key == "" {
		newKey, err := objSettings.NewKey()
		if err != nil {
			errMsg := errors.ErrMsgByError(err, http.StatusInternalServerError)
			return ResponseByError(errMsg)
		}

		key = newKey
	}

	err := storage.SetObject(s, collectionName, key, objSettings)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}

	return Response{
		Success: true,
		Data:    []byte(key),
	}
}
