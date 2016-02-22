// Package schedule contains all the logic to execute the jobs
package schedule

import (
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

// Scheduler implements the unit of execution of a job, it mus implement run method
// this method will be the one that will execute job logic
type Scheduler interface {
	Run(*job.Result, *job.Job)
}

// SchedulerFunc is a handy scheduler that runs the function itself
type SchedulerFunc func(*job.Result, *job.Job)

// Run implements the start of the execution
func (s SchedulerFunc) Run(r *job.Result, j *job.Job) {
	// Run the chain!
	s(r, j)
}

// DummyScheduler breaks the schedule chain and returns an specific  result
func DummyScheduler(exitStatus int, resultOut string, s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		r.Out = resultOut
		r.Status = exitStatus
		logrus.Infof("Dummy HTTP request: %v", j.URL.String())
		// s.Run(r, j)
	})
}

// TimingScheduler registers the start and finish time of a job
func TimingScheduler(s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		r.Start = time.Now().UTC()
		s.Run(r, j)
		r.Finish = time.Now().UTC()
	})
}

// HTTPScheduler makes an http call to the job destination and registers the result
func HTTPScheduler(s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		//TODO: HTTP handling
		s.Run(r, j)
	})
}

// LogScheduler logs the execution of a job
func LogScheduler(s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		logrus.Infof("Start running cron '%d' at %v", j.ID, time.Now().UTC())
		s.Run(r, j)
		logrus.Infof("Start running cron '%d' at %v", j.ID, time.Now().UTC())
	})
}
