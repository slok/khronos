package schedule

import "github.com/slok/khronos/job"

// SimpleRun has the simples execution flow of a job, log, time, and http
func SimpleRun(j *job.Job) *job.Result {
	var s SchedulerFunc
	r := &job.Result{}

	LogScheduler(
		TimingScheduler(
			HTTPScheduler(s))).Run(r, j)
	return r
}
