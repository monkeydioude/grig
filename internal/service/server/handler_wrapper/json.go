package with

import (
	"encoding/json"
	"fmt"
	"io"
	"monkeydioude/grig/internal/service/server"
	"net/http"
)

func JsonPayload[T any](handler func(w http.ResponseWriter, r *http.Request, payload *T) error) server.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("WithPayload: %w", err)
		}
		var payload T
		err = json.Unmarshal(data, &payload)
		if err != nil {
			return fmt.Errorf("WithPayload: %w", err)
		}
		return handler(w, r, &payload)
	}
}
