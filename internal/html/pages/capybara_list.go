package pages

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/server/config"
)

type Capybara struct {
	Titl     string
	Data     *model.Capybara
	FilePath string
}

func CapybaraList(config *config.ServerConfig) Capybara {
	p := Capybara{
		Titl: "Create a Capybara config",
		Data: &model.Capybara{
			Services: make([]model.ServiceDefinition, 1),
		},
	}

	if config == nil || config.CapybaraConfigPath == "" {
		return p
	}
	cp, err := file.UnmarshalFromPath[model.Capybara](config.CapybaraConfigPath)
	if err != nil {
		slog.Error("pages.Capybaralist", "error", err)
		return p
	}
	p.Data = &cp
	return p
}

func (c Capybara) Title() string {
	return c.Titl
}

func GetServiceInputName(it int, key string) string {
	return fmt.Sprintf("services[%d][%s]", it, key)
}

func (c Capybara) GetId(it int, key string) string {
	return fmt.Sprintf("services-%d-%s", it, key)
}
