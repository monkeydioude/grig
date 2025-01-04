package middleware

import "net/http"

type API struct {
	chain http.Handler
}

func Mux(handler http.Handler) *API {
	return &API{handler}
}

// Use declares a middleware. `Use` allows to freely chose
// when the middleware should apply, before or after (or both) the main HTTP handler function.
// At the cost of having to call ServeHttp of the http.Handler.
// To be used when needing to act before and after the main HTTP handler function.
// Some examples:
// - pass an x request id between services
// - middleware using `defer`, such as panic recover middlewares
// - log incoming/outgoing requests
func (a *API) Use(handlers ...func(http.Handler) http.Handler) {
	for _, handler := range handlers {
		a.chain = handler(a.chain)
	}
}

// UseBefore allows for a simpler way of applying a middleware.
// The middleware will always be called __BEFORE__ the main HTTP handler function.
func (a *API) UseBefore(handlers ...func(http.ResponseWriter, *http.Request)) {
	for _, handler := range handlers {
		a.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler(w, r)
				h.ServeHTTP(w, r)
			})
		})
	}
}

// UseAfter allows for a simpler way of applying a middleware.
// The middleware will always be called __AFTER__ the main HTTP handler function.
func (a *API) UseAfter(handlers ...func(http.ResponseWriter, *http.Request)) {
	for _, handler := range handlers {
		a.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h.ServeHTTP(w, r)
				handler(w, r)
			})
		})
	}
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.chain.ServeHTTP(w, r)
}
