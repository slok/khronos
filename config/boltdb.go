package config

// BoltDB  holds the configuration of storage
type BoltDB struct {
	BoltDBPath           string `envconfig:"BOLTDB_PATH"`
	BoltDBTimeoutSeconds int    `envconfig:"BOLTDB_TIMEOUT_SECONDS"`
}
