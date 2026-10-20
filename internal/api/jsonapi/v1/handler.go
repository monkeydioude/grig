package v1

import (
	"monkeydioude/grig/internal/service/restart"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/tiger/assert"
)

type Handler struct {
	Layout *server.Layout[config.ServerConfig]
	// Restarter applies a saved config by restarting the unit reading it.
	Restarter restart.Restarter
}

func New(layout *server.Layout[config.ServerConfig]) Handler {
	assert.NotNil(layout)
	return Handler{
		Layout: layout,
	}
}

func (h Handler) WithRestarter(restarter restart.Restarter) Handler {
	h.Restarter = restarter
	return h
}
