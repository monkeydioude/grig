package with

import (
	"monkeydioude/grig/internal/html/element"
	"monkeydioude/grig/internal/service/server"
	"net/http"
)

type NavWrapper element.Nav

func (nw NavWrapper) WithNav(handler func(w http.ResponseWriter, r *http.Request, nav element.Nav) error) server.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return handler(w, r, element.Nav(nw).WithCurent(r.URL.Path))
	}
}
