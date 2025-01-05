package http_errors

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrInternalServerError = errors.New("Internal Server Error")
var ErrHttpUnknownInternalServerError = errors.New("unknown internal server error")
var ErrHTTPNotFound = errors.New("Not found")
var ErrMethodNotAllowed = errors.New("method not allowed")

type HttpError struct {
	Status int32
	Err    error
}

func (he HttpError) Error() string {
	return fmt.Sprintf(`{"status": %d, "error": "%s"}`, he.Status, he.Err)
}

func (he HttpError) WriteError(w http.ResponseWriter) {
	w.WriteHeader(int(he.Status))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(he.Error()))
}

func WriteError(err error, w http.ResponseWriter) {
	t, ok := err.(HttpError)
	if !ok {
		t = InternalServerError(err)
	}
	t.WriteError(w)
}

func BadRequest(err error) HttpError {
	return HttpError{
		Status: http.StatusBadRequest,
		Err:    err,
	}
}

func InternalServerError(err error) HttpError {
	return HttpError{
		Status: http.StatusInternalServerError,
		Err:    err,
	}
}

func UnknownInternalServerError() HttpError {
	return InternalServerError(ErrHttpUnknownInternalServerError)
}

func MethodNotAllowed() HttpError {
	return HttpError{
		Status: http.StatusMethodNotAllowed,
		Err:    ErrMethodNotAllowed,
	}
}

func NotFound() HttpError {
	return HttpError{
		Status: http.StatusNotFound,
		Err:    ErrHTTPNotFound,
	}
}
