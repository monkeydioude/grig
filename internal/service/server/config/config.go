package config

import (
	"errors"
	"fmt"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/services"
	"monkeydioude/grig/pkg/tiger/assert"
	"os"
	"path/filepath"

	"monkeydioude/grig/pkg/urls"

	"gopkg.in/yaml.v3"
)

// ServerConfig holds the app config and is also
// a factory for generating some model's entities
type ServerConfig struct {
	ServerConfigPath string `json:"-" yaml:"-"`
	// BasePath is the sub-path the app is served under, e.g. "/grig" when
	// a reverse proxy forwards :80/grig to :6969/grig without stripping
	// the prefix. Empty serves the app at the root.
	BasePath            string                   `json:"base_path" yaml:"base_path"`
	AppsServicesPaths   services.AppServicePaths `json:"services_paths" yaml:"services_paths"`
	JosukeConfigPath    string                   `json:"josuke_config_path" yaml:"josuke_config_path"`
	CapybaraConfigPaths []string                 `json:"capybara_config_paths" yaml:"capybara_config_paths"`
}

func unmarshalConfig(configRaw []byte) ServerConfig {
	config := ServerConfig{}
	err := yaml.Unmarshal(configRaw, &config)
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

// Prefixer builds links for the sub-path this app is served under.
func (sc ServerConfig) Prefixer() urls.Prefixer {
	return urls.New(sc.BasePath)
}

func (sc ServerConfig) Save() error {
	data, err := yaml.Marshal(&sc)
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
