package os

import (
	"context"
	"log/slog"
	"monkeydioude/grig/pkg/errors"

	"github.com/coreos/go-systemd/v22/dbus"
)

func DaemonReload(
	ctx context.Context,
	logger *slog.Logger,
) error {
	logger.Info("connecting to dbus")
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return errors.Wrap(err, "DaemonReload::NewSystemdConnectionContext")
	}
	defer conn.Close()
	logger.Info("trying to reload daemons")
	err = conn.ReloadContext(ctx)
	if err != nil {
		return errors.Wrap(err, "DaemonReload::ReloadContext")
	}

	return nil
}

func ServiceRestart(
	ctx context.Context,
	name string,
	logger *slog.Logger,
) error {
	logger.Info("connecting to dbus")
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return errors.Wrap(err, "ServiceRestart::NewSystemdConnectionContext")
	}
	defer conn.Close()
	logger.Info("trying to restart service", "service", name)
	jobID, err := conn.RestartUnitContext(ctx, name, "replace", nil)
	if err != nil {
		return errors.Wrap(err, "ServiceRestart::RestartUnitContext")
	}
	logger.Info("service restarted", "job_id", jobID)
	return nil
}
