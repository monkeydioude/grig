package errors

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteErrorWithOkValues(t *testing.T) {
	w := httptest.NewRecorder()

	trial := HttpError{
		Status: 999,
		Err:    errors.New("waw"),
	}
	WriteError(trial, w)

	res := w.Result()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, 999, res.StatusCode)
	assert.Equal(t, string(body), `"status": 999, "error": "waw"}`)
}

func TestWriteErrorWithWeirdErr(t *testing.T) {
	w := httptest.NewRecorder()

	WriteError(errors.New("kekw"), w)

	res := w.Result()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, 500, res.StatusCode)
	assert.Equal(t, string(body), `"status": 500, "error": "kekw"}`)
}
