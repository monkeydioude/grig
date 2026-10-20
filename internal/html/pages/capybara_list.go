package pages

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/server/config"
)

type Capybara struct {
	Data *model.Capybara
	Ref  config.CapybaraRef
	Err  error
	// Restartable offers "Save & restart" (restart_on_save has its unit).
	Restartable bool
}

func CapybaraList(ref config.CapybaraRef, restartable bool) Capybara {
	p := Capybara{
		Data: &model.Capybara{
			Services: make([]model.ServiceDefinition, 1),
		},
		Ref:         ref,
		Restartable: restartable,
	}

	if ref.Path == "" {
		return p
	}
	cp, err := file.UnmarshalFromPath[model.Capybara](ref.Path)
	if err != nil {
		slog.Error("pages.CapybaraList", "error", err, "path", ref.Path)
		p.Err = err
		return p
	}
	p.Data = &cp
	return p
}

func (c Capybara) Title() string {
	return c.Ref.Name
}

func GetServiceInputName(it int, key string) string {
	return fmt.Sprintf("services[%d][%s]", it, key)
}

func (c Capybara) GetId(it int, key string) string {
	return fmt.Sprintf("services-%d-%s", it, key)
}
