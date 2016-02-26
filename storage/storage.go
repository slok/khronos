package storage

import "github.com/slok/khronos/job"

// Client implements the client of an storage
type Client interface {
	Close() error

	// Job actions
	GetJobs() ([]*job.Job, error)
	GetJob(id int) (*job.Job, error)
	SaveJob(j *job.Job) error
	DeleteJob(j *job.Job) error
}
