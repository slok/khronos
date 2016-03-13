package storage

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strings"
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
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

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

}

func TestBoltDBSaveJob(t *testing.T) {
	boltPath := randomPath()
	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}

	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

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
			key := idToByte(j.ID)
			if err := json.Unmarshal(b.Get(key), gotJob); err != nil {
				return err
			}
			return nil
		})

		// Check all the info correct
		if !reflect.DeepEqual(*j, *gotJob) {
			t.Errorf("Jobs should be equal; expected: %#v,\n got: %#v", j, gotJob)
		}

		// Check ID ok
		if gotJob.ID != i+1 {
			t.Errorf("IDS should be equal; expected: %d,\n got: %d", i+1, gotJob.ID)
		}
	}
}

func TestBoltDBGetJobs(t *testing.T) {
	boltPath := randomPath()
	totalJobs := 50

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Create a buch of jobs
	for i := 1; i <= totalJobs; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://khronos.io/job%d", i))
		j := &job.Job{
			Name:        fmt.Sprintf("job%d", i),
			Description: fmt.Sprintf("job %d", i),
			When:        fmt.Sprintf("@every %dm", i),
			Active:      true,
			URL:         u,
		}
		err := c.SaveJob(j)
		if err != nil {
			t.Error("Error saving job on database")
		}
	}

	// Prepare tests
	tests := []struct {
		givenLow   int
		givenHigh  int
		wantLength int
		wantError  bool
	}{
		{
			givenLow:   0,
			givenHigh:  0,
			wantLength: totalJobs,
			wantError:  false,
		},
		{
			givenLow:   totalJobs - 20,
			givenHigh:  totalJobs - 10,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalJobs - 10,
			givenHigh:  totalJobs,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalJobs - 10,
			givenHigh:  0,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalJobs - 20,
			givenHigh:  totalJobs - 20,
			wantLength: 0,
			wantError:  false,
		},
		{
			givenLow:  totalJobs - 30,
			givenHigh: 1,
			wantError: true,
		},
		{
			givenLow:  totalJobs - 10,
			givenHigh: totalJobs + 1,
			wantError: true,
		},
		{
			givenLow:  totalJobs - 40,
			givenHigh: totalJobs + 100,
			wantError: true,
		},
	}

	for _, test := range tests {
		jobs, err := c.GetJobs(test.givenLow, test.givenHigh)

		// Check it should error or not
		if test.wantError && err == nil {
			t.Error("job retrieval didn't error when it should")
		}

		if !test.wantError && err != nil {
			t.Errorf("job retrieval error when it shouldn't: %v", err)

		}

		// Only check trusted content, this means no errors
		if err == nil {
			// Check len
			if len(jobs) != test.wantLength {
				t.Errorf("Number of retrieved jobs is wrong; expected: %d; got: %d", test.wantLength, len(jobs))
			}
			// Check content
			for k, gotJob := range jobs {
				i := test.givenLow + k + 1 // add + 1 because jobs ids start in 1 not 0
				u, _ := url.Parse(fmt.Sprintf("http://khronos.io/job%d", i))
				j := &job.Job{
					ID:          i,
					Name:        fmt.Sprintf("job%d", i),
					Description: fmt.Sprintf("job %d", i),
					When:        fmt.Sprintf("@every %dm", i),
					Active:      true,
					URL:         u,
				}
				// Check all the info correct
				if !reflect.DeepEqual(*j, *gotJob) {
					t.Errorf("Jobs should be equal; expected: %#v,\n got: %#v", j, gotJob)
				}

			}
		}

	}
}

func TestBoltDBGetJob(t *testing.T) {
	boltPath := randomPath()
	totalJobs := 50

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Create a buch of jobs
	for i := 1; i <= totalJobs; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://khronos.io/job%d", i))
		j := &job.Job{
			Name:        fmt.Sprintf("job%d", i),
			Description: fmt.Sprintf("job %d", i),
			When:        fmt.Sprintf("@every %dm", i),
			Active:      true,
			URL:         u,
		}
		err := c.SaveJob(j)
		if err != nil {
			t.Error("Error saving job on database")
		}
	}

	// Check all the stored jobs retrieving one by one
	for id := 1; id <= totalJobs; id++ {
		gotJob, err := c.GetJob(id)

		if err != nil {
			t.Errorf("Job should be retrieved, it didn't: %v", err)
		}

		u, _ := url.Parse(fmt.Sprintf("http://khronos.io/job%d", id))
		j := &job.Job{
			ID:          id,
			Name:        fmt.Sprintf("job%d", id),
			Description: fmt.Sprintf("job %d", id),
			When:        fmt.Sprintf("@every %dm", id),
			Active:      true,
			URL:         u,
		}

		if !reflect.DeepEqual(*j, *gotJob) {
			t.Errorf("Jobs should be equal, it didn't; expected: %#v;\ngot: %#v", *j, *gotJob)
		}

	}

	// Check not existent job
	_, err = c.GetJob(totalJobs + 1)
	if err == nil {
		t.Error("Expected error but didn't got")
	}

	if !strings.Contains(err.Error(), "job does not exists") {
		t.Errorf("Expected error but not this, got: %v", err)
	}
}

func TestBoltDBDeleteJob(t *testing.T) {
	boltPath := randomPath()
	totalJobs := 5
	totalResults := 10
	jobs := make([]*job.Job, totalJobs, totalJobs)

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Create a buch of jobs
	for i := 1; i <= totalJobs; i++ {
		u, _ := url.Parse(fmt.Sprintf("http://khronos.io/job%d", i))
		j := &job.Job{
			Name:        fmt.Sprintf("job%d", i),
			Description: fmt.Sprintf("job %d", i),
			When:        fmt.Sprintf("@every %dm", i),
			Active:      true,
			URL:         u,
		}
		err := c.SaveJob(j)
		if err != nil {
			t.Error("Error saving job on database")
		}
		jobs[i-1] = j

		// Create a bunch of results for this job
		for ir := 1; ir <= totalResults; ir++ {
			r := &job.Result{
				Job:    j,
				Out:    "Linux khronos-dev 4.4.1-2-ARCH #1 SMP PREEMPT Wed Feb 3 13:12:33 UTC 2016 x86_64 GNU/Linux",
				Status: job.ResultOK,
				Start:  time.Now().UTC(),
				Finish: time.Now().UTC(),
			}
			c.SaveResult(r)
		}
	}

	// Delete one by one and check is correct
	for _, j := range jobs {
		// Check in database
		if _, err := c.GetJob(j.ID); err != nil {
			t.Errorf("Job should exists, got error: %v", err)
		}

		if rs, err := c.GetResults(j, 0, 0); err != nil || len(rs) != totalResults {
			t.Errorf("Job results should exists,\n expected length %d; got %d \nerror: %v", totalResults, len(rs), err)
		}

		// Delete
		if err := c.DeleteJob(j); err != nil {
			t.Errorf("Job should be deleted, got error: %v", err)
		}

		// Check not in database
		if _, err := c.GetJob(j.ID); err == nil {
			t.Errorf("Job shouldn't exists, should got error, didn't")
		}
		if rs, err := c.GetResults(j, 0, 0); err != nil || len(rs) != 0 {
			t.Errorf("Job results should't exists,\n expected length %d; got %d \nerror: %v", 0, len(rs), err)
		}

	}
}

func TestBoltDBSaveResultJob(t *testing.T) {
	boltPath := randomPath()
	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}

	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Create job & results
	u1, _ := url.Parse("http://khronos.io/job1")
	j := &job.Job{
		Name:        "job1",
		Description: "Job 1",
		When:        "@every 1m",
		Active:      true,
		URL:         u1,
	}

	rs := []*job.Result{
		&job.Result{
			Job:    j,
			Out:    "Linux khronos-dev 4.4.1-2-ARCH #1 SMP PREEMPT Wed Feb 3 13:12:33 UTC 2016 x86_64 GNU/Linux",
			Status: job.ResultOK,
			Start:  time.Now().UTC(),
			Finish: time.Now().UTC(),
		},
		&job.Result{
			Job:    j,
			Out:    "ls: cannot open directory /root/: Permission denied",
			Status: job.ResultError,
			Start:  time.Now().UTC(),
			Finish: time.Now().UTC(),
		},
	}

	for i, r := range rs {
		// Save result on boltdb
		err := c.SaveJob(j)
		if err != nil {
			t.Errorf("Error saving job %d: %v", j.ID, err)
		}

		// Save the result
		err = c.SaveResult(r)
		if err != nil {
			t.Errorf("Error saving result %d: %v", j.ID, err)
		}

		gotRes := &job.Result{}

		// Retrieve result
		err = c.DB.View(func(tx *bolt.Tx) error {
			// Get main results bucket
			rsB := tx.Bucket([]byte(resultsBucket))
			// Get results bucket
			rbKey := fmt.Sprintf(jobResultsBuckets, string(idToByte(r.Job.ID)))
			rB := rsB.Bucket([]byte(rbKey))

			// Get result
			resKey := idToByte(r.ID)
			if err := json.Unmarshal(rB.Get(resKey), gotRes); err != nil {
				return err
			}
			return nil
		})

		// Check all the info correct
		if !reflect.DeepEqual(*r, *gotRes) {
			t.Errorf("Results should be equal; expected: %#v,\n got: %#v", r, gotRes)
		}

		// Check ID ok
		if gotRes.ID != i+1 {
			t.Errorf("IDS should be equal; expected: %d,\n got: %d", i+1, gotRes.ID)
		}
	}
}

func TestBoltDBGetResult(t *testing.T) {
	boltPath := randomPath()
	totalResults := 50
	u, _ := url.Parse("http://khronos.io/job1")
	j := &job.Job{
		Name:        "job1",
		Description: "Job 1",
		When:        "@every 1m",
		Active:      true,
		URL:         u,
	}

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Save our job
	if err := c.SaveJob(j); err != nil {
		t.Error("Error saving job on database")
	}

	// Create a buch of results
	for i := 1; i <= totalResults; i++ {
		r := &job.Result{
			Job:    j,
			Out:    fmt.Sprintf("Out %d", i),
			Status: job.ResultOK,
			Start:  time.Now().UTC(),
			Finish: time.Now().UTC(),
		}
		if err := c.SaveResult(r); err != nil {
			t.Error("Error saving result on database")
		}
	}

	// Check all the stored jobs retrieving one by one
	for id := 1; id <= totalResults; id++ {
		gotRes, err := c.GetResult(j, id)

		if err != nil {
			t.Errorf("Result should be retrieved, it didn't: %v", err)
		}

		if gotRes.Out != fmt.Sprintf("Out %d", id) {
			t.Errorf("Resul out didn't match; expected: %#v;\ngot: %#v", fmt.Sprintf("Out %d", id), gotRes.Out)
		}

		// Check Job instance is the same (pointer)
		if gotRes.Job != j {
			t.Error("Jobs should be the same, they aren't")
		}

	}

	// Check not existent job
	if _, err = c.GetResult(j, totalResults+1); err == nil {
		t.Error("Expected error but didn't got")
	}

	if !strings.Contains(err.Error(), "result does not exists") {
		t.Errorf("Expected error but not this, got: %v", err)
	}
}

func TestBoltDBGetResults(t *testing.T) {
	boltPath := randomPath()
	totalResults := 50
	u, _ := url.Parse("http://khronos.io/job1")
	j := &job.Job{
		Name:        "job1",
		Description: "Job 1",
		When:        "@every 1m",
		Active:      true,
		URL:         u,
	}

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Save our job
	if err := c.SaveJob(j); err != nil {
		t.Error("Error saving job on database")
	}

	// Create a buch of results
	for i := 1; i <= totalResults; i++ {
		r := &job.Result{
			Job:    j,
			Out:    fmt.Sprintf("Out %d", i),
			Status: job.ResultOK,
			Start:  time.Now().UTC(),
			Finish: time.Now().UTC(),
		}
		if err := c.SaveResult(r); err != nil {
			t.Error("Error saving result on database")
		}
	}

	// Prepare tests
	tests := []struct {
		givenLow   int
		givenHigh  int
		wantLength int
		wantError  bool
	}{
		{
			givenLow:   0,
			givenHigh:  0,
			wantLength: totalResults,
			wantError:  false,
		},
		{
			givenLow:   totalResults - 20,
			givenHigh:  totalResults - 10,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalResults - 10,
			givenHigh:  totalResults,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalResults - 10,
			givenHigh:  0,
			wantLength: 10,
			wantError:  false,
		},
		{
			givenLow:   totalResults - 20,
			givenHigh:  totalResults - 20,
			wantLength: 0,
			wantError:  false,
		},
		{
			givenLow:  totalResults - 30,
			givenHigh: 1,
			wantError: true,
		},
		{
			givenLow:  totalResults - 10,
			givenHigh: totalResults + 1,
			wantError: true,
		},
		{
			givenLow:  totalResults - 40,
			givenHigh: totalResults + 100,
			wantError: true,
		},
	}

	for _, test := range tests {
		res, err := c.GetResults(j, test.givenLow, test.givenHigh)

		// Check it should error or not
		if test.wantError && err == nil {
			t.Error("result retrieval didn't error when it should")
		}

		if !test.wantError && err != nil {
			t.Errorf("result retrieval error when it shouldn't: %v", err)

		}

		// Only check trusted content, this means no errors
		if err == nil {
			// Check len
			if len(res) != test.wantLength {
				t.Errorf("Number of retrieved result is wrong; expected: %d; got: %d", test.wantLength, len(res))
			}
			// Check content
			for k, gotRes := range res {
				i := test.givenLow + k + 1 // add + 1 because jobs ids start in 1 not 0

				// Check result ok
				if gotRes.Out != fmt.Sprintf("Out %d", i) {
					t.Errorf("Result out didn't match; expected: %#v;\ngot: %#v", fmt.Sprintf("Out %d", i), gotRes.Out)
				}

				// Check Job instance is the same (pointer)
				if gotRes.Job != j {
					t.Error("Jobs should be the same, they aren't")
				}
			}
		}
	}
}

func TestBoltDBDeleteResult(t *testing.T) {
	boltPath := randomPath()
	totalResults := 5
	results := make([]*job.Result, totalResults, totalResults)
	u, _ := url.Parse("http://khronos.io/job1")
	j := &job.Job{
		Name:        "job1",
		Description: "Job 1",
		When:        "@every 1m",
		Active:      true,
		URL:         u,
	}

	// Create a new boltdb connection
	c, err := NewBoltDB(boltPath, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating bolt connection: %v", err)
	}
	// Close ok
	defer func() {
		if err := tearDownBoltDB(c.DB); err != nil {
			t.Error(err)
		}
	}()

	// Save our job
	if err := c.SaveJob(j); err != nil {
		t.Error("Error saving job on database")
	}

	// Create a buch of results
	for i := 1; i <= totalResults; i++ {
		r := &job.Result{
			Job:    j,
			Out:    fmt.Sprintf("Out %d", i),
			Status: job.ResultOK,
			Start:  time.Now().UTC(),
			Finish: time.Now().UTC(),
		}
		if err := c.SaveResult(r); err != nil {
			t.Error("Error saving result on database")
		}

		results[i-1] = r
	}

	// Delete one by one and check is correct
	for _, r := range results {
		// Check in datbase
		if _, err := c.GetResult(j, r.ID); err != nil {
			t.Errorf("Result should exists, got error: %v", err)
		}

		// Delete
		if err := c.DeleteResult(r); err != nil {
			t.Errorf("Result should be deleted, got error: %v", err)
		}

		// Check not in database
		if _, err := c.GetResult(j, r.ID); err == nil {
			t.Errorf("Result shouldn't exists, should got error, didn't")
		}
	}
}
