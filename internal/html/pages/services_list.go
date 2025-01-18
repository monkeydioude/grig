package pages

import (
	"log/slog"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/parser"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/html/elements/htmx"
)

type Services struct {
	Data        []model.Service
	FilePath    string
	HtmxFactory htmx.Factory
}

func ServicesList(config *config.ServerConfig, logger *slog.Logger) Services {
	p := Services{
		Data:        []model.Service{},
		HtmxFactory: htmx.NewJsonFactory(),
	}
	if config == nil || len(config.AppsServicesPaths) == 0 {
		return p
	}

	srvcs, err := parser.IniServicesParser(config.AppsServicesPaths)
	if err != nil {
		logger.Error("in ServicesList, IniServicesParser", slog.String("error", err.Error()))
	}
	p.Data = srvcs
	return p
}

func (Services) Title() string {
	return "Services"
}
