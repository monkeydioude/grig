package repo_setup

import (
	"context"
	"errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/forge"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeForge struct {
	repoErr  error
	hookErr  error
	hookRepo string
	hook     forge.Webhook
}

func (f *fakeForge) CreateRepo(_ context.Context, req forge.RepoRequest) (forge.Repo, error) {
	return forge.Repo{FullName: "me/" + req.Name}, f.repoErr
}

func (f *fakeForge) AddWebhook(_ context.Context, fullName string, hook forge.Webhook) error {
	f.hookRepo, f.hook = fullName, hook
	return f.hookErr
}

func TestFormRepoRequest(t *testing.T) {
	reg := forge.Registry{"github": &fakeForge{}}

	req, err := Form{Forge: "github", Owner: " calq ", Name: "my-app", Visibility: "public"}.RepoRequest(reg)
	require.NoError(t, err)
	assert.Equal(t, forge.RepoRequest{Owner: "calq", Name: "my-app", Private: false}, req)

	req, err = Form{Forge: "github", Name: "app"}.RepoRequest(reg)
	require.NoError(t, err)
	assert.True(t, req.Private, "private unless public is asked")

	_, err = Form{Forge: "gitlab", Name: "app"}.RepoRequest(reg)
	assert.ErrorIs(t, err, ErrUnknownForge)
	_, err = Form{Forge: "github", Name: "../app"}.RepoRequest(reg)
	assert.ErrorIs(t, err, ErrInvalidName)
	_, err = Form{Forge: "github", Owner: "a/b", Name: "app"}.RepoRequest(reg)
	assert.ErrorIs(t, err, ErrInvalidOwner)
}

func TestJosukeWebhook(t *testing.T) {
	jk := model.Josuke{Hook: []model.Hook{
		{Type: "gogs", Path: "/josuke/gogs", Secret: "a"},
		{Type: "github", Path: "/josuke/github", Secret: "b"},
	}}

	hook, err := JosukeWebhook(jk, "https://example.com/")
	require.NoError(t, err)
	assert.Equal(t, forge.Webhook{URL: "https://example.com/josuke/github", Secret: "b", Events: []string{"push"}}, hook)

	_, err = JosukeWebhook(jk, "")
	assert.ErrorIs(t, err, ErrNoWebhookBaseURL)
	_, err = JosukeWebhook(model.Josuke{}, "https://example.com")
	assert.ErrorIs(t, err, ErrNoJosukeHook)
}

func TestRunAddsTheWebhookToTheCreatedRepo(t *testing.T) {
	f := &fakeForge{}
	hook := forge.Webhook{URL: "https://h/josuke/github"}

	res := Run(context.Background(), f, forge.RepoRequest{Name: "app"}, hook, nil)

	assert.NoError(t, res.RepoErr)
	assert.True(t, res.HookDone)
	assert.Equal(t, "me/app", f.hookRepo)
	assert.Equal(t, hook, f.hook)
}

func TestRunStopsWhenTheRepoFails(t *testing.T) {
	f := &fakeForge{repoErr: errors.New("exists")}

	res := Run(context.Background(), f, forge.RepoRequest{Name: "app"}, forge.Webhook{}, nil)

	assert.Error(t, res.RepoErr)
	assert.False(t, res.HookDone)
	assert.Empty(t, f.hookRepo)
}

func TestRunKeepsTheRepoWhenTheWebhookCannotBeBuilt(t *testing.T) {
	f := &fakeForge{}

	res := Run(context.Background(), f, forge.RepoRequest{Name: "app"}, forge.Webhook{}, ErrNoJosukeHook)

	assert.NoError(t, res.RepoErr)
	assert.ErrorIs(t, res.HookErr, ErrNoJosukeHook)
	assert.Empty(t, f.hookRepo)
}
