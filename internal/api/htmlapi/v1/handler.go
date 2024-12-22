package v1

import (
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/tiger/assert"
)

type Handler struct {
	Layout *server.Layout
}

// func (h Handler) Ok()

func New(layout *server.Layout) Handler {
	assert.NotNil(layout)
	return Handler{
		Layout: layout,
	}
}
