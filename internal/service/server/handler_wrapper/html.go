package with

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	element "monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/tiger/assert"
	"net/http"
	"unicode"
)

// NavWrapper extends `element.Nav` and is used in the routing phase, to build
// the site's navigation.
// Does some magic, like trying to derive the page's name
// from the link.
//
// `element.Nav` contains a slice of `element.Link`.
// `element.Link.Href` cannot be empty.
// If `element.Link.Text` is empty, `element.Link.Href` cannot be empty or a single '/'.
type NavWrapper element.Nav

func NewNavWrapper() NavWrapper {
	return NavWrapper(element.Nav{})
}

func transform(input string) string {
	if len(input) == 0 || input[0] != '/' {
		return input // Return as is if input is empty or doesn't start with /
	}
	// Remove the leading /
	trimmed := input[1:]
	if len(trimmed) == 0 {
		return "" // Return empty if there's nothing after the /
	}
	// Capitalize the first letter
	runes := []rune(trimmed)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (nw *NavWrapper) WithNav(handler func(w http.ResponseWriter, r *http.Request, nav element.Nav) error, link element.Link) server.Handler {
	assert.NotEmpty(string(link.Href), errors.ErrEmptyLinkHref)
	if link.Target == "" {
		link.Target = element.Self
	}
	if link.Text == nil || link.Text.String() == "" {
		link.Text = element.Text(transform(string(link.Href)))
	}
	assert.NotEmpty(link.Text.String(), errors.ErrEmptyLinkText, fmt.Errorf("Here element.Link.Href is '%v'", link.Href))
	nw.Links = append(nw.Links, link)
	return func(w http.ResponseWriter, r *http.Request) error {
		// w.WriteHeader(202)
		return handler(w, r, element.Nav(*nw).WithCurent(r.URL.Path))
	}
}
