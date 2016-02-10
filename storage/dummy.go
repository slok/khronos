package storage

import (
	"errors"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

// Our memory database
var (
	httpJobs = []*job.HTTPJob{}
)

// Dummy implements the Storage interface everything to a local memory map
type Dummy struct {
}

// NewDummy creates a client that stores on memory
func NewDummy() *Dummy {
	logrus.Debug("New Dummy storage client created")
	return &Dummy{}
}

// GetHTTPJobs returns all the http jobs stored on memory
func (c *Dummy) GetHTTPJobs() (jobs []*job.HTTPJob, err error) {
	if httpJobs == nil {
		return nil, errors.New("Error retrieving jobs")
	}
	return httpJobs, nil
}

// GetHTTPJob returns a job from memory
func (c *Dummy) GetHTTPJob(id int) (job *job.HTTPJob, err error) {
	return httpJobs[id], nil
}

// SaveHTTPJob stores a job on memory
func (c *Dummy) SaveHTTPJob(j *job.HTTPJob) error {
	id := len(httpJobs)
	j.ID = id
	httpJobs = append(httpJobs, j)
	return nil
}

// UpdateHTTPJob updates a present job on memory
func (c *Dummy) UpdateHTTPJob(j *job.HTTPJob) error {
	httpJobs[j.ID] = j
	return nil
}

// DeleteHTTPJob Deletes a job on memory
func (c *Dummy) DeleteHTTPJob(j *job.HTTPJob) error {
	httpJobs[j.ID] = nil
	return nil
}
