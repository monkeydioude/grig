package os

import (
	"context"
	"log/slog"
)

func DaemonReload(
	_ context.Context,
	logger *slog.Logger,
) error {
	logger.Info("nothing to reload on darwin")
	return nil
}

func ServiceRestart(
	_ context.Context,
	name string,
	logger *slog.Logger,
) error {
	logger.Info("nothing to restart on darwin", "service", name)
	return nil
}
