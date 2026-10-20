package v1

import (
	"errors"
	"log/slog"
	"monkeydioude/grig/internal/html/blocks"
	"monkeydioude/grig/internal/html/layouts"
	"monkeydioude/grig/internal/html/pages"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/db_schema"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/forge"
	"monkeydioude/grig/internal/service/repo_setup"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/html/elements"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
)

func (h Handler) ServiceNew(w http.ResponseWriter, r *http.Request, _ *slog.Logger, nav elements.Nav) error {
	hook, hookErr := h.josukeWebhook()
	page := pages.NewServiceNew(h.Forges.Names(), hook.URL, hookErr)
	if h.Schemas.Enabled() {
		password, err := db_schema.GeneratePassword()
		if err != nil {
			return err
		}
		page = page.WithDatabase(db_schema.Redacted(h.Schemas.URL), password)
	}
	return layouts.Main(nav, page).Render(r.Context(), w)
}

// ServiceSchemaCreate creates the database schema and user of a new
// service, and renders the DB_URL the service connects with.
func (h Handler) ServiceSchemaCreate(w http.ResponseWriter, r *http.Request, logger *slog.Logger, req *db_schema.Request) error {
	if !h.Schemas.Enabled() {
		return http_errors.BadRequest(errors.New("ServiceSchemaCreate: DB_URL is not set"))
	}
	if err := req.Validate(); err != nil {
		return http_errors.BadRequest(pkgErrors.Wrap(err, "ServiceSchemaCreate"))
	}
	res, err := h.Schemas.Create(r.Context(), *req)
	logger.Info("schema setup", "schema", req.Schema, "user", req.User, "error", err)
	return blocks.SchemaSetupResult(res, err).Render(r.Context(), w)
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
