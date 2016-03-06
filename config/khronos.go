package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

var (
	// ValidStorageEngines contains the selectable storage engines
	ValidStorageEngines = []string{"dummy", "boltdb"}
)

// Khronos holds the configuration of the main application
type Khronos struct {
	// ResultBufferLen is default result channel length
	ResultBufferLen int `envconfig:"KHRONOS_RESULT_BUFFER_LEN" default:"100"`

	// StorageEngine is the engine used to store the data
	StorageEngine string `envconfig:"KHRONOS_STORAGE_ENGINE" default:"boltdb"`
}

// LoadKhronosConfig Loads the configuration for the application
func (k *Khronos) LoadKhronosConfig(cfg *AppConfig) {

	if cfg.ConfigFilePath != "" {
		config.LoadJSONFile(cfg.ConfigFilePath, k)
	}

	config.LoadEnvConfig(k)

	// Check storage Engine
	valid := false
	for _, v := range ValidStorageEngines {
		if v == k.StorageEngine {
			valid = true
			break
		}
	}
	if !valid {
		logrus.Fatal("Incorrect storage engine")
	}

	logrus.Infof("Using '%s' storage engine", cfg.StorageEngine)
	logrus.Infof("Set result buffer length to %d", cfg.ResultBufferLen)
}
