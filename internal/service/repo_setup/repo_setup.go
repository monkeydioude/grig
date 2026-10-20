// Package repo_setup creates the repository of a new service on a forge,
// with the webhook that makes Josuke deploy it on push.
package repo_setup

import (
	"context"
	"errors"
	"fmt"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/forge"
	"regexp"
	"strings"
)

var (
	ErrUnknownForge     = errors.New("unknown forge")
	ErrInvalidName      = errors.New("invalid repository name")
	ErrInvalidOwner     = errors.New("invalid owner")
	ErrNoWebhookBaseURL = errors.New("webhook_base_url is not set in grig's config")
	ErrNoJosukeHook     = errors.New(`josuke config has no hook of type "github"`)
)

const josukeHookType = "github"

var validName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// Form is what the "+ Service" page posts.
type Form struct {
	Forge      string `json:"forge"`
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
}

// RepoRequest checks the form and turns it into a forge request.
func (f Form) RepoRequest(reg forge.Registry) (forge.RepoRequest, error) {
	if _, ok := reg[f.Forge]; !ok {
		return forge.RepoRequest{}, fmt.Errorf("%w: %s", ErrUnknownForge, f.Forge)
	}
	name, owner := strings.TrimSpace(f.Name), strings.TrimSpace(f.Owner)
	if !validName.MatchString(name) {
		return forge.RepoRequest{}, ErrInvalidName
	}
	if owner != "" && !validName.MatchString(owner) {
		return forge.RepoRequest{}, ErrInvalidOwner
	}
	return forge.RepoRequest{Owner: owner, Name: name, Private: f.Visibility != "public"}, nil
}

// JosukeWebhook builds the webhook pointing at Josuke's "github" hook,
// which also accepts Forgejo's deliveries.
func JosukeWebhook(jk model.Josuke, baseURL string) (forge.Webhook, error) {
	if baseURL == "" {
		return forge.Webhook{}, ErrNoWebhookBaseURL
	}
	for _, hook := range jk.Hook {
		if hook.Type == josukeHookType {
			return forge.Webhook{
				URL:    strings.TrimRight(baseURL, "/") + hook.Path,
				Secret: hook.Secret,
				Events: []string{"push"},
			}, nil
		}
	}
	return forge.Webhook{}, ErrNoJosukeHook
}

// Result reports each step. A step that did not run has a nil error and
// is told apart by the step before it failing.
type Result struct {
	Repo     forge.Repo
	RepoErr  error
	HookURL  string
	HookErr  error
	HookDone bool
}

// Run creates the repository, then its webhook. A webhook that cannot
// be built (hookErr) is reported without blocking the repository.
func Run(ctx context.Context, f forge.Forge, req forge.RepoRequest, hook forge.Webhook, hookErr error) Result {
	res := Result{HookURL: hook.URL}
	res.Repo, res.RepoErr = f.CreateRepo(ctx, req)
	if res.RepoErr != nil {
		return res
	}
	if hookErr != nil {
		res.HookErr = hookErr
		return res
	}
	res.HookErr = f.AddWebhook(ctx, res.Repo.FullName, hook)
	res.HookDone = res.HookErr == nil
	return res
}
