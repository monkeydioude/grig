package main

import (
	"monkeydioude/grig/internal/api"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"
	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/service/server/middleware"
	"monkeydioude/grig/internal/tiger/assert"
	"net/http"
)

func routing(layout *server.Layout) http.Handler {
	assert.NotNil(layout)
	mux := http.NewServeMux()
	serveStatic(mux)
	json := jsonApi.New(layout)
	html := htmlApi.New(layout)

	// generale routes definition
	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))

	// json api routes definition
	mux.HandleFunc("/api/v1/capybara", layout.Post(json.CapybaraSave))

	// html routes definition
	mux.HandleFunc("/capybara", layout.Get(html.CapybaraList))
	mux.HandleFunc("/", layout.Get(html.Index))

	// Apply the middleware to your server
	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}
