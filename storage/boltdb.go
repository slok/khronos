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

// GetJobs returns all the HTTP jobs from boltdb
func (c *BoltDB) GetJobs() (jobs []*job.Job, err error) {
	return nil, errors.New("Not implemented")
}

// GetJob returns an specific HTTP job based on the ID
func (c *BoltDB) GetJob(id int) (job *job.Job, err error) {
	return nil, errors.New("Not implemented")
}

// SaveJob stores an HTTP job on boltdb
func (c *BoltDB) SaveJob(j *job.Job) error {
	return errors.New("Not implemented")
}

// UpdateJob updates an HTTP job on boltdb
func (c *BoltDB) UpdateJob(j *job.Job) error {
	return errors.New("Not implemented")
}

// DeleteJob deletes an HTTP job from boltdb
func (c *BoltDB) DeleteJob(j *job.Job) error {
	return errors.New("Not implemented")
}
