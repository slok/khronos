package schedule

import (
	"net/url"
	"testing"

	"github.com/slok/khronos/job"
)

func TestDummyScheduler(t *testing.T) {

	tests := []struct {
		givenExitStatus int
		givenOut        string
	}{
		{givenExitStatus: job.ResultOK, givenOut: "test1"},
		{givenExitStatus: job.ResultError, givenOut: "test2"},
		{givenExitStatus: job.ResultUnknow, givenOut: "test3"},
	}

	for _, test := range tests {
		u, _ := url.Parse("http://test.org/test")
		j := &job.Job{URL: u}
		r := &job.Result{}

		DummyScheduler(test.givenExitStatus, test.givenOut,
			SchedulerFunc(func(r *job.Result, j *job.Job) {})).
			Run(r, j)

		if test.givenOut != r.Out {
			t.Errorf("result output should be: %s; got: %s", test.givenOut, r.Out)
		}
		if test.givenExitStatus != r.Status {
			t.Errorf("result exit status should be: %d; got: %d", test.givenExitStatus, r.Status)
		}
	}
}

func TestTimingScheduler(t *testing.T) {
	j := &job.Job{}
	r := &job.Result{}

	TimingScheduler(SchedulerFunc(func(r *job.Result, j *job.Job) {})).Run(r, j)

	if !r.Start.Before(r.Finish) || !r.Finish.After(r.Start) {
		t.Errorf("Wrong timing on job execution, start: %v; finish: %v", r.Start, r.Finish)
	}
}
