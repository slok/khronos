/*
BoltDB storage for jobs and results consist on.

Jobs are stored in a bucket named "jobs"; in this bucket the job key is an incremental ID.
Results are stored in a bucket named "results", in this bucket will be more buckets
named "job:ID:results" that will have the results identified with an incremental ID.

.
├── jobs
│   ├── 1
│   ├── 2
│   └── 3
└── results
   ├── job:1:results
   │   ├── 1
   │   ├── 2
   │   └── 3
   ├── job:2:results
   │   └── 1
   └── job:3:results
	   ├── 1
	   └── 2

*/

package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"

	"github.com/slok/khronos/job"
)

const (
	jobsBucket       = "jobs"
	resultsBucket    = "results"
	jobResultBuckets = "job:%d:result"
	jobKey           = "%d"
	resultKey        = "%d"
)

// BoltDB client to store jobs on database
type BoltDB struct {
	BoltPath string
	Timeout  time.Duration
	DB       *bolt.DB
}

// NewBoltDB creates a boltdb client
func NewBoltDB(path string, timeout time.Duration) (*BoltDB, error) {

	// Open connection
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: timeout})
	if err != nil {
		return nil, err
	}

	// Create top level buckets if necessary
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(jobsBucket))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(resultsBucket))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}

		return nil
	})

	// Create our client
	c := &BoltDB{
		BoltPath: path,
		Timeout:  timeout,
		DB:       db,
	}
	return c, nil
}

// Close closes boltdb connection to database
func (c *BoltDB) Close() error {
	return c.Close()
}

// GetJobs returns all the HTTP jobs from boltdb
func (c *BoltDB) GetJobs() ([]*job.Job, error) {
	return nil, errors.New("Not implemented")
}

// GetJob returns an specific HTTP job based on the ID
func (c *BoltDB) GetJob(id int) (*job.Job, error) {
	return nil, errors.New("Not implemented")
}

// SaveJob stores an HTTP job on boltdb
func (c *BoltDB) SaveJob(j *job.Job) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobsBucket))

		// Create a new ID for the new job, not new ID if it has already (update)
		// Starts in 1, so its safe to check with 0
		if j.ID == 0 {
			id, _ := b.NextSequence()
			j.ID = int(id)
		}

		// Marshal to json our job
		buf, err := json.Marshal(j)
		if err != nil {
			return err
		}

		// save as always (insert or update doesn't matter)
		key := fmt.Sprintf(jobKey, j.ID)
		return b.Put([]byte(key), buf)
	})

	if err != nil {
		err = fmt.Errorf("error storing job '%d': %v", j.ID, err)
		logrus.Error(err.Error())
		return err
	}

	return nil
}

// DeleteJob deletes an HTTP job from boltdb
func (c *BoltDB) DeleteJob(j *job.Job) error {
	return errors.New("Not implemented")
}
