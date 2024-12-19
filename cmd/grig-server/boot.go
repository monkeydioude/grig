package main

import (
	"flag"
	"fmt"

	"monkeydioude/grig/internal/service/os"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/tiger/assert"
)

type AppConfigPath = string

func parseFlags() AppConfigPath {
	appConfig := flag.String("c", "grig_server.config.json", "-c <path to config>")
	flag.Parse()

	return *appConfig
}

func boot() server.Layout {
	config := server.NewServerConfigFromPath(parseFlags())

	appServices, err := config.ProcessAppsServicesDir()
	assert.NoError(err)
	fmt.Printf("appServices, %+v\n", appServices)
	layout := server.Layout{
		OS:           os.FindoutOS(),
		ServerConfig: config,
		AppsServices: appServices,
	}
	return layout
}
