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

	cfg.ConfigureApp()

	return cfg
}

func (a *AppConfig) loadConfigFromFile() {
	config.LoadJSONFile(a.ConfigFilePath, a)
	config.LoadJSONFile(a.ConfigFilePath, a.Server)
}

func (a *AppConfig) loadConfigFromEnv() {
	config.LoadEnvConfig(a)
	config.LoadEnvConfig(a.Server)
}

// ConfigureApp loads all the application settings with a priority: First loads
// settings from file, then loads the settings from env vars.
// Initializes the required stuff (like databases)
func (a *AppConfig) ConfigureApp() {

	// Load configurations
	// Only load config file if present
	if a.ConfigFilePath != "" {
		a.loadConfigFromFile()
	}
	a.loadConfigFromEnv()
}
