// Package forge creates repositories and their webhooks on git hosts
// (GitHub, Codeberg/Forgejo) through their REST APIs.
package forge

import (
	"context"
	"net/http"
	"slices"
)

// RepoRequest describes the repository to create.
type RepoRequest struct {
	// Owner is the organization to create the repository in.
	// Empty creates it in the account of the token.
	Owner   string
	Name    string
	Private bool
}

// Repo is a created repository.
type Repo struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
	SSHURL   string `json:"ssh_url"`
}

// Webhook is the hook to register on a repository.
type Webhook struct {
	URL    string
	Secret string
	Events []string
}

type Forge interface {
	CreateRepo(ctx context.Context, req RepoRequest) (Repo, error)
	// AddWebhook registers hook on the repository named "owner/name".
	AddWebhook(ctx context.Context, fullName string, hook Webhook) error
}

// Registry holds the forges grig can use, by name.
type Registry map[string]Forge

// Names lists the registered forges, sorted.
func (r Registry) Names() []string {
	names := make([]string, 0, len(r))
	for name := range r {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

const (
	GitHubName   = "github"
	CodebergName = "codeberg"

	defaultCodebergURL = "https://codeberg.org"
)

// FromEnv registers the forges whose token is set:
// GITHUB_TOKEN, and CODEBERG_TOKEN (CODEBERG_URL for another Forgejo host).
func FromEnv(getenv func(string) string, client *http.Client) Registry {
	reg := Registry{}
	if token := getenv("GITHUB_TOKEN"); token != "" {
		reg[GitHubName] = NewGitHub(token, client)
	}
	if token := getenv("CODEBERG_TOKEN"); token != "" {
		url := getenv("CODEBERG_URL")
		if url == "" {
			url = defaultCodebergURL
		}
		reg[CodebergName] = NewForgejo(url, token, client)
	}
	return reg
}
