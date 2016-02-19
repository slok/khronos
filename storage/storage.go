package storage

import (
	"github.com/Sirupsen/logrus"
	"github.com/slok/khronos/job"
)

// Client implements the client of an storage
type Client interface {
	GetJobs() (jobs []*job.Job, err error)
	GetJob(id int) (job *job.Job, err error)
	SaveJob(j *job.Job) error
	UpdateJob(j *job.Job) error
	DeleteJob(j *job.Job) error
}

// Nil implements the Storage interface everything to nil
type Nil struct{}

// NewNil creates a nil storege client instance
func NewNil() *Nil {
	logrus.Debug("New Nil storage client created")
	return &Nil{}
}

func (c *Nil) GetJobs() (jobs []*job.Job, err error)   { return nil, nil }
func (c *Nil) GetJob(id int) (job *job.Job, err error) { return nil, nil }
func (c *Nil) SaveJob(j *job.Job) error                { return nil }
func (c *Nil) UpdateJob(j *job.Job) error              { return nil }
func (c *Nil) DeleteJob(j *job.Job) error              { return nil }
