package handlers

import (
	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type RequestProcessor interface {
	Process(s storage.Storage) Response
}

// response for all requests
type Response struct {
	Data    []byte       `json:"data"`
	Success bool         `json:"success"`
	Error   errors.Error `json:"error"`
}

func ResponseByError(errMsg errors.Error) Response {
	return Response{
		Error: errMsg,
	}
}
