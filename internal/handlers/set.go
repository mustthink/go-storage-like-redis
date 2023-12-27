package handlers

import (
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
	"github.com/mustthink/go-storage-like-redis/internal/storage/object"
)

type (
	// POST, PUT one object request
	PostRequest struct {
		Type       string                 `json:"type"`
		Collection string                 `json:"collection"`
		Key        string                 `json:"key"`
		Object     object.RequestSettings `json:"object"`
	}
)

func (r PostRequest) Process(s storage.Storage) Response {
	switch r.Type {
	case TypeCollection:
		return postCollectionResponse(r.Collection, s)
	case TypeObject:
		return postObjectResponse(r.Collection, r.Key, r.Object, s)
	}

	errMsg := errors.ErrMsgUnknownType(r.Type)
	return ResponseByError(errMsg)
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
	err := storage.SetObject(s, collectionName, key, objSettings)
	if err != nil {
		errMsg := errors.ErrMsgByError(err, http.StatusBadRequest)
		return ResponseByError(errMsg)
	}

	return Response{
		Success: true,
	}
}
