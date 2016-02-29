package service

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/slok/khronos/service/validate"
)

const (
	errorRetrievingAllJobsMsg = "Error retrieving all jobs"
	errorCreatingJobMsg       = "Error creating job"
)

//Ping informs service is alive
func (s *KhronosService) Ping(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling ping endpoint")
	return http.StatusOK, "pong", nil
}

//GetJobs returns a list of jobs
func (s *KhronosService) GetJobs(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling GetAllJobs endpoint")

	jobs, err := s.Storage.GetJobs(0, 0)

	if err != nil {
		logrus.Errorf("Error retrieving all jobs: %v", err)
		return http.StatusInternalServerError, errorRetrievingAllJobsMsg, nil
	}

	return http.StatusOK, jobs, nil
}

//CreateNewJob Creates and registers a new job
func (s *KhronosService) CreateNewJob(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling CreateNewJob endpoint")
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// Unmarshall tje received json
	v, err := validate.NewJobValidatorFromJSON(string(b))
	if err != nil {
		logrus.Errorf("Error unmarshalling json: %v", err)
		return http.StatusInternalServerError, errorCreatingJobMsg, nil
	}

	// Validate received json
	if err = v.Validate(); err != nil {
		result := map[string][]string{
			"errors": []string{},
		}

		for _, e := range v.Errors {
			result["errors"] = append(result["errors"], fmt.Sprintf("%v", e))
		}

		return http.StatusBadRequest, result, nil

	}

	// Store the received json
	j, err := v.Instance()
	if err != nil {
		logrus.Errorf("Error Creating valid job instance: %v", err)
		return http.StatusInternalServerError, errorCreatingJobMsg, nil

	}
	err = s.Storage.SaveJob(j)
	if err != nil {
		logrus.Errorf("Error storing job: %v", err)
		return http.StatusInternalServerError, errorCreatingJobMsg, nil

	}

	// Register a new cron job!
	s.Cron.RegisterCronJob(j)

	return http.StatusCreated, j, nil
}

// GetJob returns a single job by id
func (s *KhronosService) GetJob(r *http.Request) (int, interface{}, error) {
	// Get resul ID
	jobID, _ := mux.Vars(r)["id"]
	logrus.Debug("Calling GetJob with id: %s", jobID)

	return http.StatusNotImplemented, nil, nil
}

// GetResults returns the jobs from an specific job
func (s *KhronosService) GetResults(r *http.Request) (int, interface{}, error) {
	// Get job ID
	jobID := r.URL.Query().Get("job")
	logrus.Debug("Calling GetResults from jobid: %s", jobID)

	return http.StatusNotImplemented, nil, nil
}

// GetResult returns a single result by id
func (s *KhronosService) GetResult(r *http.Request) (int, interface{}, error) {
	// Get resul ID
	resultID, _ := mux.Vars(r)["id"]
	logrus.Debug("Calling GetResult with id: %s", resultID)

	return http.StatusNotImplemented, nil, nil
}
