package api

import "net/http"

func Healthcheck(w http.ResponseWriter, req *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"health\": \"OK\"}"))
	return nil
}
