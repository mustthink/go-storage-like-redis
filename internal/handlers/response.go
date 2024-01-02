package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mustthink/go-storage-like-redis/internal/errors"
)

type (
	DataCode interface {
		DataAndCode() (data []byte, code int)
	}

	// Response - single response
	Response struct {
		Data    []byte       `json:"data"`
		Success bool         `json:"success"`
		Error   errors.Error `json:"error"`
	}

	// Responses - slice of responses
	Responses []Response
)

func ResponseByError(errMsg errors.Error) Response {
	return Response{
		Error: errMsg,
	}
}

func (r Response) DataAndCode() ([]byte, int) {
	responseData, err := json.Marshal(r)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	if r.Error.Code != 0 {
		return responseData, r.Error.Code
	}

	return responseData, http.StatusOK
}

func (r Responses) DataAndCode() ([]byte, int) {
	responseData, err := json.Marshal([]Response(r))
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return responseData, http.StatusOK
}
