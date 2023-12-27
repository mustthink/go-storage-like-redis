package handlers

import (
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type (
	// DELETE one object request
	DeleteRequest struct {
		Type       string `json:"type"`
		Collection string `json:"collection"`
		Key        string `json:"key"`
	}
)

func (r DeleteRequest) Process(s storage.Storage) Response {
	switch r.Type {
	case TypeCollection:
		return deleteCollectionResponse(r.Collection, s)
	case TypeObject:
		return deleteObjectResponse(r.Collection, r.Key, s)
	}

	errMsg := errors.ErrMsgUnknownType(r.Type)
	return ResponseByError(errMsg)
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
