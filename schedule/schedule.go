// Package schedule contains all the logic to execute the jobs
package schedule

import (
	"io/ioutil"
	"net/http"
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
// http executed jobs should return 200 if job went ok and 500 if it went wrong,
// everything else will be interpreted as unknown
func HTTPScheduler(timeout time.Duration, s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		// Create a custom client for each request based on the timeout
		c := http.Client{Timeout: timeout}

		// Using if else because we never return, always continue the scheduler flow to the next one
		// Get the http call
		req, err := http.NewRequest("GET", j.URL.String(), nil)
		if err != nil {
			logrus.Errorf("Error creating response '%s': %s", j.URL.String(), err)
			r.Status = job.ResultInternalError
			r.Out = err.Error()
		} else {
			resp, err := c.Do(req)
			// If err then set as internal error executing job
			if err != nil {
				logrus.Errorf("Error making call to '%s': %s", j.URL.String(), err)
				r.Status = job.ResultInternalError
				r.Out = err.Error()
			} else {
				defer resp.Body.Close()

				switch resp.StatusCode {
				//if 200 then ok
				case 200:
					r.Status = job.ResultOK
					// if 500 then error
				case 500:
					r.Status = job.ResultError
					//else then unknown
				default:
					r.Status = job.ResultUnknow
				}
				// If error getting the job result then mark as wrong
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logrus.Errorf("Error reading body of '%s': º%s", j.URL.String(), err)
					r.Status = job.ResultInternalError
					r.Out = err.Error()
				} else {
					r.Out = string(b)
				}
			}
		}
		s.Run(r, j)
	})
}

// LogScheduler logs the execution of a job
func LogScheduler(s Scheduler) Scheduler {
	return SchedulerFunc(func(r *job.Result, j *job.Job) {
		logrus.Infof("Start running cron '%d' at %v", j.ID, time.Now().UTC())
		s.Run(r, j)
		logrus.Infof("Stop running cron '%d' at %v", j.ID, time.Now().UTC())
	})
}
