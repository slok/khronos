package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

const (
	// DefaultResultBufferLen is the default length that the result buffer will have if no settings are set
	DefaultResultBufferLen = 100
	// DefaultStorageEngine is the default storage engine used to store the data
	DefaultStorageEngine = "boltdb"
)

var (
	// ValidStorageEngines contains the selectable storage engines
	ValidStorageEngines = []string{"dummy", "boltdb"}
)

// Khronos holds the configuration of the main application
type Khronos struct {
	// ResultBufferLen is default result channel length
	ResultBufferLen int `envconfig:"KHRONOS_RESULT_BUFFER_LEN"`

	// StorageEngine is the engine used to store the data
	StorageEngine string `envconfig:"KHRONOS_STORAGE_ENGINE"`
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

	// Check storage Engine
	if k.StorageEngine == "" {
		k.StorageEngine = DefaultStorageEngine
	}

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
