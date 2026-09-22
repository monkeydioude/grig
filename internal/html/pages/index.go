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

// CapybaraCard is one capybara config as shown on the index: the ref
// identifying it, plus either its parsed content or the load error.
type CapybaraCard struct {
	Ref  config.CapybaraRef
	Data *model.Capybara
	Err  error
}

type IndexPage struct {
	Capybaras []CapybaraCard
	Josuke    *model.Josuke
	Err       error
}

func loadCapybaraCard(ref config.CapybaraRef) CapybaraCard {
	cp, err := file.UnmarshalFromPath[model.Capybara](ref.Path)
	if err != nil {
		slog.Error("pages.Index", "error", err, "path", ref.Path)
		return CapybaraCard{Ref: ref, Err: err}
	}
	return CapybaraCard{Ref: ref, Data: &cp}
}

func Index(
	config *config.ServerConfig,
) IndexPage {
	p := IndexPage{}
	if config == nil {
		p.Err = errors.Join(p.Err, fmt.Errorf("pages.Index: config: %w", customErrors.ErrNilPointer))
		return p
	}
	// capybara: a broken config is reported on its own card, so the
	// other configs still render
	for _, ref := range config.CapybaraRefs() {
		p.Capybaras = append(p.Capybaras, loadCapybaraCard(ref))
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
	return "Grig"
}
