package main

import (
	"flag"
	"log/slog"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/server/config"
	"monkeydioude/grig/internal/service/server/logger"
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

func setLogger() {
	logger := slog.New(logger.SlogTintWithContext(consts.X_REQUEST_ID_LABEL))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func boot() *server.Layout[config.ServerConfig] {
	setLogger()
	mainConfigPath := parseFlags()
	conf := config.NewServerConfigFromPath(mainConfigPath)
	if conf.CapybaraConfigPath == "" {
		conf.CapybaraConfigPath = file.AppendToThisFileDirectory(consts.DEFAULT_CAPYBARA_FILENAME, conf.ServerConfigPath)
		assert.NoError(file.CreateAndWriteFile(conf.CapybaraConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	if conf.JosukeConfigPath == "" {
		conf.JosukeConfigPath = file.AppendToThisFileDirectory(consts.DEFAULT_JOSUKE_FILENAME, conf.ServerConfigPath)
		assert.NoError(file.CreateAndWriteFile(conf.JosukeConfigPath, []byte("{}"), nativeOs.ModePerm))
	}
	assert.NoError(conf.Save())
	layout := server.Layout[config.ServerConfig]{
		OS:           os.FindoutOS(),
		ServerConfig: conf,
	}
	return &layout
}
