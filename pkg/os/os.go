package os

import (
	"log"
	"runtime"

	"monkeydioude/grig/internal/errors"
)

type OS string

const (
	Linux   OS = "linux"
	MacOS   OS = "darwin"
	Unknown OS = "unknown"
)

func FindoutOS() OS {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return MacOS
	}
	log.Fatalf("%s: %s\n", runtime.GOOS, errors.ErrUnknownOS)
	return Unknown
}
