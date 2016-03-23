package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

const jobKeyFmt = "job:%d"
const jobResultsKeyFmt = "job:%d:results"
const resultKeyFmt = "result:%d"

// Dummy implements the Storage interface everything to a local memory map
type Dummy struct {
	// Our memory database
	jobsMutex  *sync.Mutex
	Jobs       map[string]*job.Job
	JobCounter int

	resultsMutex   *sync.Mutex
	Results        map[string]map[string]*job.Result
	ResultsCounter map[string]int

	TokenMutex *sync.Mutex
	Tokens     map[string]struct{}
}

// NewDummy creates a client that stores on memory
func NewDummy() *Dummy {
	logrus.Debug("New Dummy storage client created")
	return &Dummy{
		jobsMutex:  &sync.Mutex{},
		Jobs:       map[string]*job.Job{},
		JobCounter: 0,

		resultsMutex:   &sync.Mutex{},
		Results:        map[string]map[string]*job.Result{},
		ResultsCounter: map[string]int{},

		TokenMutex: &sync.Mutex{},
		Tokens:     map[string]struct{}{},
	}
}

// Close doens't do nothing on dummy client
func (c *Dummy) Close() error {
	return nil
}

// GetJobs returns all the http jobs stored on memory
func (c *Dummy) GetJobs(low, high int) (jobs []*job.Job, err error) {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()
	if c.Jobs == nil {
		return nil, errors.New("Error retrieving jobs")
	}

	// High on top means all
	if high == 0 {
		high = len(c.Jobs)
	}

	// Check indexes ok
	if low > high || low > len(c.Jobs) || high > len(c.Jobs) {
		return nil, errors.New("wrong parameters")
	}

	// Add 1 to match the indexes
	low++
	high++
	for i := low; i < high; i++ {
		jobs = append(jobs, c.Jobs[fmt.Sprintf(jobKeyFmt, i)])
	}
	return jobs, nil
}

// GetJob returns a job from memory
func (c *Dummy) GetJob(id int) (job *job.Job, err error) {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, id)
	j, ok := c.Jobs[key]
	if !ok {
		return nil, errors.New("Not existent job")
	}
	return j, nil
}

// SaveJob stores a job on memory
func (c *Dummy) SaveJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	c.JobCounter++
	j.ID = c.JobCounter
	key := fmt.Sprintf(jobKeyFmt, j.ID)
	c.Jobs[key] = j

	// Never conflict (always creates a new id)
	return nil
}

// UpdateJob updates a present job on memory
func (c *Dummy) UpdateJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)

	if _, ok := c.Jobs[key]; !ok {
		return errors.New("Not existent job")
	}

	c.Jobs[key] = j
	return nil
}

// DeleteJob Deletes a job on memory
func (c *Dummy) DeleteJob(j *job.Job) error {
	c.jobsMutex.Lock()
	defer c.jobsMutex.Unlock()

	key := fmt.Sprintf(jobKeyFmt, j.ID)
	delete(c.Jobs, key)
	delete(c.Results, fmt.Sprintf(jobResultsKeyFmt, j.ID))
	// Don't return error if the job doesn't exists

	return nil
}

// JobsLength returns the number of jobs stored
func (c *Dummy) JobsLength() int {
	return len(c.Jobs)
}

// GetResults returns results from a job on memory
func (c *Dummy) GetResults(j *job.Job, low, high int) ([]*job.Result, error) {
	c.resultsMutex.Lock()
	defer c.resultsMutex.Unlock()

	results, ok := c.Results[fmt.Sprintf(jobResultsKeyFmt, j.ID)]

	if !ok {
		return nil, errors.New("Wrong job")
	}

	if results == nil {
		return nil, errors.New("Error retrieving jobs")
	}

	// High on top means all
	if high == 0 {
		high = len(results)
	}

	// Check indexes ok
	if low > high || low > len(results) || high > len(results) {
		return nil, errors.New("wrong parameters")
	}

	// Add 1 to match the indexes
	low++
	high++
	res := []*job.Result{}
	for i := low; i < high; i++ {
		res = append(res, results[fmt.Sprintf(resultKeyFmt, i)])
	}
	return res, nil
}

// GetResult returns a result from a job on memory
func (c *Dummy) GetResult(j *job.Job, id int) (*job.Result, error) {
	c.resultsMutex.Lock()
	defer c.resultsMutex.Unlock()

	results, ok := c.Results[fmt.Sprintf(jobResultsKeyFmt, j.ID)]
	if !ok {
		return nil, errors.New("Wrong job")
	}

	res, ok := results[fmt.Sprintf(resultKeyFmt, id)]
	if !ok {
		return nil, errors.New("Not existent result")
	}

	return res, nil
}

// SaveResult saves a result on a job in memory
func (c *Dummy) SaveResult(r *job.Result) error {
	c.resultsMutex.Lock()
	defer c.resultsMutex.Unlock()

	resultsKey := fmt.Sprintf(jobResultsKeyFmt, r.Job.ID)
	c.ResultsCounter[resultsKey]++
	r.ID = c.ResultsCounter[resultsKey]
	results, ok := c.Results[resultsKey]
	if !ok {
		results = map[string]*job.Result{}
		c.Results[resultsKey] = results
	}
	results[fmt.Sprintf(resultKeyFmt, r.ID)] = r

	return nil
}

// DeleteResult deletes a job on ajob in memory
func (c *Dummy) DeleteResult(r *job.Result) error {
	c.resultsMutex.Lock()
	defer c.resultsMutex.Unlock()

	results, ok := c.Results[fmt.Sprintf(jobResultsKeyFmt, r.Job.ID)]
	if !ok {
		return errors.New("Wrong job")
	}

	delete(results, fmt.Sprintf(resultKeyFmt, r.ID))

	// Don't return error if the result doens't exists
	return nil
}

// ResultsLength returns the number of results stored
func (c *Dummy) ResultsLength(j *job.Job) int {
	resultsKey := fmt.Sprintf(jobResultsKeyFmt, j.ID)
	if l, ok := c.Results[resultsKey]; ok {
		return len(l)
	}

	// No key, means 0 size
	return 0
}

// SaveAuthenticationToken stores an authentication token on database
func (c *Dummy) SaveAuthenticationToken(token string) error {
	if token == "" {
		return errors.New("Wrong token")
	}
	c.TokenMutex.Lock()
	defer c.TokenMutex.Unlock()
	c.Tokens[token] = struct{}{}
	return nil
}

// DeleteAuthenticationToken Deletes an authentication token from database
func (c *Dummy) DeleteAuthenticationToken(token string) error {
	c.TokenMutex.Lock()
	defer c.TokenMutex.Unlock()
	delete(c.Tokens, token)
	return nil
}

// AuthenticationTokenExists Checks if an authentication token exists
func (c *Dummy) AuthenticationTokenExists(token string) bool {
	c.TokenMutex.Lock()
	defer c.TokenMutex.Unlock()
	if _, ok := c.Tokens[token]; ok {
		return true
	}
	return false
}
