package os

import (
	"errors"
	"log"
	"runtime"
)

var ErrUnknownOS = errors.New("unknown OS")

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
	log.Fatalf("%s: %s\n", runtime.GOOS, ErrUnknownOS)
	return Unknown
}
