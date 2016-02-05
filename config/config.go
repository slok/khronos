package config

import (
	"github.com/NYTimes/gizmo/config"
)

// AppConfig holds the configuration of the application
type AppConfig struct {
	*config.Server

	configFilePath string
}

// NewAppConfig creates a new app configuration with all the settings loaded
func NewAppConfig(configFile string) *AppConfig {
	cfg := &AppConfig{
		configFilePath: configFile,
		Server:         &config.Server{},
	}

	cfg.LoadConfig()

	return cfg
}

// LoadConfig loads all the application settings with a priority:
// First loads settings from file, then loads the settings from env vars
func (a *AppConfig) LoadConfig() {

	// If there is a config file load it
	if a.configFilePath != "" {
		config.LoadJSONFile(a.configFilePath, a)
		config.LoadJSONFile(a.configFilePath, a.Server)
	}

	// Load settings from env var
	config.LoadEnvConfig(a)
	config.LoadEnvConfig(a.Server)

}
