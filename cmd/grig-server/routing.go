package main

import (
	"monkeydioude/grig/internal/api"
	cmdApi "monkeydioude/grig/internal/api/cmdapi/v1"
	htmlApi "monkeydioude/grig/internal/api/htmlapi/v1"
	jsonApi "monkeydioude/grig/internal/api/jsonapi/v1"
	"monkeydioude/grig/internal/service/server/config"
	with "monkeydioude/grig/internal/service/server/handler_wrapper"
	element "monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/server/middleware"
	"monkeydioude/grig/pkg/tiger/assert"
	"monkeydioude/grig/pkg/urls"
	"net/http"

	"github.com/a-h/templ"
)

func routing(layout *server.Layout[config.ServerConfig]) http.Handler {
	assert.NotNil(layout)
	mux := http.NewServeMux()
	serveStatic(mux)

	mux.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))
	routing_json(layout, mux)
	routing_html(layout, mux)
	routing_command(layout, mux)

	prefixer := layout.ServerConfig.Prefixer()

	// Apply middlewares to server. The prefixer goes on first so every
	// template can build links under the base path.
	app := middleware.Mux(mountUnderBasePath(layout, mux, prefixer))
	app.Use(
		middleware.PanicRecover,
		middleware.JsonApiXRequestID,
		middleware.JsonApiLogRequest,
		urls.Middleware(prefixer),
	)
	return app
}

// mountUnderBasePath serves the app's routes under the configured
// sub-path. Routes stay registered relative to the app root and the
// prefix is stripped on the way in, so handlers and patterns are
// unaffected by where the app is mounted.
func mountUnderBasePath(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
	prefixer urls.Prefixer,
) http.Handler {
	base := prefixer.Base()
	if base == "" {
		return mux
	}
	root := http.NewServeMux()
	root.Handle(base+"/", http.StripPrefix(base, mux))
	root.Handle(base, http.RedirectHandler(base+"/", http.StatusMovedPermanently))
	// keep an unprefixed healthcheck for probes that do not know the base path
	root.HandleFunc("/healthcheck", layout.Get(api.Healthcheck))
	// anything else sends the browser to the app root
	root.Handle("/", http.RedirectHandler(base+"/", http.StatusFound))
	return root
}

func routing_json(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	json := jsonApi.New(layout)
	mux.HandleFunc("/api/v1/capybara/{id}", layout.Post(with.JsonPayload(json.CapybaraSave)))
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
	for _, ref := range layout.ServerConfig.CapybaraRefs() {
		mux.HandleFunc(ref.URL(), layout.Get(nw.WithNav(html.CapybaraList, element.Link{
			Href: templ.SafeURL(ref.URL()),
			Text: element.Text(ref.Name),
		})))
	}
	mux.HandleFunc("/josuke", layout.Get(nw.WithNav(html.JosukeList, element.Link{Href: "/josuke"})))
	mux.HandleFunc("/services", layout.Get(nw.WithNav(html.ServicesList, element.Link{Href: "/services"})))
	mux.HandleFunc("/services/by_filepath", layout.Post(with.JsonPayload(html.AddServiceByFilepath)))
	routing_html_blocks(layout, mux)
}

func routing_html_blocks(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	html := htmlApi.New(layout)
	// capybara blocks
	mux.HandleFunc("/capybara/service/block", layout.Get(html.CapybaraServiceBlock))
	mux.HandleFunc("/capybara/tls_host/block", layout.Get(html.CapybaraTLSHostBlock))
	mux.HandleFunc("/capybara/extra/block", layout.Get(html.CapybaraExtraFieldBlock))

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

func routing_command(
	layout *server.Layout[config.ServerConfig],
	mux *http.ServeMux,
) {
	cmd := cmdApi.New(layout)
	mux.HandleFunc("/cmd/services/restart/{service}", layout.Post(cmd.CmdServiceRestart))
	mux.HandleFunc("/cmd/services/restart", layout.Post(cmd.CmdServiceRestartAll))
}
