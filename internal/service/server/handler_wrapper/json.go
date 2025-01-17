package with

import (
	"encoding/json"
	"io"
	"log/slog"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/server"
	"net/http"
	"regexp"
)

func cleanJSONFromNull(input []byte) []byte {
	re := regexp.MustCompile(`,\s*null|null,\s*|\[null\]`)
	return re.ReplaceAll(input, []byte{})
}

func JsonPayload[T any](handler func(http.ResponseWriter, *http.Request, *slog.Logger, *T) error) server.Handler {
	return func(w http.ResponseWriter, r *http.Request, logger *slog.Logger) error {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.Wrap(err, "JsonPayload")
		}
		data = cleanJSONFromNull(data)
		var payload T
		err = json.Unmarshal(data, &payload)
		if err != nil {
			return errors.Wrap(err, "JsonPayload")
		}
		return handler(w, r, logger, &payload)
	}
}
