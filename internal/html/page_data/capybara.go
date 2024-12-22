package page_data

import (
	"fmt"
	"log"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/server"
)

type CapybaraData struct {
	Titl     string
	Data     *model.Capybara
	FilePath string
}

func CapybaraCreate(
	config *server.ServerConfig,
) CapybaraData {
	p := CapybaraData{
		Titl: "Create a Capybara config",
	}
	if config == nil {
		return p
	}
	p.Data = &model.Capybara{}
	return p
}

func CapybaraList(config *server.ServerConfig) CapybaraData {
	p := CapybaraData{
		Titl: "Create a Capybara config",
		Data: &model.Capybara{},
	}
	if config == nil || config.CapybaraConfigPath == "" {
		return p
	}

	data, err := fs.UnmarshalFromPath[model.Capybara](config.CapybaraConfigPath)
	if err != nil {
		log.Printf("[ERR ] pages.Capybaralist: %s", err)
		return p
	}

	p.Data = &data
	return p
}

func (c CapybaraData) Title() string {
	return c.Titl
}

func (c CapybaraData) GetServiceInputName(it int, key string) string {
	return fmt.Sprintf("services[%d][%s]", it, key)
}
