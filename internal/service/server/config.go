package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/service/parser"
	"monkeydioude/grig/internal/tiger/assert"
	"os"
	"path/filepath"
)

// ServerConfig holds the app config and is also
// a factory for generating some model's entities
type ServerConfig struct {
	ServerConfigPath   string `json:"-"`
	AppsServicesDir    string `json:"parsers_dir"`
	JosukeConfigPath   string `json:"josuke_config_path"`
	CapybaraConfigPath string `json:"capybara_config_path"`
}

func unmarshalConfig(configRaw []byte) ServerConfig {
	config := ServerConfig{}
	err := json.Unmarshal(configRaw, &config)
	assert.NoError(err)
	return config
}

func readConfigFile(appConfigPath string) []byte {
	configRaw, err := os.ReadFile(appConfigPath)
	// create the config file if does not exist
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(appConfigPath)
		assert.NoError(err)
		assert.NotNil(file)
		configRaw = []byte("{}")
		n, err := file.Write(configRaw)
		assert.NotEmpty(n)
		assert.NoError(err)
		assert.NoError(file.Close())
	} else {
		assert.NoError(err)
	}
	return configRaw
}

// NewServerConfigFromPath tries to parse a file located at `appConfigPath`,
// holding the server config.
func NewServerConfigFromPath(mainConfigPath string) ServerConfig {
	// sanitize file path (in case of change dir etc...)
	mainConfigPath, err := filepath.Abs(mainConfigPath)
	assert.NoError(err)
	// read the file and return a bag of raw bytes
	configRaw := readConfigFile(mainConfigPath)
	// try to unmarshal config raw bytes
	config := unmarshalConfig(configRaw)
	config.ServerConfigPath = mainConfigPath
	return config
}

func (sc ServerConfig) ProcessAppsServicesDir() (*fs.Dir[model.Service], error) {
	if sc.AppsServicesDir == "" {
		return nil, nil
	}
	sc.AppsServicesDir = fs.AppendToThisFileDirectory(sc.AppsServicesDir, sc.ServerConfigPath)
	res, err := fs.NewDirFromPathAndFileParser(sc.AppsServicesDir, parser.ServiceFileParser)
	if err != nil {
		return nil, fmt.Errorf("server.ServerConfig.ProcessAppsServicesDir: %w", err)
	}
	return &res, nil
}

func (sc ServerConfig) ProcessJosuke() (*model.Josuke, error) {
	if sc.JosukeConfigPath == "" {
		return nil, nil
	}
	sc.JosukeConfigPath = fs.AppendToThisFileDirectory(sc.JosukeConfigPath, sc.ServerConfigPath)
	jojo, err := fs.UnmarshalFromPath[model.Josuke](sc.JosukeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("server.ServerConfig.ProcessJosuke: %w", err)
	}
	return &jojo, nil
}

func (sc ServerConfig) ProcessCapybara() (*model.Capybara, error) {
	if sc.CapybaraConfigPath == "" {
		return nil, nil
	}
	sc.CapybaraConfigPath = fs.AppendToThisFileDirectory(sc.CapybaraConfigPath, sc.ServerConfigPath)
	capy, err := fs.UnmarshalFromPath[model.Capybara](sc.CapybaraConfigPath)
	if err != nil {
		return nil, fmt.Errorf("server.ServerConfig.ProcessJosuke: %w", err)
	}
	return &capy, nil
}
