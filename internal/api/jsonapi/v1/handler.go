package v1

import (
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/tiger/assert"
)

type Handler struct {
	Layout *server.Layout[config.ServerConfig]
}

func New(layout *server.Layout[config.ServerConfig]) Handler {
	assert.NotNil(layout)
	return Handler{
		Layout: layout,
	}
}
