package restart

import (
	"context"
	"log/slog"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/trans_types"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeJosuke serves jobs on GET /jobs, for the token "tok" only.
func fakeJosuke(t *testing.T, jobs string) model.Josuke {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jobs" || r.Header.Get("Authorization") != "Bearer tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(jobs))
	}))
	t.Cleanup(srv.Close)
	host, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
	portNum, _ := strconv.Atoi(port)
	return model.Josuke{
		Host:  host,
		Port:  trans_types.StringInt(portNum),
		Path:  "/conf/josuke.yaml",
		Extra: map[string]any{"jobs_api": map[string]any{"token": "tok"}},
	}
}

// recorder is a Restarter recording the units it restarts.
func recorder(restarted *[]string) Restarter {
	return Restarter{
		Units: map[string]string{"/conf/josuke.yaml": "josuke.service", "/conf/capybara.yaml": "capybara.service"},
		Restart: func(_ context.Context, unit string, _ *slog.Logger) error {
			*restarted = append(*restarted, unit)
			return nil
		},
		Client: http.DefaultClient,
	}
}

func TestRestartForUsesTheMappedUnit(t *testing.T) {
	restarted := []string{}
	r := recorder(&restarted)

	require.NoError(t, r.RestartFor(context.Background(), "/conf/capybara.yaml", slog.Default()))
	assert.ErrorIs(t, r.RestartFor(context.Background(), "/conf/garden.yaml", slog.Default()), ErrNoUnit)
	assert.Equal(t, []string{"capybara.service"}, restarted)
}

func TestRestartJosukeWhenOnlyWaitingJobs(t *testing.T) {
	restarted := []string{}
	jk := fakeJosuke(t, `[{"id":2,"status":"waiting"}]`)

	require.NoError(t, recorder(&restarted).RestartJosuke(context.Background(), jk, slog.Default()))
	assert.Equal(t, []string{"josuke.service"}, restarted)
}

func TestRestartJosukeRefusesWhileAJobRuns(t *testing.T) {
	for _, status := range []string{"running", "stopping"} {
		restarted := []string{}
		jk := fakeJosuke(t, `[{"id":1,"status":"`+status+`"}]`)

		err := recorder(&restarted).RestartJosuke(context.Background(), jk, slog.Default())

		assert.ErrorIs(t, err, ErrJosukeBusy, status)
		assert.Empty(t, restarted, status)
	}
}

func TestRestartJosukeNeedsTheJobsAPI(t *testing.T) {
	restarted := []string{}
	jk := fakeJosuke(t, `[]`)
	jk.Extra = nil

	err := recorder(&restarted).RestartJosuke(context.Background(), jk, slog.Default())

	assert.ErrorIs(t, err, ErrNoJobsAPI)
	assert.Empty(t, restarted)
}
