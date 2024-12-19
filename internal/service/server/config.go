package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/app_service"
	"monkeydioude/grig/internal/service/fs"
	"monkeydioude/grig/internal/tiger/assert"
	"os"
	"path/filepath"
)

// ServerConfig holds the app config and is also
// a factory for generating some model's entities
type ServerConfig struct {
	ServerconfigPath   string `json:"-"`
	AppsServicesDir    string `json:"app_services_dir"`
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
		// file, err := os.Create(appConfigPath)
		// assert.NoError(err)
		// assert.NotNil(file)
		// _, err = file.WriteString("{}")
		// assert.NoError(err)
		// assert.NoError(file.Close())
		configRaw = []byte("{}")
	} else {
		assert.NoError(err)
	}
	return configRaw
}

// NewServerConfigFromPath tries to parse a file located at `appConfigPath`,
// holding the server config.
func NewServerConfigFromPath(appConfigPath string) ServerConfig {
	// sanitize file path (in case of change dir etc...)
	appConfigPath, err := filepath.Abs(appConfigPath)
	assert.NoError(err)
	// read the file and return a bag of raw bytes
	configRaw := readConfigFile(appConfigPath)
	// try to unmarshal config raw bytes
	config := unmarshalConfig(configRaw)
	config.ServerconfigPath = appConfigPath
	return config
}

func (sc ServerConfig) ProcessAppsServicesDir() (*fs.Dir[model.Service], error) {
	if sc.AppsServicesDir == "" {
		return nil, nil
	}
	res, err := fs.NewDirFromPath(sc.AppsServicesDir, app_service.NewServiceFromPath)
	if err != nil {
		return nil, fmt.Errorf("server.ServerConfig.ProcessAppsServicesDir: %w", err)
	}
	return &res, nil
}

func (sc ServerConfig) ProcessJosuke() *model.Josuke {
	if sc.JosukeConfigPath == "" {
		return nil
	}
	return nil
}

func (sc ServerConfig) ProcessCapybara() *fs.Dir[model.Service] {
	if sc.CapybaraConfigPath == "" {
		return nil
	}
	return nil
}
