package pages

import (
	"errors"
	"fmt"
	"log"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/fs"
)

type IndexPage struct {
	Capybara    *model.Capybara
	Josuke      *model.Josuke
	ServicesLen int
	Err         error
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
	cp, err := fs.UnmarshalFromPath[model.Capybara](config.CapybaraConfigPath)
	if err != nil {
		log.Printf("[ERR ] pages.Index: %q", err)
		p.Err = errors.Join(p.Err, err)
	} else {
		p.Capybara = &cp
	}
	p.ServicesLen = len(cp.Services)

	// josuke
	jk, err := fs.UnmarshalFromPath[model.Josuke](config.JosukeConfigPath)
	if err != nil {
		log.Printf("[ERR ] pages.Index: %q", err)
		p.Err = errors.Join(p.Err, err)
	} else {
		p.Josuke = &jk
	}
	return p
}

func (IndexPage) Title() string {
	return "Blitz Grig"
}
