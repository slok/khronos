package storage

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"

	"github.com/slok/khronos/job"
)

// BoltDB client to store jobs on database
type BoltDB struct {
	BoltPath string
	Timeout  time.Duration
	DB       *bolt.DB
}

// NewBoltDB creates a boltdb client
func NewBoltDB(path string, timeout time.Duration) (c *BoltDB, err error) {

	// Open connection
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: timeout})
	if err != nil {
		return nil, err
	}

	// Create our client
	c = &BoltDB{
		BoltPath: path,
		Timeout:  timeout,
		DB:       db,
	}
	return
}

// GetHTTPJobs returns all the HTTP jobs from boltdb
func (c *BoltDB) GetHTTPJobs() (jobs []*job.HTTPJob, err error) {
	return nil, errors.New("Not implemented")
}

// GetHTTPJob returns an specific HTTP job based on the ID
func (c *BoltDB) GetHTTPJob(id int) (job *job.HTTPJob, err error) {
	return nil, errors.New("Not implemented")
}

// SaveHTTPJob stores an HTTP job on boltdb
func (c *BoltDB) SaveHTTPJob(j *job.HTTPJob) error {
	return errors.New("Not implemented")
}

// UpdateHTTPJob updates an HTTP job on boltdb
func (c *BoltDB) UpdateHTTPJob(j *job.HTTPJob) error {
	return errors.New("Not implemented")
}

// DeleteHTTPJob deletes an HTTP job from boltdb
func (c *BoltDB) DeleteHTTPJob(j *job.HTTPJob) error {
	return errors.New("Not implemented")
}
