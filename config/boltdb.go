package config

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/Sirupsen/logrus"
)

// BoltDB  holds the configuration of storage
type BoltDB struct {
	BoltDBPath           string `envconfig:"BOLTDB_PATH" default:"data/khronos.db"`
	BoltDBTimeoutSeconds int    `envconfig:"BOLTDB_TIMEOUT_SECONDS" default:"1"`
}

// LoadBoltDBConfig loads boltdb env and jsonfile config
func (b *BoltDB) LoadBoltDBConfig(cfg *AppConfig) {
	if cfg.ConfigFilePath != "" {
		config.LoadJSONFile(cfg.ConfigFilePath, b)
	}

	config.LoadEnvConfig(b)

	logrus.Infof("Boltdb database path: %s", cfg.BoltDBPath)
	logrus.Infof("Boltdb timeout set to: %ds", cfg.BoltDBTimeoutSeconds)
}
