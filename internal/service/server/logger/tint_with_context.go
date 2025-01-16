package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func SlogTintWithContext(ctxKeys ...string) slog.Handler {
	return SlogWithContext(tint.NewHandler(os.Stdout, nil), ctxKeys...)
}
