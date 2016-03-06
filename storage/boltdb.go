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
	jobsBucket        = "jobs"
	resultsBucket     = "results"
	jobResultsBuckets = "job:%s:results"
)

// BoltDB client to store jobs on database
type BoltDB struct {
	BoltPath string
	Timeout  time.Duration
	DB       *bolt.DB
}

// idToByte returns an 8-byte big endian representation of int id.
// this is used to preserve the int decimal order in a byte sequence, the one
// that is used by boltdb. This representation is internal.
// examples:
// key=000000000000001f, ID=31
// key=0000000000000020, ID=32
// key=0000000000000021, ID=33
// key=0000000000000022, ID=34
func idToByte(id int) []byte {
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

	logrus.Debug("New Boltdb storage client created")
	return c, nil
}

// Close closes boltdb connection to database
func (c *BoltDB) Close() error {
	return c.Close()
}

// GetJobs returns all the HTTP jobs from boltdb. Use low and high params as slice operator
func (c *BoltDB) GetJobs(low, high int) ([]*job.Job, error) {
	jobs := []*job.Job{}

	// In database id start from one, not 0, so we convert to help in the logic
	low++
	high++

	// if low and high the same then return empty slice (not when asking for all)
	if low != 1 && low == high {
		return jobs, nil
	}

	// Check indexes ok
	if high != 1 && low >= high {
		return nil, errors.New("wrong parameters")
	}

	// Get all asked jobs
	err := c.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(jobsBucket))
		c := b.Cursor()

		var lastKey []byte
		if high == 1 {
			lastKey = nil

		} else {
			lastKey = idToByte(high)
		}

		firstKey := idToByte(low)

		// Until we reach the number asked for (we compare that the returned key is not the last one or is nil, that means the end)
		for k, v := c.Seek([]byte(firstKey)); k != nil && string(k) != string(lastKey); k, v = c.Next() {
			j := &job.Job{}
			if err := json.Unmarshal(v, j); err != nil {
				return err
			}
			jobs = append(jobs, j)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("error retrieving jobs form boltdb: %v", err)
		return nil, err
	}

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
	j := &job.Job{}
	err := c.DB.View(func(tx *bolt.Tx) error {
		// Get job from job bucket
		b := tx.Bucket([]byte(jobsBucket))
		jb := b.Get(idToByte(id))

		// Check if job is present
		if jb == nil {
			return errors.New("job does not exists")
		}

		// if present then decode from json
		if err := json.Unmarshal(jb, j); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("error retrieving job '%d' form boltdb: %v", id, err)
		return nil, err
	}

	return j, nil
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
		key := idToByte(j.ID)
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

// DeleteJob deletes an HTTP job from boltdb, doesn't return error if job doesn't exists
func (c *BoltDB) DeleteJob(j *job.Job) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(jobsBucket)).Delete(idToByte(j.ID))
	})

	if err != nil {
		err = fmt.Errorf("error deleting job '%d': %v", j.ID, err)
		logrus.Error(err.Error())
		return err
	}

	logrus.Debugf("Job '%d' deleted boltdb", j.ID)
	return nil
}

// GetResults returns all the results of a jobs from boltdb. Use low and high params as slice operator
func (c *BoltDB) GetResults(j *job.Job, low, high int) ([]*job.Result, error) {
	res := []*job.Result{}

	// In database id start from one, not 0, so we convert to help in the logic
	low++
	high++

	// if low and high the same then return empty slice (not when asking for all)
	if low != 1 && low == high {
		return res, nil
	}

	// Check indexes ok
	if high != 1 && low >= high {
		return nil, errors.New("wrong parameters")
	}

	// Get all asked jobs
	err := c.DB.View(func(tx *bolt.Tx) error {
		// Get main results bucket
		rsB := tx.Bucket([]byte(resultsBucket))
		// Get results bucket
		rbKey := fmt.Sprintf(jobResultsBuckets, string(idToByte(j.ID)))
		rB := rsB.Bucket([]byte(rbKey))

		c := rB.Cursor()

		var lastKey []byte
		if high == 1 {
			lastKey = nil

		} else {
			lastKey = idToByte(high)
		}

		firstKey := idToByte(low)

		// Until we reach the number asked for (we compare that the returned key is not the last one or is nil, that means the end)
		for k, v := c.Seek([]byte(firstKey)); k != nil && string(k) != string(lastKey); k, v = c.Next() {
			r := &job.Result{}
			if err := json.Unmarshal(v, r); err != nil {
				return err
			}
			// Set same job instance to save memory, let GC remove the Unmarshal created job
			r.Job = j
			res = append(res, r)
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("error retrieving jobs form boltdb: %v", err)
		return nil, err
	}

	// return error if not retrieved all asked for (if high is 0 means: want all from low,
	// doesn't matter how many, so in this case no errors)
	if high > 1 && len(res) != high-low {
		return res, fmt.Errorf("error retrieving all asked for; expected: %d; got: %d", high-low, len(res))
	}

	logrus.Debugf("Retrieved '%d' results from boltdb without errors", len(res))
	return res, nil

}

// GetResult retursn a single result from boltdb
func (c *BoltDB) GetResult(j *job.Job, id int) (*job.Result, error) {
	r := &job.Result{}

	err := c.DB.View(func(tx *bolt.Tx) error {
		// Get main results bucket
		rsB := tx.Bucket([]byte(resultsBucket))

		// Get results bucket
		rbKey := fmt.Sprintf(jobResultsBuckets, string(idToByte(j.ID)))
		rB := rsB.Bucket([]byte(rbKey))

		// Get result
		res := rB.Get(idToByte(id))
		if res == nil {
			return errors.New("result does not exists")
		}

		if err := json.Unmarshal(res, r); err != nil {
			return err
		}

		// Set same job instance to save memory, let GC remove the Unmarshal created job
		r.Job = j
		return nil
	})
	if err != nil {
		logrus.Errorf("error retrieving result '%d' form boltdb: %v", id, err)
		return nil, err
	}

	return r, nil
}

// SaveResult stores a result of a job on boltdb
func (c *BoltDB) SaveResult(r *job.Result) error {
	// First check if the job is present, if not, then error
	if _, err := c.GetJob(r.Job.ID); err != nil {
		return err
	}

	err := c.DB.Update(func(tx *bolt.Tx) error {
		resB := tx.Bucket([]byte(resultsBucket))

		// Create job results bucket if doens't exist
		jobresKey := fmt.Sprintf(jobResultsBuckets, string(idToByte(r.Job.ID)))
		b, err := resB.CreateBucketIfNotExists([]byte(jobresKey))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}

		// Create a new ID for the new result, not new ID if it has already (update)
		// Starts in 1, so its safe to check with 0
		if r.ID == 0 {
			id, _ := b.NextSequence()
			r.ID = int(id)
		}

		// Marshal the result
		buf, err := json.Marshal(r)
		if err != nil {
			return err
		}

		// Save the result
		return b.Put(idToByte(r.ID), buf)
	})
	if err != nil {
		err = fmt.Errorf("error storing result '%d': %v", r.ID, err)
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("Stored result '%d' boltdb", r.ID)
	return nil
}

// DeleteResult deletes a result from boltdv, doesn't return error if result doesn't exist
func (c *BoltDB) DeleteResult(r *job.Result) error {
	err := c.DB.Update(func(tx *bolt.Tx) error {
		resB := tx.Bucket([]byte(resultsBucket))
		jobresKey := fmt.Sprintf(jobResultsBuckets, string(idToByte(r.Job.ID)))
		return resB.Bucket([]byte(jobresKey)).Delete(idToByte(r.ID))
	})

	if err != nil {
		err = fmt.Errorf("error deleting result '%d': %v", r.ID, err)
		logrus.Error(err.Error())
		return err
	}

	logrus.Debugf("Res '%d' deleted boltdb", r.ID)
	return nil
}
