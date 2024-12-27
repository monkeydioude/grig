package main

import (
	"fmt"
	"monkeydioude/grig/internal/api"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"

	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/service/server/middleware"
	"monkeydioude/grig/internal/tiger/assert"
	"net/http"
	"os"
	"time"
)

func apiRouting(layout *server.Layout) http.Handler {
	assert.NotNil(layout)
	mux := http.NewServeMux()
	json := jsonApi.New(layout)
	html := htmlApi.New(layout)

	// generale routes definition
	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))

	// json api routes definition
	mux.HandleFunc("/api/v1/capybara/create", layout.Post(json.CapybaraCreate))

	// html routes definition
	mux.HandleFunc("/capybara", layout.Get(html.CapybaraList))
	mux.HandleFunc("/", layout.Get(html.Index))

	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}

func setupApiServer(layout *server.Layout) *http.Server {
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
