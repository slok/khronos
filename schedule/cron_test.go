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
	"github.com/slok/khronos/storage"
)

func TestRegisterCronJob(t *testing.T) {
	// Create configuration, storage and test vars
	cfg := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	stCli := storage.NewDummy()
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
	dCron := NewDummyCron(cfg, stCli, wantExitStatus, wantOut)
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

func TestRegisterCronJobStore(t *testing.T) {
	// Create configuration, storage and test vars
	cfg := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	stCli := storage.NewDummy()
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
	dCron := NewDummyCron(cfg, stCli, wantExitStatus, wantOut)
	dCron.Start(nil)

	// Register the job
	dCron.RegisterCronJob(j)

	// wait the number of iterations
	time.Sleep(time.Duration(iterations) * time.Second)

	// stop and check iterations where ok
	dCron.Stop()

	// Check stored results
	results, err := stCli.GetResults(j, 0, 0)
	if err != nil {
		t.Errorf("Error retrieving stored results: %v", err)
	}

	if len(results) != iterations {
		t.Errorf("Wrong result list; expected: %d; got: %d", iterations, len(results))
	}

	for _, r := range results {
		if r.Job != j {
			t.Errorf("Wrong result Job; expected: %d; got: %d", j, r.Job)
		}
		if r.Status != wantExitStatus {
			t.Errorf("Wrong result status; expected: %d; got: %d", wantExitStatus, r.Status)
		}
		if r.Out != wantOut {
			t.Errorf("Wrong result out; expected: %s; got: %s", wantOut, r.Out)
		}
	}

}

func TestStarCronEngine(t *testing.T) {
	cfg := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	stCli := storage.NewDummy()
	dCron := NewDummyCron(cfg, stCli, 0, "")

	// First time should start ok
	if err := dCron.Start(nil); err != nil {
		t.Errorf("Starting the first time should not get an error: %v", err)
	}

	// Second time should fail
	if err := dCron.Start(nil); err == nil {
		t.Errorf("Starting the second time should get an error")
	}

	// Stop and start again perfectly
	dCron.Stop()
	if err := dCron.Start(nil); err != nil {
		t.Errorf("Starting after stopping should not get an error: %v", err)
	}

}

func TestStopCronEngine(t *testing.T) {
	cfg := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	stCli := storage.NewDummy()
	dCron := NewDummyCron(cfg, stCli, 0, "")

	// Stopping without starting should fail
	if err := dCron.Stop(); err == nil {
		t.Errorf("Stopping without starting should get an error")
	}

	// Stopping after starting should go ok
	dCron.Start(nil)
	if err := dCron.Stop(); err != nil {
		t.Errorf("Stopping after starting should not get an error %v", err)
	}

	// Stopping after stopping should fail
	dCron.Stop()
	if err := dCron.Stop(); err == nil {
		t.Errorf("Stopping after stopping should get an error")
	}

}
