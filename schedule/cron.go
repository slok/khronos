package schedule

import (
	"errors"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/job"
	"github.com/slok/khronos/storage"
)

// Cron is be the registerer of jobs
type Cron struct {
	// runner is the lowlevel cron that runs the jobs on separate gorutines
	runner *cron.Cron

	// Scheduler is the scheduler (chain of schedulers) that will be wrapped around the job
	// see `schedule/run` for the available default schedulers
	scheduler Scheduler

	// Results is a blocking channel where the crons will post their results as
	// a notification event of finished (ideally the actions of a finish result should
	// be in the chain of the scheduler (see `scheduler/run`)
	Results chan *job.Result

	// started flag is up if any other cron is up
	started    bool
	startMutex *sync.Mutex

	// loaded flag is up when stored jobs from database are loaded
	storedlJobsLoaded bool

	// application context
	cfg *config.AppConfig

	// Storage client
	storage storage.Client
}

// NewSimpleCron creates a new instance of a cron initialized with the basic functionality
func NewSimpleCron(cfg *config.AppConfig, storage storage.Client) *Cron {
	return &Cron{
		runner:            cron.New(),
		scheduler:         SimpleRun(),
		Results:           make(chan *job.Result, cfg.ResultBufferLen),
		started:           false,
		startMutex:        &sync.Mutex{},
		storedlJobsLoaded: false,
		cfg:               cfg,
		storage:           storage,
	}
}

// NewDummyCron creates a new instance of a cron that will execute dummy jobs
// (do nothing) when the time comes
func NewDummyCron(cfg *config.AppConfig, storage storage.Client, exitStatus int, out string) *Cron {
	return &Cron{
		runner:            cron.New(),
		scheduler:         DummyRun(exitStatus, out),
		Results:           make(chan *job.Result, cfg.ResultBufferLen),
		started:           false,
		startMutex:        &sync.Mutex{},
		storedlJobsLoaded: false,
		cfg:               cfg,
		storage:           storage,
	}
}

// startResultProcesser starts the processor for the results (runs in a goroutine)
func (c *Cron) startResultProcesser(f func(*job.Result)) error {
	// Started already aquired from Start
	if c.started {
		return errors.New("Already running")
	}

	// Apply default logic for default processing
	if f == nil {
		f = func(r *job.Result) {
			logrus.Debugf("received result from job '%d' with:\nstatus:%d;\nOutput:%s", r.Job.ID, r.Status, r.Out)

			// Save result
			err := c.storage.SaveResult(r)
			if err != nil {
				logrus.Errorf("error saving result '%d' from job '%d'", r.ID, r.Job.ID)
			}
			logrus.Debugf("Saved result '%d' for job id '%d'", r.ID, r.Job.ID)
			// TODO: apply result processing logic
		}
	}

	logrus.Info("Result processing started...")
	// Start job runner in a gouroutine. This anom func will execute the received
	// func for each result
	go func() {
		for r := range c.Results {
			f(r)
		}
	}()
	return nil
}

// Start starts cron job scheduler
// starts the result listener.
// loads stored jobs from database if DontScheduleJobsStart option is disabled
// f parameter is the function that will be executed for each result, could be nil.
func (c *Cron) Start(f func(*job.Result)) error {
	c.startMutex.Lock()
	defer c.startMutex.Unlock()

	if c.started {
		return errors.New("Already running")
	}

	// Start the cron runner
	c.runner.Start()

	// Start the result processor
	if err := c.startResultProcesser(f); err != nil {
		return err
	}

	// Register database cron jobs
	if !c.cfg.DontScheduleJobsStart && !c.storedlJobsLoaded {
		if err := c.registerStoredCronJobs(); err != nil {
			return err
		}
		c.storedlJobsLoaded = true
	}

	c.started = true
	return nil
}

// Stop stops cron job scheduler and result listener
func (c *Cron) Stop() error {
	c.startMutex.Lock()
	defer c.startMutex.Unlock()

	if !c.started {
		return errors.New("Not running")
	}
	c.runner.Stop()
	close(c.Results)
	c.started = false
	return nil
}

// RegisterCronJob registers a cron to be run when its it's time
func (c *Cron) RegisterCronJob(j *job.Job) {
	logrus.Debugf("Registering cron job: '%d'", j.ID)

	// Wrap the execution of job
	jobExec := func() {
		logrus.Debugf("Start running cron '%d' at %v", j.ID, time.Now().UTC())
		r := &job.Result{Job: j}
		c.scheduler.Run(r, j)
		c.Results <- r
		logrus.Debugf("Finished running cron '%d' at %v", j.ID, time.Now().UTC())
	}

	// Add job to  cron
	c.runner.AddFunc(j.When, jobExec)
}

// registerStoredCronJobs registers all the stored cron jobs
func (c *Cron) registerStoredCronJobs() error {
	logrus.Debug("Registering stored jobs")

	// Get storage jobs
	js, err := c.storage.GetJobs(0, 0)
	if err != nil {
		return err
	}

	// Register all the jobs
	for _, j := range js {
		// TODO: Register only active ones (check or retrieve only active ones from database?)
		c.RegisterCronJob(j)
	}
	return nil
}
