package schedule

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/job"
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

	// application context
	cfg *config.AppConfig
}

// NewSimpleCron creates a new instance of a cron initialized with the basic functionality
func NewSimpleCron(cfg *config.AppConfig) *Cron {
	return &Cron{
		runner:    cron.New(),
		scheduler: SimpleRun(),
		Results:   make(chan *job.Result, cfg.ResultBufferLen),
		cfg:       cfg,
	}
}

// NewDummyCron creates a new instance of a cron that will execute dummy jobs
// (do nothing) when the time comes
func NewDummyCron(cfg *config.AppConfig, exitStatus int, out string) *Cron {
	return &Cron{
		runner:    cron.New(),
		scheduler: DummyRun(exitStatus, out),
		Results:   make(chan *job.Result, cfg.ResultBufferLen),
		cfg:       cfg,
	}
}

// StartResultProcesser starts the processor for the results (it blocks, needs a goroutine)
func (c *Cron) StartResultProcesser() {
	logrus.Info("Result processing started...")
	for r := range c.Results {
		logrus.Debugf("received result from job '%d' with:\nstatus:%d;\nOutput:%s", r.Job.ID, r.Status, r.Out)
		// TODO: apply result processing logic
	}
}

// Start starts cron job scheduler adn the result listener
func (c *Cron) Start() {
	c.runner.Start()
	go c.StartResultProcesser()
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
