package main

import (
	"flag"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/pkg/fs"
	"monkeydioude/grig/pkg/os"
	"monkeydioude/grig/pkg/server"
	"monkeydioude/grig/pkg/tiger/assert"
	nativeOs "os"
)

func parseFlags() string {
	mainConfigPath := flag.String("c", "grig_server.config.json", "-c <path to config>")
	flag.Parse()

	return *mainConfigPath
}

func boot() *server.Layout[config.ServerConfig] {
	mainConfigPath := parseFlags()
	conf := config.NewServerConfigFromPath(mainConfigPath)
	if conf.CapybaraConfigPath == "" {
		conf.CapybaraConfigPath = fs.AppendToThisFileDirectory(consts.DEFAULT_CAPYBARA_FILENAME, conf.ServerConfigPath)
		assert.NoError(fs.CreateAndWriteFile(conf.CapybaraConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	if conf.JosukeConfigPath == "" {
		conf.JosukeConfigPath = fs.AppendToThisFileDirectory(consts.DEFAULT_JOSUKE_FILENAME, conf.ServerConfigPath)
		assert.NoError(fs.CreateAndWriteFile(conf.JosukeConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	assert.NoError(conf.Save())
	layout := server.Layout[config.ServerConfig]{
		OS:           os.FindoutOS(),
		ServerConfig: conf,
	}
	return &layout
}
