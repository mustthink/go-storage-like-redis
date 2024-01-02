package handlers

import (
	"encoding/json"
	"fmt"
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

type (
	RequestProcessor interface {
		ProcessCollection(s storage.Storage) Response
		ProcessObjects(s storage.Storage) []Response
	}

	// Request - struct of requests
	Request struct {
		Type string `json:"type"`
		RequestProcessor
	}
)

func RequestByMethod(method string) Request {
	var processor RequestProcessor
	switch method {
	case http.MethodPost, http.MethodPut:
		processor = &PostRequest{}
	case http.MethodGet:
		processor = &GetRequest{}
	case http.MethodDelete:
		processor = &DeleteRequest{}
	}
	return Request{
		RequestProcessor: processor,
	}
}

func (r *Request) readRequestBody(requestBody io.ReadCloser) (responseError errors.Error) {
	body, err := io.ReadAll(requestBody)
	if err != nil {
		responseError = errors.ErrMsgReadBody(err)
	}
	fmt.Println("body: ", string(body))

	if err := json.Unmarshal(body, r); err != nil {
		responseError = errors.ErrMsgUnmarshalBody(err)
	}

	return
}

func (r *Request) getResponse(s storage.Storage) ([]byte, int) {
	var response DataCode
	switch r.Type {
	case TypeCollection:
		response = r.ProcessCollection(s)
	case TypeObject:
		response = Responses(r.ProcessObjects(s))
	default:
		errMsg := errors.ErrMsgUnknownType(r.Type)
		response = ResponseByError(errMsg)
	}
	return response.DataAndCode()
}

func (r *Request) UnmarshalJSON(data []byte) error {
	type decoderType Request
	var decoder decoderType
	if err := json.Unmarshal(data, &decoder); err != nil {
		return err
	}
	r.Type = decoder.Type
	if err := json.Unmarshal(data, r.RequestProcessor); err != nil {
		return err
	}
	return nil
}
