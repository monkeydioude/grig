package pages

import (
	"errors"
	"fmt"
	"log/slog"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/server/config"
)

type IndexPage struct {
	Capybara *model.Capybara
	Josuke   *model.Josuke
	Err      error
}

func Index(
	config *config.ServerConfig,
) IndexPage {
	p := IndexPage{}
	if config == nil {
		p.Err = errors.Join(p.Err, fmt.Errorf("pages.Index: config: %w", customErrors.ErrNilPointer))
		return p
	}
	// capybara
	cp, err := file.UnmarshalFromPath[model.Capybara](config.CapybaraConfigPath)
	if err != nil {
		slog.Error("pages.Index", "error", err)
		p.Err = errors.Join(p.Err, err)
	} else {
		p.Capybara = &cp
	}

	// josuke
	jk, err := file.UnmarshalFromPath[model.Josuke](config.JosukeConfigPath)
	if err != nil {
		slog.Error("pages.Index", "error", err)
		p.Err = errors.Join(p.Err, err)
	} else {
		p.Josuke = &jk
	}
	return p
}

func (IndexPage) Title() string {
	return "Blitz Grig"
}
