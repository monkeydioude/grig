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

	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))
	routing_json(layout, mux)
	routing_html(layout, mux)
	routing_blocks(layout, mux)

	// Apply middlewares to server
	app := middleware.Mux(mux)
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiXRequestID,
		middleware.JsonApiLogRequest,
	)
	return app
}

func routing_json(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	json := jsonApi.New(layout)
	mux.HandleFunc("/api/v1/capybara", layout.Post(with.JsonPayload(json.CapybaraSave)))
	mux.HandleFunc("/api/v1/josuke", layout.Post(with.JsonPayload(json.JosukeSave)))
	mux.HandleFunc("/api/v1/services", layout.Post(with.JsonPayload(json.ServicesSave)))
}

func routing_html(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	html := htmlApi.New(layout)
	nw := with.NewNavWrapper()

	mux.HandleFunc("/", layout.Get(nw.WithNav(html.Index, element.Link{Href: "/", Text: element.Text("Index")})))
	mux.HandleFunc("/capybara", layout.Get(nw.WithNav(html.CapybaraList, element.Link{Href: "/capybara"})))
	mux.HandleFunc("/josuke", layout.Get(nw.WithNav(html.JosukeList, element.Link{Href: "/josuke"})))
	mux.HandleFunc("/services", layout.Get(nw.WithNav(html.ServicesList, element.Link{Href: "/services"})))
	mux.HandleFunc("/services/by_filepath", layout.Post(with.JsonPayload(html.AddServiceByFilepath)))
}

func routing_blocks(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	html := htmlApi.New(layout)
	// capybara blocks
	mux.HandleFunc("/capybara/service/block", layout.Get(html.CapybaraServiceBlock))

	// josuke blocks
	mux.HandleFunc("/josuke/hook/block", layout.Get(html.JosukeHookBlock))
	mux.HandleFunc("/josuke/deployment/block", layout.Get(html.JosukeDeploymentBlock))
	mux.HandleFunc("/josuke/branch/block", layout.Get(html.JosukeBranchBlock))
	mux.HandleFunc("/josuke/action/block", layout.Get(html.JosukeActionBlock))
	mux.HandleFunc("/josuke/command/block", layout.Get(html.JosukeCommandBlock))

	// sys services blocks
	mux.HandleFunc("/services/environment/block", layout.Get(html.ServicesEnvironmentBlock))
	mux.HandleFunc("/services/service/block", layout.Get(html.ServicesServiceBlock))
}
