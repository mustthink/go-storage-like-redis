package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

const (
	// types
	TypeCollection = "collection"
	TypeObject     = "object"
)

func Handler(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	var processor RequestProcessor
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		processor = &PostRequest{}
	case http.MethodGet:
		processor = &GetRequest{}
	case http.MethodDelete:
		processor = &DeleteRequest{}
	}

	response := getResponse(r.Body, processor, s)
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(response.Error.Code)
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseData)
}

func getResponse(requestBody io.ReadCloser, processor RequestProcessor, s storage.Storage) Response {
	var errMsg errors.Error
	body, err := io.ReadAll(requestBody)
	if err != nil {
		errMsg = errors.ErrMsgReadBody(err)
		return ResponseByError(errMsg)
	}

	if err := json.Unmarshal(body, processor); err != nil {
		errMsg = errors.ErrMsgUnmarshalBody(err)
		return ResponseByError(errMsg)
	}

	return processor.Process(s)
}
