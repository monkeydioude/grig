package v1

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) CapybaraList(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/capybara/"))
	if err != nil {
		return errors.Wrap(err, "invalid capybara ID")
	}
	ref, ok := h.Layout.ServerConfig.CapybaraRefByID(id)
	if !ok {
		return errors.Wrap(fmt.Errorf("ID out of bounds"), "invalid capybara ID")
	}
	layout := layouts.Main(nav, pages.CapybaraList(ref))
	layout.Render(r.Context(), w)
	return nil
}

// CapybaraTLSHostBlock renders one more empty row for the proxy's TLS
// host list.
func (h Handler) CapybaraTLSHostBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		index = 0
	}
	parent := r.URL.Query().Get("parent_name")
	if parent == "" {
		parent = "proxy"
	}
	return blocks.CapybaraTLSHost(parent, index, "").Render(r.Context(), w)
}

// CapybaraExtraFieldBlock renders one empty key/value row for a key grig
// does not model. The row carries no field name until the browser sees a
// valid key, so the service index is resolved client side.
func (h Handler) CapybaraExtraFieldBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	return blocks.CapybaraNewExtraField().Render(r.Context(), w)
}

func (h Handler) CapybaraServiceBlock(w http.ResponseWriter, r *http.Request, _ *slog.Logger) error {
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		it, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.Wrap(err, "CapybaraServiceBlock")
		}
		index = it
	}
	return blocks.CapybaraService(pages.GetServiceInputName, index, model.ServiceDefinition{}).Render(r.Context(), w)
}
