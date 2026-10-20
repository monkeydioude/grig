package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recorded struct {
	path string
	auth string
	body map[string]any
}

// fakeAPI answers every request with status and answer, and records it.
func fakeAPI(t *testing.T, status int, answer string) (*httptest.Server, *recorded) {
	t.Helper()
	rec := &recorded{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.path = r.URL.Path
		rec.auth = r.Header.Get("Authorization")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&rec.body))
		w.WriteHeader(status)
		w.Write([]byte(answer))
	}))
	t.Cleanup(srv.Close)
	return srv, rec
}

func testGitHub(url string) restForge {
	f := NewGitHub("tok", http.DefaultClient).(restForge)
	f.apiURL = url
	return f
}

func TestCreateRepoInOwnAccount(t *testing.T) {
	srv, rec := fakeAPI(t, 201, `{"full_name":"me/app","html_url":"https://x/me/app","ssh_url":"git@x:me/app.git"}`)

	repo, err := testGitHub(srv.URL).CreateRepo(context.Background(), RepoRequest{Name: "app", Private: true})

	require.NoError(t, err)
	assert.Equal(t, Repo{FullName: "me/app", HTMLURL: "https://x/me/app", SSHURL: "git@x:me/app.git"}, repo)
	assert.Equal(t, "/user/repos", rec.path)
	assert.Equal(t, "Bearer tok", rec.auth)
	assert.Equal(t, map[string]any{"name": "app", "private": true}, rec.body)
}

func TestCreateRepoInOrganization(t *testing.T) {
	srv, rec := fakeAPI(t, 201, `{"full_name":"calq/app"}`)
	f := NewForgejo(srv.URL, "tok", http.DefaultClient)

	_, err := f.CreateRepo(context.Background(), RepoRequest{Owner: "calq", Name: "app"})

	require.NoError(t, err)
	assert.Equal(t, "/api/v1/orgs/calq/repos", rec.path)
	assert.Equal(t, "token tok", rec.auth)
}

func TestCreateRepoReportsTheForgeMessage(t *testing.T) {
	srv, _ := fakeAPI(t, 422, `{"message":"Repository creation failed.","errors":[{"message":"name already exists on this account"}]}`)

	_, err := testGitHub(srv.URL).CreateRepo(context.Background(), RepoRequest{Name: "app"})

	var apiErr APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, 422, apiErr.Status)
	assert.Equal(t, "Repository creation failed., name already exists on this account", apiErr.Message)
}

func TestAddWebhookBodies(t *testing.T) {
	hook := Webhook{URL: "https://h/josuke/github", Secret: "s", Events: []string{"push"}}
	config := map[string]any{"url": hook.URL, "content_type": "json", "secret": "s"}

	srv, rec := fakeAPI(t, 201, `{}`)
	require.NoError(t, testGitHub(srv.URL).AddWebhook(context.Background(), "me/app", hook))
	assert.Equal(t, "/repos/me/app/hooks", rec.path)
	assert.Equal(t, "web", rec.body["name"])
	assert.Equal(t, config, rec.body["config"])

	srv, rec = fakeAPI(t, 201, `{}`)
	require.NoError(t, NewForgejo(srv.URL, "tok", http.DefaultClient).AddWebhook(context.Background(), "calq/app", hook))
	assert.Equal(t, "/api/v1/repos/calq/app/hooks", rec.path)
	assert.Equal(t, "forgejo", rec.body["type"])
	assert.Equal(t, config, rec.body["config"])
}

func TestFromEnvOnlyRegistersForgesWithAToken(t *testing.T) {
	env := map[string]string{"CODEBERG_TOKEN": "tok"}
	reg := FromEnv(func(key string) string { return env[key] }, http.DefaultClient)

	assert.Equal(t, []string{CodebergName}, reg.Names())
	assert.Equal(t, "https://codeberg.org/api/v1", reg[CodebergName].(restForge).apiURL)
}
