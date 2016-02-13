package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

const jobKeyFmt = "job:%d"

// Dummy implements the Storage interface everything to a local memory map
type Dummy struct {
	// Our memory database
	HTTPJobs       map[string]*job.HTTPJob
	HTTPJobCounter int
	httpJobsMutex  *sync.Mutex
}

// NewDummy creates a client that stores on memory
func NewDummy() *Dummy {
	logrus.Debug("New Dummy storage client created")
	return &Dummy{
		HTTPJobs:       map[string]*job.HTTPJob{},
		HTTPJobCounter: 0,
		httpJobsMutex:  &sync.Mutex{},
	}
}

// GetHTTPJobs returns all the http jobs stored on memory
func (c *Dummy) GetHTTPJobs() (jobs []*job.HTTPJob, err error) {
	c.httpJobsMutex.Lock()
	defer c.httpJobsMutex.Unlock()
	if c.HTTPJobs == nil {
		return nil, errors.New("Error retrieving jobs")
	}

	for _, v := range c.HTTPJobs {
		jobs = append(jobs, v)
	}
	return jobs, nil
}

// GetHTTPJob returns a job from memory
func (c *Dummy) GetHTTPJob(id int) (job *job.HTTPJob, err error) {
	c.httpJobsMutex.Lock()
	defer c.httpJobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, id)
	j, ok := c.HTTPJobs[key]
	if !ok {
		return nil, errors.New("Not existent job")
	}
	return j, nil
}

// SaveHTTPJob stores a job on memory
func (c *Dummy) SaveHTTPJob(j *job.HTTPJob) error {
	c.httpJobsMutex.Lock()
	defer c.httpJobsMutex.Unlock()

	c.HTTPJobCounter++
	j.ID = c.HTTPJobCounter
	key := fmt.Sprintf(jobKeyFmt, j.ID)
	c.HTTPJobs[key] = j

	// Never conflict (always creates a new id)
	return nil
}

// UpdateHTTPJob updates a present job on memory
func (c *Dummy) UpdateHTTPJob(j *job.HTTPJob) error {
	c.httpJobsMutex.Lock()
	defer c.httpJobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)

	if _, ok := c.HTTPJobs[key]; !ok {
		return errors.New("Not existent job")
	}

	c.HTTPJobs[key] = j
	return nil
}

// DeleteHTTPJob Deletes a job on memory
func (c *Dummy) DeleteHTTPJob(j *job.HTTPJob) error {
	c.httpJobsMutex.Lock()
	defer c.httpJobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)

	if _, ok := c.HTTPJobs[key]; !ok {
		return errors.New("Not existent job")
	}

	delete(c.HTTPJobs, key)
	return nil
}
