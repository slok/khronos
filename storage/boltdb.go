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
	"encoding/binary"
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
	resultKey        = "%d"
)

// BoltDB client to store jobs on database
type BoltDB struct {
	BoltPath string
	Timeout  time.Duration
	DB       *bolt.DB
}

// jobKey returns an 8-byte big endian representation of int id.
// this is used to preserve the int decimal order in a byte sequence, the one
// that is used by boltdb. This representation is internal, the users of this
// client will not know anything.
// examples:
// key=000000000000001f, ID=31
// key=0000000000000020, ID=32
// key=0000000000000021, ID=33
// key=0000000000000022, ID=34

func jobKey(id int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return b
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

// GetJobs returns all the HTTP jobs from boltdb. Use low and high params as slice operator
func (c *BoltDB) GetJobs(low, high int) ([]*job.Job, error) {
	// In database id start from one, not 0, so we convert to help in the logic
	low++
	high++

	jobs := []*job.Job{}

	// Check indexes ok
	if high != 1 && low >= high {
		return nil, errors.New("wrong parameters")
	}

	// Get all asked jobs
	c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobsBucket))
		c := b.Cursor()

		var lastKey []byte
		if high == 1 {
			lastKey = nil

		} else {
			lastKey = jobKey(high)
		}

		firstKey := jobKey(low)

		// Until we reach the number asked for (we compare that the returned key is not the last one)
		for k, v := c.Seek([]byte(firstKey)); string(k) != string(lastKey); k, v = c.Next() {
			j := &job.Job{}
			if err := json.Unmarshal(v, j); err != nil {
				return err
			}
			jobs = append(jobs, j)
		}
		return nil
	})

	// return error if not retrieved all asked for (if high is 0 means: want all from low,
	// doesn't matter how many, so in this case no errors)
	if high > 1 && len(jobs) != high-low {
		return jobs, fmt.Errorf("error retrieving all asked for; expected: %d; got: %d", high-low, len(jobs))
	}

	logrus.Debugf("Retrieved '%d' jobs from boltdb without errors", len(jobs))
	return jobs, nil
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
		key := jobKey(j.ID)
		return b.Put(key, buf)
	})

	if err != nil {
		err = fmt.Errorf("error storing job '%d': %v", j.ID, err)
		logrus.Error(err.Error())
		return err
	}

	logrus.Debugf("Stored job '%d' boltdb", j.ID)
	return nil
}

// DeleteJob deletes an HTTP job from boltdb
func (c *BoltDB) DeleteJob(j *job.Job) error {
	return errors.New("Not implemented")
}
