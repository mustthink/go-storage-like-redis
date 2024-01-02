package errors

import (
	"fmt"
	"net/http"
)

// errors
func ErrNoCollection(collectionName string) error {
	return fmt.Errorf("no collection w name: %s", collectionName)
}

func ErrCollectionAlreadyExist(collectionName string) error {
	return fmt.Errorf("collection w name: %s already exist", collectionName)
}

func ErrNoObject(objectKey string) error {
	return fmt.Errorf("no object w key: %s", objectKey)
}

func ErrEmptyField(field string) error {
	return fmt.Errorf("%s is empty", field)
}

var (
	ErrDeleteDefaultCollection = fmt.Errorf("couldn't delete default collection")
	ErrMaxCollectionsCount     = fmt.Errorf("too many collections")
)

// error struct for response
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func ErrMsgReadBody(err error) Error {
	return Error{
		Message: fmt.Sprintf("couldn't read request.Body w err: %s", err.Error()),
		Code:    http.StatusBadRequest,
	}
}

func ErrMsgUnmarshalBody(err error) Error {
	return Error{
		Message: fmt.Sprintf("couldn't unmarshal request.Body w err: %s", err.Error()),
		Code:    http.StatusBadRequest,
	}
}

func ErrMsgMarshalCollection(err error) Error {
	return Error{
		Message: fmt.Sprintf("couldn't marshal collection w err: %s", err.Error()),
		Code:    http.StatusInternalServerError,
	}
}

func ErrMsgUnknownType(type_ string) Error {
	return Error{
		Message: fmt.Sprintf("unknown type: %s", type_),
		Code:    http.StatusBadRequest,
	}
}

func ErrMsgByError(err error, code int) Error {
	return Error{
		Message: err.Error(),
		Code:    code,
	}
}

var (
	ErrInvalidPostRequest = Error{
		Message: "request MUST have ONLY map[key]object or array of objects",
		Code:    http.StatusBadRequest,
	}
)
