package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// restForge talks to a GitHub-like REST API. GitHub and Forgejo share
// the same routes and repository fields; they differ in the API prefix,
// the auth header and the webhook body.
type restForge struct {
	apiURL   string
	auth     string
	headers  map[string]string
	hookBody func(Webhook) any
	client   *http.Client
}

func NewGitHub(token string, client *http.Client) Forge {
	return restForge{
		apiURL: "https://api.github.com",
		auth:   "Bearer " + token,
		headers: map[string]string{
			"Accept":               "application/vnd.github+json",
			"X-GitHub-Api-Version": "2022-11-28",
		},
		hookBody: gitHubHookBody,
		client:   client,
	}
}

func NewForgejo(baseURL, token string, client *http.Client) Forge {
	return restForge{
		apiURL:   strings.TrimRight(baseURL, "/") + "/api/v1",
		auth:     "token " + token,
		headers:  map[string]string{"Accept": "application/json"},
		hookBody: forgejoHookBody,
		client:   client,
	}
}

func gitHubHookBody(hook Webhook) any {
	return map[string]any{
		"name":   "web",
		"active": true,
		"events": hook.Events,
		"config": hookConfig(hook),
	}
}

func forgejoHookBody(hook Webhook) any {
	return map[string]any{
		"type":   "forgejo",
		"active": true,
		"events": hook.Events,
		"config": hookConfig(hook),
	}
}

func hookConfig(hook Webhook) map[string]string {
	return map[string]string{
		"url":          hook.URL,
		"content_type": "json",
		"secret":       hook.Secret,
	}
}

func (f restForge) CreateRepo(ctx context.Context, req RepoRequest) (Repo, error) {
	path := "/user/repos"
	if req.Owner != "" {
		path = "/orgs/" + req.Owner + "/repos"
	}
	body := map[string]any{"name": req.Name, "private": req.Private}
	repo := Repo{}
	if err := f.post(ctx, path, body, &repo); err != nil {
		return Repo{}, fmt.Errorf("forge.CreateRepo: %w", err)
	}
	return repo, nil
}

func (f restForge) AddWebhook(ctx context.Context, fullName string, hook Webhook) error {
	if err := f.post(ctx, "/repos/"+fullName+"/hooks", f.hookBody(hook), nil); err != nil {
		return fmt.Errorf("forge.AddWebhook: %w", err)
	}
	return nil
}

// post sends body as JSON and decodes the answer into out, when not nil.
func (f restForge) post(ctx context.Context, path string, body, out any) error {
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.apiURL+path, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	f.setHeaders(req)
	res, err := f.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return newAPIError(res)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (f restForge) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", f.auth)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range f.headers {
		req.Header.Set(key, value)
	}
}

// APIError is a non-2xx answer of a forge.
type APIError struct {
	Status  int
	Message string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

// newAPIError reads the message both GitHub and Forgejo put in their
// error bodies, plus GitHub's detailed errors ("name already exists...").
func newAPIError(res *http.Response) APIError {
	body, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	parsed := struct {
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}{}
	if json.Unmarshal(body, &parsed) != nil || parsed.Message == "" {
		return APIError{Status: res.StatusCode, Message: http.StatusText(res.StatusCode)}
	}
	msg := parsed.Message
	for _, detail := range parsed.Errors {
		if detail.Message != "" {
			msg += ", " + detail.Message
		}
	}
	return APIError{Status: res.StatusCode, Message: msg}
}
