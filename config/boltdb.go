package config

import "time"

// BoltDB  holds the configuration of storage
type BoltDB struct {
	BoltDBPath    string        `envconfig:"BOLTDB_PATH"`
	BoltDBTimeout time.Duration `envconfig:"BOLTDB_TIMEOUT"`
}
