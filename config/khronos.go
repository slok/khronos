package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

var (
	// ValidStorageEngines contains the selectable storage engines
	ValidStorageEngines = []string{"dummy", "boltdb"}

	// Defaults
	resultBufferLenDefault     = 100
	storageEngineDefault       = "boltdb"
	apiResourcesPerPageDefault = 20
)

// Khronos holds the configuration of the main application
type Khronos struct {
	// ResultBufferLen is default result channel length
	ResultBufferLen int `envconfig:"KHRONOS_RESULT_BUFFER_LEN"`

	// StorageEngine is the engine used to store the data
	StorageEngine string `envconfig:"KHRONOS_STORAGE_ENGINE"`

	//DontScheduleJobsStart flag, specifies to not schedule jobs at app startup
	DontScheduleJobsStart bool `envconfig:"KHRONOS_DONT_SCHEDULE_JOBS_ON_START"`

	//APIResourcesPerPage integer, specifies how many objects will the API return
	APIResourcesPerPage int `envconfig:"KHRONOS_API_RESOURCES_PER_PAGE"`

	// Disable API security
	APIDisableSecurity bool `envconfig:"KHRONOS_API_DISABLE_SECURITY"`
}

// LoadKhronosConfig Loads the configuration for the application
func (k *Khronos) LoadKhronosConfig(cfg *AppConfig) {

	if cfg.ConfigFilePath != "" {
		config.LoadJSONFile(cfg.ConfigFilePath, k)
	}

	config.LoadEnvConfig(k)

	// Load defaults
	k.LoadDefaults()

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

	if k.DontScheduleJobsStart {
		logrus.Warning("Not loading jobs on startup")
	} else {
		logrus.Infof("Loading jobs on startup active")
	}

	if k.APIDisableSecurity {
		logrus.Warning("Security of the API is disabled!")
	}
}

// LoadDefaults loads defaults settings
func (k *Khronos) LoadDefaults() {
	if k.ResultBufferLen == 0 {
		k.ResultBufferLen = resultBufferLenDefault
	}

	if k.StorageEngine == "" {
		k.StorageEngine = storageEngineDefault
	}

	if k.APIResourcesPerPage == 0 {
		k.APIResourcesPerPage = apiResourcesPerPageDefault
	}
}
