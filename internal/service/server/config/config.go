package config

import (
	"errors"
	"fmt"
	"log/slog"
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
	// WebhookBaseURL is where forges reach Josuke, e.g. "https://example.com":
	// the webhook of a new repository is this URL plus Josuke's hook path.
	WebhookBaseURL string `json:"webhook_base_url,omitempty" yaml:"webhook_base_url,omitempty"`
	// RestartOnSave maps a config file to the systemd unit reading it,
	// offering "Save & restart" on its page.
	RestartOnSave map[string]string `json:"restart_on_save,omitempty" yaml:"restart_on_save,omitempty"`
	// EnvBasePath holds one folder per service, with its .env:
	// <env_base_path>/<project>/.env, written by "+ Service".
	EnvBasePath string `json:"env_base_path,omitempty" yaml:"env_base_path,omitempty"`
	// DockerComposeBasePath holds the docker-compose.<type>.yml files
	// "+ Service" adds services to.
	DockerComposeBasePath string `json:"docker_compose_base_path,omitempty" yaml:"docker_compose_base_path,omitempty"`
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

// LogValue lets slog print the config as a group of its fields.
func (sc ServerConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("path", sc.ServerConfigPath),
		slog.String("base_path", sc.BasePath),
		slog.String("josuke_config_path", sc.JosukeConfigPath),
		slog.Any("capybara_config_paths", sc.CapybaraConfigPaths),
		slog.Any("services_paths", []string(sc.AppsServicesPaths)),
		slog.String("webhook_base_url", sc.WebhookBaseURL),
		slog.Any("restart_on_save", sc.RestartOnSave),
		slog.String("env_base_path", sc.EnvBasePath),
		slog.String("docker_compose_base_path", sc.DockerComposeBasePath),
	)
}

// CanRestart tells whether the config at path has a unit to restart on save.
func (sc ServerConfig) CanRestart(path string) bool {
	_, ok := sc.RestartOnSave[path]
	return ok
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
