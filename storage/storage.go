package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/slok/khronos/job"
)

// Client implements the client of an storage
type Client interface {
	GetHTTPJobs() (jobs []*job.HTTPJob, err error)
	GetHTTPJob(id int) (job *job.HTTPJob, err error)
	SaveHTTPJob(j *job.HTTPJob) error
	UpdateHTTPJob(j *job.HTTPJob) error
	DeleteHTTPJob(j *job.HTTPJob) error
}

// Nil implements the Storage interface everything to nil
type Nil struct{}

// NewNil creates a nil storege client instance
func NewNil() *Nil {
	logrus.Debug("New Nil storage client created")
	return &Nil{}
}

func (c *Nil) GetHTTPJobs() (jobs []*job.HTTPJob, err error)   { return nil, nil }
func (c *Nil) GetHTTPJob(id int) (job *job.HTTPJob, err error) { return nil, nil }
func (c *Nil) SaveHTTPJob(j *job.HTTPJob) error                { return nil }
func (c *Nil) UpdateHTTPJob(j *job.HTTPJob) error              { return nil }
func (c *Nil) DeleteHTTPJob(j *job.HTTPJob) error              { return nil }
