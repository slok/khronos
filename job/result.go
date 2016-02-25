package job

import "time"

const (
	// ResultOK means that the result was ok
	ResultOK = iota
	// ResultError means that the result end in error
	ResultError
	// ResultInternalError means that the job execution had an internal error
	ResultInternalError
	// ResultUnknow means that the result was not clear how it ended
	ResultUnknow
)

// Result has the result of a job
type Result struct {
	Job    *Job
	Out    string
	Status int
	Start  time.Time
	Finish time.Time
}
