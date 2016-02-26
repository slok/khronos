package schedule

import (
	"time"

	"github.com/slok/khronos/job"
)

const timeout = 2 * time.Second

// SimpleRun has the simples execution flow of a job, log, time, and http
func SimpleRun() Scheduler {
	final := SchedulerFunc(func(r *job.Result, j *job.Job) {})
	s := LogScheduler(
		TimingScheduler(
			HTTPScheduler(timeout, final))) // TODO: Custom timeout per job
	return s
}

// DummyRun has a dummy chain for tests
func DummyRun(exitStatus int, result string) Scheduler {
	final := SchedulerFunc(func(r *job.Result, j *job.Job) {})
	s := LogScheduler(
		TimingScheduler(
			DummyScheduler(exitStatus, result, final)))
	return s
}
