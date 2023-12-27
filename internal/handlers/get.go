package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type (
	// GET one object request
	GetRequest struct {
		Type       string `json:"type"`
		Collection string `json:"collection"`
		Key        string `json:"key"`
	}
)

func (r GetRequest) Process(s storage.Storage) Response {
	switch r.Type {
	case TypeCollection:
		return getCollectionResponse(r.Collection, s)
	case TypeObject:
		return getObjectResponse(r.Collection, r.Key, s)
	}

	errMsg := errors.ErrMsgUnknownType(r.Type)
	return ResponseByError(errMsg)
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
