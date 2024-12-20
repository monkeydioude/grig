package main

import (
	"flag"
	"fmt"
	"monkeydioude/grig/internal/service/os"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/tiger/assert"
)

func parseFlags() string {
	mainConfigPath := flag.String("c", "grig_server.config.json", "-c <path to config>")
	flag.Parse()

	return *mainConfigPath
}

func boot() server.Layout {
	mainConfigPath := parseFlags()
	config := server.NewServerConfigFromPath(mainConfigPath)
	layout := server.Layout{
		OS: os.FindoutOS(),
	}
	{
		appServices, err := config.ProcessAppsServicesDir()
		assert.NoError(err)
		layout.AppsServices = appServices
	}
	{
		josukeConfig, err := config.ProcessJosuke()
		assert.NoError(err)
		layout.JosukeConfig = josukeConfig
	}
	{
		capyConfig, err := config.ProcessCapybara()
		assert.NoError(err)
		layout.CapybaraConfig = capyConfig
	}
	fmt.Printf("appServices: %+v\njosuke: %+v\ncapybara: %+v\n", layout.AppsServices, layout.JosukeConfig, layout.CapybaraConfig)
	return layout
}
