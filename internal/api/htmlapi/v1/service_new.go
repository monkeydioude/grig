package v1

import (
	"cmp"
	"errors"
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/compose"
	"monkeydioude/grig/internal/service/db_schema"
	"monkeydioude/grig/internal/service/env_file"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/forge"
	"monkeydioude/grig/internal/service/repo_setup"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server/http_errors"
	"monkeydioude/grig/pkg/trans_types"
	"net/http"
	"strings"
)

func (h Handler) ServiceNew(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	hook, hookErr := h.josukeWebhook()
	page := pages.NewServiceNew(h.Forges.Names(), hook.URL, hookErr)
	if h.Schemas.Enabled() {
		password, err := db_schema.GeneratePassword()
		if err != nil {
			return err
		}
		page = page.WithDatabase(db_schema.Redacted(h.Schemas.URL), password, h.Envs.Base)
	}
	if h.Compose.Enabled() {
		page = page.WithCompose(h.Compose.Types())
	}
	return layouts.Main(nav, page).Render(r.Context(), w)
}

const dbURLKey = "DB_URL"

// schemaForm is the Database form: the schema, and the project whose
// .env receives the service's DB_URL.
type schemaForm struct {
	db_schema.Request
	Project string `json:"project"`
}

// ServiceSchemaCreate creates the database schema and user of a new
// service, writes its DB_URL to the project's .env and renders the result.
func (h Handler) ServiceSchemaCreate(w http.ResponseWriter, r *http.Request, logger *slog.Logger, form *schemaForm) error {
	if !h.Schemas.Enabled() {
		return http_errors.BadRequest(errors.New("ServiceSchemaCreate: DB_URL is not set"))
	}
	if err := form.Validate(); err != nil {
		return http_errors.BadRequest(pkgErrors.Wrap(err, "ServiceSchemaCreate"))
	}
	envPath, err := h.schemaEnvPath(form.Project)
	if err != nil {
		return err
	}
	res, err := h.Schemas.Create(r.Context(), form.Request)
	envErr := writeDBURL(envPath, res, err)
	logger.Info("schema setup", "schema", form.Schema, "user", form.User, "error", err,
		"env", envPath, "env_error", envErr)
	return blocks.SchemaSetupResult(res, err, envPath, envErr).Render(r.Context(), w)
}

// schemaEnvPath is the .env to write DB_URL to, empty when env_base_path
// is not set. It is checked before the schema is created, which would
// otherwise be created for a .env that already has a DB_URL.
func (h Handler) schemaEnvPath(project string) (string, error) {
	if !h.Envs.Enabled() {
		return "", nil
	}
	path, err := h.Envs.Path(project)
	if err != nil {
		return "", http_errors.BadRequest(pkgErrors.Wrap(err, "ServiceSchemaCreate"))
	}
	has, err := env_file.HasKey(path, dbURLKey)
	if err != nil {
		return "", pkgErrors.Wrap(err, "ServiceSchemaCreate")
	}
	if has {
		return "", http_errors.Conflict(fmt.Errorf("%s: %s %w", path, dbURLKey, env_file.ErrKeyExists))
	}
	return path, nil
}

// writeDBURL adds the created schema's DB_URL to the .env at envPath.
func writeDBURL(envPath string, res db_schema.Result, createErr error) error {
	if envPath == "" || createErr != nil {
		return nil
	}
	return env_file.Add(envPath, dbURLKey, res.ServiceURL)
}

// ServiceRepoCreate creates the repository of a new service and its
// Josuke webhook, and renders how each step went.
func (h Handler) ServiceRepoCreate(w http.ResponseWriter, r *http.Request, logger *slog.Logger, form *repo_setup.Form) error {
	req, err := form.RepoRequest(h.Forges)
	if err != nil {
		return http_errors.BadRequest(pkgErrors.Wrap(err, "ServiceRepoCreate"))
	}
	hook, hookErr := h.josukeWebhook()
	res := repo_setup.Run(r.Context(), h.Forges[form.Forge], req, hook, hookErr)
	logger.Info("repository setup", "forge", form.Forge, "repo", res.Repo.FullName,
		"repo_error", res.RepoErr, "webhook_error", res.HookErr)
	return blocks.RepoSetupResult(res).Render(r.Context(), w)
}

// josukeWebhook builds the webhook from the Josuke config on disk.
func (h Handler) josukeWebhook() (forge.Webhook, error) {
	conf := h.Layout.ServerConfig
	jk, err := file.UnmarshalFromPath[model.Josuke](conf.JosukeConfigPath)
	if err != nil {
		return forge.Webhook{}, err
	}
	return repo_setup.JosukeWebhook(jk, conf.WebhookBaseURL)
}

// composeForm is the Docker compose form, with the page's project name.
type composeForm struct {
	Project      string                `json:"project"`
	Type         string                `json:"type"`
	Port         trans_types.StringInt `json:"port"`
	BuildContext string                `json:"build_context"`
	Dockerfile   string                `json:"dockerfile"`
	HealthPath   string                `json:"health_path"`
}

// ServiceComposeCreate adds the project's service to a compose file, with
// the project's .env, and renders the added service.
func (h Handler) ServiceComposeCreate(w http.ResponseWriter, r *http.Request, logger *slog.Logger, form *composeForm) error {
	if !h.Compose.Enabled() {
		return http_errors.BadRequest(errors.New("ServiceComposeCreate: docker_compose_base_path is not set"))
	}
	svc, err := h.composeService(*form)
	if err != nil {
		return http_errors.BadRequest(pkgErrors.Wrap(err, "ServiceComposeCreate"))
	}
	path, err := h.Compose.Add(form.Type, svc)
	logger.Info("compose setup", "service", svc.Name, "file", path, "error", err)
	return blocks.ComposeSetupResult(svc, path, err).Render(r.Context(), w)
}

// composeService builds the service from the form; its .env is the one
// the Database step writes, when env_base_path is set.
func (h Handler) composeService(form composeForm) (compose.Service, error) {
	svc := compose.Service{
		Name:         strings.TrimSpace(form.Project),
		BuildContext: strings.TrimSuffix(strings.TrimSpace(form.BuildContext), "/"),
		Dockerfile:   cmp.Or(strings.TrimSpace(form.Dockerfile), "Dockerfile.prod"),
		Port:         int(form.Port),
		HealthPath:   strings.TrimSpace(form.HealthPath),
	}
	if h.Envs.Enabled() {
		envPath, err := h.Envs.Path(svc.Name)
		if err != nil {
			return compose.Service{}, err
		}
		svc.EnvFile = envPath
	}
	return svc, svc.Validate()
}
