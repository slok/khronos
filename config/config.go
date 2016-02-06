package config

import "github.com/NYTimes/gizmo/config"

const defaultLogLevel = "info"

// AppConfig holds the configuration of the application
type AppConfig struct {
	*config.Server
	ConfigFilePath string
}

// NewAppConfig creates a new app configuration with all the settings loaded
func NewAppConfig(configFile string) *AppConfig {
	cfg := &AppConfig{
		ConfigFilePath: configFile,
		Server:         &config.Server{},
	}

	cfg.LoadConfig()

	return cfg
}

// LoadConfig loads all the application settings with a priority:
// First loads settings from file, then loads the settings from env vars
func (a *AppConfig) LoadConfig() {

	// If there is a config file load it
	if a.ConfigFilePath != "" {
		config.LoadJSONFile(a.ConfigFilePath, a)
		config.LoadJSONFile(a.ConfigFilePath, a.Server)
	}

	// Load settings from env var
	config.LoadEnvConfig(a)
	config.LoadEnvConfig(a.Server)

}
