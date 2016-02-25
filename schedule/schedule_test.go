package schedule

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

func TestHTTPScheduler(t *testing.T) {
	tests := []struct {
		givenBody  string
		givenCode  int
		wantStatus int
		timeout    time.Duration
	}{
		{
			givenBody:  "Job finished without errors",
			givenCode:  http.StatusOK,
			wantStatus: job.ResultOK,
			timeout:    2 * time.Second,
		},
		{
			givenBody:  "Job finished with errors",
			givenCode:  http.StatusInternalServerError,
			wantStatus: job.ResultError,
			timeout:    2 * time.Second,
		},
		{
			givenBody:  "dunno what do you want...",
			givenCode:  http.StatusBadRequest,
			wantStatus: job.ResultUnknow,
			timeout:    2 * time.Second,
		},
		{
			wantStatus: job.ResultInternalError,
			timeout:    3,
		},
	}

	for _, test := range tests {
		// Create our fake server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.givenCode)
			fmt.Fprint(w, test.givenBody)
		}))
		defer ts.Close()

		// Create the job
		u, _ := url.Parse(ts.URL)
		j := &job.Job{
			URL: u,
		}
		r := &job.Result{}

		// Use the scheduler
		HTTPScheduler(test.timeout, SchedulerFunc(func(r *job.Result, j *job.Job) {})).Run(r, j)

		// Check result is ok
		// Only check the body of calls that where returned ok
		if test.wantStatus != job.ResultInternalError && test.givenBody != r.Out {
			t.Errorf("result output should be: %s; got: %s", test.givenBody, r.Out)
		}
		// Always check result exit code
		if test.wantStatus != r.Status {
			t.Errorf("result exit status should be: %d; got: %d", test.wantStatus, r.Status)
		}
	}
}
