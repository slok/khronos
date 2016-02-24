package schedule

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/job"
)

func TestRegisterCronJob(t *testing.T) {
	// Create configuration and test vars
	cfg := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	wantExitStatus := job.ResultOK
	wantOut := fmt.Sprintf("Result: %d", rand.Int())
	iterations := 3

	// Create a job & a result
	u, _ := url.Parse("http://test.org/test")
	j := &job.Job{
		ID:     1,
		URL:    u,
		When:   "@every 1s",
		Active: true,
	}

	// Create our cron engine and start
	dCron := NewDummyCron(cfg, wantExitStatus, wantOut)
	var results []*job.Result
	dCron.Start(func(r *job.Result) {
		results = append(results, r)
	})

	// Register the job
	dCron.RegisterCronJob(j)

	// wait the number of iterations
	time.Sleep(time.Duration(iterations) * time.Second)

	// stop and check iterations where ok
	dCron.Stop()

	// Check the number of jobs executed and completed is the expected
	if len(results) != iterations {
		t.Errorf("Wrong result list; expected: %d; got: %d", iterations, len(results))
	}

	// Check the the dummy schedule was executed ok by checking the expected results
	for _, r := range results {
		if r.Status != wantExitStatus {
			t.Errorf("Wrong result status; expected: %d; got: %d", wantExitStatus, r.Status)
		}
		if r.Out != wantOut {
			t.Errorf("Wrong result out; expected: %s; got: %s", wantOut, r.Out)
		}
	}
}
