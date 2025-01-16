package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/pkg/tiger/assert"
	"os"
	"path/filepath"
)

// ServerConfig holds the app config and is also
// a factory for generating some model's entities
type ServerConfig struct {
	ServerConfigPath   string   `json:"-"`
	AppsServicesPaths  []string `json:"services_paths"`
	JosukeConfigPath   string   `json:"josuke_config_path"`
	CapybaraConfigPath string   `json:"capybara_config_path"`
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

func (sc ServerConfig) Save() error {
	data, err := json.Marshal(&sc)
	if err != nil {
		return fmt.Errorf("ServerConfig.Save(): %w", err)
	}
	if err := file.CreateAndWriteFile(sc.ServerConfigPath, data, os.ModePerm); err != nil {
		return fmt.Errorf("ServerConfig.Save(): %w", err)
	}
	return nil
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
