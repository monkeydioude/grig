package main

import (
	"monkeydioude/grig/internal/api"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"
	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	"monkeydioude/grig/internal/service/server/config"
	with "monkeydioude/grig/internal/service/server/handler_wrapper"
	element "monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/server/middleware"
	"monkeydioude/grig/pkg/tiger/assert"
	"net/http"
)

func routing(layout *server.Layout[config.ServerConfig]) http.Handler {
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
	mux.HandleFunc("/api/v1/josuke", layout.Post(with.JsonPayload(json.JosukeSave)))

	// html routes definition
	mux.HandleFunc("/", layout.Get(nw.WithNav(html.Index, element.Link{Href: "/", Text: element.Text("Index")})))
	mux.HandleFunc("/capybara", layout.Get(nw.WithNav(html.CapybaraList, element.Link{Href: "/capybara"})))
	mux.HandleFunc("/josuke", layout.Get(nw.WithNav(html.JosukeList, element.Link{Href: "/josuke"})))

	mux.HandleFunc("/blocks/josuke/hook", layout.Get(html.JosukeHookBlock))
	mux.HandleFunc("/blocks/josuke/deployment", layout.Get(html.JosukeDeploymentBlock))
	mux.HandleFunc("/blocks/josuke/branch", layout.Get(html.JosukeBranchBlock))
	mux.HandleFunc("/blocks/josuke/action", layout.Get(html.JosukeActionBlock))
	mux.HandleFunc("/blocks/josuke/command", layout.Get(html.JosukeCommandBlock))

	mux.HandleFunc("/blocks/capybara/service", layout.Get(html.CapybaraServiceBlock))

	// Apply the middleware to your server
	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiLogRequest,
		middleware.JsonApiXRequestID,
	)
	return app
}
