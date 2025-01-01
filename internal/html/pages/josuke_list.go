package pages

import (
	"log"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/server"
)

type Josuke struct {
	Titl     string
	Data     *model.Josuke
	FilePath string
}

func JosukeList(config *server.ServerConfig) Josuke {
	p := Josuke{
		Titl: "Create a Josuke config",
		Data: &model.Josuke{},
	}

	if config == nil || config.JosukeConfigPath == "" {
		return p
	}
	jk, err := fs.UnmarshalFromPath[model.Josuke](config.JosukeConfigPath)
	if err != nil {
		log.Printf("[ERR ] pages.JosukeList: %s", err)
		return p
	}
	p.Data = &jk
	return p
}

func (c Josuke) Title() string {
	return c.Titl
}
