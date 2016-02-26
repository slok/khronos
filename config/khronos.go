package config

import "github.com/NYTimes/gizmo/config"

const (
	// DefaultResultBufferLen is the default length that the result buffer will have if no settings are set
	DefaultResultBufferLen = 100
)

// Khronos holds the configuration of the main application
type Khronos struct {
	// The default result channel length
	ResultBufferLen int `envconfig:"KHRONOS_RESULT_BUFFER_LEN"`
}

// LoadKhronosConfig Loads the configuration for the application
func (k *Khronos) LoadKhronosConfig(cfg *AppConfig) {

	if cfg.ConfigFilePath != "" {
		config.LoadJSONFile(cfg.ConfigFilePath, k)
	}

	config.LoadEnvConfig(k)

	if k.ResultBufferLen == 0 {
		k.ResultBufferLen = DefaultResultBufferLen
	}

}
