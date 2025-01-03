package pages

import (
	"fmt"
	"log"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/server"
)

type Capybara struct {
	Titl     string
	Data     *model.Capybara
	FilePath string
}

func CapybaraList(config *server.ServerConfig) Capybara {
	p := Capybara{
		Titl: "Create a Capybara config",
		Data: &model.Capybara{
			Services: make([]model.ServiceDefinition, 1),
		},
	}

	if config == nil || config.CapybaraConfigPath == "" {
		return p
	}
	cp, err := fs.UnmarshalFromPath[model.Capybara](config.CapybaraConfigPath)
	if err != nil {
		log.Printf("[ERR ] pages.Capybaralist: %s", err)
		return p
	}
	p.Data = &cp
	return p
}

func (c Capybara) Title() string {
	return c.Titl
}

func (c Capybara) GetServiceInputName(it int, key string) string {
	return fmt.Sprintf("services[%d][%s]", it, key)
}

func (c Capybara) GetId(it int, key string) string {
	return fmt.Sprintf("services-%d-%s", it, key)
}
