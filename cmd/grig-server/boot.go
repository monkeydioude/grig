package main

import (
	"flag"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/html/element"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/os"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/tiger/assert"
)

func MainNavigation() element.Nav {
	return element.Nav{
		Links: []element.Link{
			{
				Href:   "/",
				Text:   element.Text("Home"),
				Target: element.Self,
			},
			{
				Href:   "/capybara",
				Text:   element.Text("Capybara"),
				Target: element.Self,
			},
		},
	}
}

func parseFlags() string {
	mainConfigPath := flag.String("c", "grig_server.config.json", "-c <path to config>")
	flag.Parse()

	return *mainConfigPath
}

func boot() *server.Layout {
	mainConfigPath := parseFlags()
	config := server.NewServerConfigFromPath(mainConfigPath)
	if config.CapybaraConfigPath == "" {
		config.CapybaraConfigPath = fs.AppendToThisFileDirectory(consts.DEFAULT_CAPYBARA_FILENAME, config.ServerConfigPath)
	}
	assert.NoError(config.Save())
	layout := server.Layout{
		OS:           os.FindoutOS(),
		Navigation:   MainNavigation(),
		ServerConfig: config,
	}
	// fmt.Printf("appServices: %+v\njosuke: %+v\ncapybara: %+v\n", layout.AppsServices, layout.JosukeConfig, layout.CapybaraConfig)
	return &layout
}
