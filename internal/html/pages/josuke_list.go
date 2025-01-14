package pages

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/internal/service/utils"
	pkgModel "monkeydioude/grig/pkg/model"
)

type Josuke struct {
	Titl     string
	Data     *model.Josuke
	FilePath string
}

func JosukeList(config *config.ServerConfig) Josuke {
	p := Josuke{
		Titl: "Create a Josuke config",
		Data: &model.Josuke{},
	}

	if config == nil || config.JosukeConfigPath == "" {
		return p
	}
	jk, err := file.UnmarshalFromPath[model.Josuke](config.JosukeConfigPath)
	if err != nil {
		slog.Error("pages.JosukeList", "error", err)
		return p
	}
	jk.FillBaseData()
	p.Data = &jk
	return p
}

func (jk Josuke) Title() string {
	return jk.Titl
}

func GetHookInputName(it int, key string) string {
	return fmt.Sprintf("hook[%d][%s]", it, key)
}

func GetInputNameWithKey(hp pkgModel.IndexBuilder, key string) string {
	return utils.GetInputName(hp, key)
}
