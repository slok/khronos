package storage

import "github.com/slok/khronos/job"

// Client implements the client of an storage
type Client interface {
	Close() error

	// Job actions
	// GetJobs returns an slice of job instances; the low parmeter will be the
	// first job and the high will be the next one to the last job that will be
	// returned; this acts like an slice operator. 0 on high parameter means all.
	// this would be translated as jobs[low:] and 0 on low would be jobs[:high]
	GetJobs(low, high int) ([]*job.Job, error)

	// GetJob returns a job by ID
	GetJob(id int) (*job.Job, error)

	// SaveJob stores the job; this method works as an insert or update, the
	// method will know if the job needs to be updated or inserted by identifying
	// the presence of the ID. This wil save as a batch so on an update the
	// instance should have all the fields set
	SaveJob(j *job.Job) error

	// DeleteJob deletes a job
	DeleteJob(j *job.Job) error

	// JobsLength returns the number of jobs stored
	JobsLength() int

	// Result actions
	// GetResults returns an slice of results from a job; The low parmeter will be the
	// first result and the high will be the next one to the last result that will be
	// returned; this acts like an slice operator. 0 on high parameter means all.
	// this would be translated as results[low:] and 0 on low would be results[:high]
	GetResults(j *job.Job, low, high int) ([]*job.Result, error)

	// GetResult returns a result based on the result id of a job
	GetResult(j *job.Job, id int) (*job.Result, error)

	// SaveResult stores the result.  Results cannot be updated
	SaveResult(r *job.Result) error

	// DeleteResult deletes a result
	DeleteResult(r *job.Result) error

	// ResultsLength returns the number of results (of a job) stored
	ResultsLength(j *job.Job) int
}
