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
