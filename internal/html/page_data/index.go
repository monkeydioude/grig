package page_data

import (
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/server"
)

type IndexPage struct {
	CapybaraConfig *model.Capybara
	JosukeConfig   *model.Josuke
	Services       *fs.Dir[model.Service]
}

func Index(
	config *server.ServerConfig,
) IndexPage {
	p := IndexPage{}
	if config == nil {
		return p
	}

	return p
}

func (IndexPage) Title() string {
	return "Blitz Grig"
}
