package v1

import (
	"monkeydioude/grig/internal/service/compose"
	"monkeydioude/grig/internal/service/db_schema"
	"monkeydioude/grig/internal/service/env_file"
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
	// Envs is where the .env files of new services are written.
	Envs env_file.Dir
	// Compose is where the docker compose services of new services are added.
	Compose compose.Dir
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

func (h Handler) WithEnvs(envs env_file.Dir) Handler {
	h.Envs = envs
	return h
}

func (h Handler) WithCompose(dir compose.Dir) Handler {
	h.Compose = dir
	return h
}

func (h Handler) WithSchemas(schemas db_schema.Admin) Handler {
	h.Schemas = schemas
	return h
}
