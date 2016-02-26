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
	Jobs       map[string]*job.Job
	JobCounter int
	jobsMutex  *sync.Mutex
}

// NewDummy creates a client that stores on memory
func NewDummy() *Dummy {
	logrus.Debug("New Dummy storage client created")
	return &Dummy{
		Jobs:       map[string]*job.Job{},
		JobCounter: 0,
		jobsMutex:  &sync.Mutex{},
	}
}

// Close doens't do nothing on dummy client
func (c *Dummy) Close() error {
	return nil
}

// GetJobs returns all the http jobs stored on memory
func (c *Dummy) GetJobs() (jobs []*job.Job, err error) {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()
	if c.Jobs == nil {
		return nil, errors.New("Error retrieving jobs")
	}

	for _, v := range c.Jobs {
		jobs = append(jobs, v)
	}
	return jobs, nil
}

// GetJob returns a job from memory
func (c *Dummy) GetJob(id int) (job *job.Job, err error) {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, id)
	j, ok := c.Jobs[key]
	if !ok {
		return nil, errors.New("Not existent job")
	}
	return j, nil
}

// SaveJob stores a job on memory
func (c *Dummy) SaveJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	c.JobCounter++
	j.ID = c.JobCounter
	key := fmt.Sprintf(jobKeyFmt, j.ID)
	c.Jobs[key] = j

	// Never conflict (always creates a new id)
	return nil
}

// UpdateJob updates a present job on memory
func (c *Dummy) UpdateJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)

	if _, ok := c.Jobs[key]; !ok {
		return errors.New("Not existent job")
	}

	c.Jobs[key] = j
	return nil
}

// DeleteJob Deletes a job on memory
func (c *Dummy) DeleteJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)

	if _, ok := c.Jobs[key]; !ok {
		return errors.New("Not existent job")
	}

	delete(c.Jobs, key)
	return nil
}
