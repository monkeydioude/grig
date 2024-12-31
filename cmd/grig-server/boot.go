package main

import (
	"flag"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/os"
	"monkeydioude/grig/internal/service/server"
	"monkeydioude/grig/internal/tiger/assert"
	nativeOs "os"
)

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
		assert.NoError(fs.CreateAndWriteFile(config.CapybaraConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	if config.JosukeConfigPath == "" {
		config.JosukeConfigPath = fs.AppendToThisFileDirectory(consts.DEFAULT_JOSUKE_FILENAME, config.ServerConfigPath)
		assert.NoError(fs.CreateAndWriteFile(config.JosukeConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	assert.NoError(config.Save())
	layout := server.Layout{
		OS:           os.FindoutOS(),
		ServerConfig: config,
	}
	// fmt.Printf("appServices: %+v\njosuke: %+v\ncapybara: %+v\n", layout.AppsServices, layout.JosukeConfig, layout.CapybaraConfig)
	return &layout
}
