package storage

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/boltdb/bolt"

	"github.com/slok/khronos/job"
)

// tearDownBoltDB closes and deletes boltdb
func tearDownBoltDB(db *bolt.DB) error {
	p := db.Path()
	err := db.Close()
	if err != nil {
		return err
	}

	return os.Remove(p)
}

func randomPath() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("/tmp/khronos_boltdb_test_%d.db", r.Int())
}

func TestBoltDBConnection(t *testing.T) {
	boltPath := randomPath()
	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}

	// Check root buckets are present
	checkBuckets := []string{jobsBucket, resultsBucket}
	err = c.DB.View(func(tx *bolt.Tx) error {
		for _, cb := range checkBuckets {
			if b := tx.Bucket([]byte(cb)); b == nil {
				t.Errorf("Bucket %s not present", cb)
			}
		}
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	// Close ok
	if err := tearDownBoltDB(c.DB); err != nil {
		t.Error(err)
	}
}

func TestBoltdbSaveJob(t *testing.T) {
	boltPath := randomPath()
	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}

	// Create jobs
	u1, _ := url.Parse("http://khronos.io/job1")
	u2, _ := url.Parse("http://khronos.io/job2")
	jobs := []*job.Job{
		&job.Job{
			Name:        "job1",
			Description: "Job 1",
			When:        "@every 1m",
			Active:      true,
			URL:         u1,
		},
		&job.Job{
			Name:        "job2",
			Description: "Job 2",
			When:        "@every 2m",
			Active:      false,
			URL:         u2,
		},
	}

	for i, j := range jobs {
		// Save jobs on boltdb
		err := c.SaveJob(j)
		if err != nil {
			t.Errorf("Error saving job %d: %v", j.ID, err)
		}

		gotJob := &job.Job{}

		// Retrieve job
		err = c.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(jobsBucket))
			key := fmt.Sprintf(jobKey, j.ID)
			if err := json.Unmarshal(b.Get([]byte(key)), gotJob); err != nil {
				return err
			}
			return nil
		})

		// Check all the info correct
		if !reflect.DeepEqual(*j, *gotJob) {
			t.Errorf("URLS should be equal; expected: %#v,\n got: %#v", j, gotJob)
		}

		// Check ID ok
		if gotJob.ID != i+1 {
			t.Errorf("IDS should be equal; expected: %d,\n got: %d", i+1, gotJob.ID)
		}

	}
}
