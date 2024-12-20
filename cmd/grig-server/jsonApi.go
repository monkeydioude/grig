package main

import (
	"fmt"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/service/server/middleware"
	"net/http"
	"os"
	"time"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"health\": \"OK\"}"))
}

func apiRouting(layout *server.Layout) http.Handler {
	mux := http.NewServeMux()
	// routes definition
	mux.HandleFunc("/grig/healthcheck", healthcheck)

	app := middleware.Mux(mux)
	app.Use(
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}

func setupJsonApiServer(layout *server.Layout) *http.Server {
	// setup multiplexer
	mux := apiRouting(layout)
	port := consts.DEFAULT_GRIG_SERVER_PORT
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           mux,
	}
}
