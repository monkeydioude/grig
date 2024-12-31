package main

import (
	"monkeydioude/grig/internal/api"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"
	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	"monkeydioude/grig/internal/html/element"
	"monkeydioude/grig/internal/service/server"
	with "monkeydioude/grig/internal/service/server/handler_wrapper"
	"monkeydioude/grig/internal/service/server/middleware"
	"monkeydioude/grig/internal/tiger/assert"
	"net/http"
)

func navWrapper() with.NavWrapper {
	return with.NavWrapper(element.Nav{
		Links: []element.Link{
			{
				Href:   "/",
				Text:   element.Text("Home"),
				Target: element.Self,
			},
			{
				Href:   "/capybara",
				Text:   element.Text("Capybara"),
				Target: element.Self,
			},
		},
	})
}

func routing(layout *server.Layout) http.Handler {
	assert.NotNil(layout)
	mux := http.NewServeMux()
	serveStatic(mux)
	json := jsonApi.New(layout)
	html := htmlApi.New(layout)
	nw := navWrapper()

	// generale routes definition
	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))

	// json api routes definition
	mux.HandleFunc("/api/v1/capybara", layout.Post(with.JsonPayload(json.CapybaraSave)))

	// html routes definition
	mux.HandleFunc("/capybara", layout.Get(nw.WithNav(html.CapybaraList)))
	mux.HandleFunc("/", layout.Get(nw.WithNav(html.Index)))

	// Apply the middleware to your server
	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}
