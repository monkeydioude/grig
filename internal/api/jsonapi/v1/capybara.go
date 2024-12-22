package v1

import (
	"fmt"
	"io"
	"net/http"
)

func (h Handler) CapybaraCreate(w http.ResponseWriter, r *http.Request) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
