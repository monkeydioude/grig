package server

import (
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/os"
)

type Layout struct {
	OS             os.OS
	AppsServices   *fs.Dir[model.Service]
	JosukeConfig   *model.Josuke
	CapybaraConfig *model.Capybara
	ServerConfig   ServerConfig
}
