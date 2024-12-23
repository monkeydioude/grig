package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError struct {
	Status int32
	Err    error
}

func (he HttpError) Error() string {
	return fmt.Sprintf(`"status": %d, "error": "%s"}`, he.Status, he.Err)
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

func InternalServerError(err error) HttpError {
	return HttpError{
		Status: 500,
		Err:    err,
	}
}

func UnknownInternalServerError() HttpError {
	return InternalServerError(ErrHttpUnknownInternalServerError)
}

func NotFound() HttpError {
	return HttpError{
		Status: 404,
		Err:    errors.New("NOT FOUND"),
	}
}