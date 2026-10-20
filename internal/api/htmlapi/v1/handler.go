package v1

import (
	"monkeydioude/grig/internal/service/db_schema"
	"monkeydioude/grig/internal/service/forge"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/tiger/assert"
)

type Handler struct {
	Layout *server.Layout[config.ServerConfig]
	// Forges are the git hosts new service repositories can be created on.
	Forges forge.Registry
	// Schemas creates the database schema and user of new services.
	Schemas db_schema.Admin
}

// func (h Handler) Ok()

func New(layout *server.Layout[config.ServerConfig]) Handler {
	assert.NotNil(layout)
	return Handler{
		Layout: layout,
	}
}

func (h Handler) WithForges(forges forge.Registry) Handler {
	h.Forges = forges
	return h
}

func (h Handler) WithSchemas(schemas db_schema.Admin) Handler {
	h.Schemas = schemas
	return h
}
