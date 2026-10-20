// Package restart restarts the systemd unit reading a config file, so a
// config saved from grig is applied: Josuke and Capybara read theirs at start.
package restart

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/model"
	"net/http"
)

var (
	ErrNoUnit     = errors.New("no unit to restart for this config: see restart_on_save")
	ErrJosukeBusy = errors.New("josuke is running a deployment, restart it once done")
)

// RestartFunc restarts a systemd unit, e.g. pkg/os.ServiceRestart.
type RestartFunc func(ctx context.Context, unit string, logger *slog.Logger) error

type Restarter struct {
	// Units maps a config file path to the unit reading it.
	Units   map[string]string
	Restart RestartFunc
	// Client asks Josuke for its jobs before restarting it.
	Client *http.Client
}

// Can tells whether the config at path has a unit to restart.
func (r Restarter) Can(path string) bool {
	_, ok := r.Units[path]
	return ok
}

// RestartFor restarts the unit reading the config at path.
func (r Restarter) RestartFor(ctx context.Context, path string, logger *slog.Logger) error {
	unit, ok := r.Units[path]
	if !ok {
		return ErrNoUnit
	}
	return r.Restart(ctx, unit, logger)
}

// RestartJosuke restarts Josuke unless it is running a job, which the
// restart would kill. Waiting jobs are kept in Josuke's queue file.
func (r Restarter) RestartJosuke(ctx context.Context, jk model.Josuke, logger *slog.Logger) error {
	if !r.Can(jk.Path) {
		return ErrNoUnit
	}
	jobs, err := NewJobsClient(jk, r.Client)
	if err != nil {
		return err
	}
	busy, err := jobs.Busy(ctx)
	if err != nil {
		return fmt.Errorf("could not check josuke jobs: %w", err)
	}
	if busy {
		return ErrJosukeBusy
	}
	return r.RestartFor(ctx, jk.Path, logger)
}
