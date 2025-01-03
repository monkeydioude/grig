package main

import (
	"monkeydioude/grig/internal/api"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"
	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	element "monkeydioude/grig/internal/html/elements"
	"monkeydioude/grig/internal/service/server"
	with "monkeydioude/grig/internal/service/server/handler_wrapper"
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
	nw := with.NewNavWrapper()

	// generale routes definition
	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))

	// json api routes definition
	mux.HandleFunc("/api/v1/capybara", layout.Post(with.JsonPayload(json.CapybaraSave)))

	// html routes definition
	mux.HandleFunc("/", layout.Get(nw.WithNav(html.Index, element.Link{Href: "/", Text: element.Text("Index")})))
	mux.HandleFunc("/capybara", layout.Get(nw.WithNav(html.CapybaraList, element.Link{Href: "/capybara"})))
	mux.HandleFunc("/josuke", layout.Get(nw.WithNav(html.JosukeList, element.Link{Href: "/josuke"})))

	// Apply the middleware to your server
	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}
