package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/slok/khronos/service/validate"
)

const (
	errorRetrievingAllJobsMsg    = "Error retrieving all jobs"
	errorCreatingJobMsg          = "Error creating job"
	errorDeletingJobMsg          = "Error deleting job"
	errorDeletingResultMsg       = "Error deleting result"
	errorRetrievingJobMsg        = "Error retrieving job"
	errorRetrievingJobResultsMsg = "Error retrieving job results"
	wrongParamsMsg               = "Wrong params"
)

//Ping informs service is alive
func (s *KhronosService) Ping(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling ping endpoint")
	return http.StatusOK, "pong", nil
}

//GetJobs returns a list of jobs
func (s *KhronosService) GetJobs(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling GetAllJobs endpoint")

	// Page will set the offset
	p := r.URL.Query().Get("page")
	var err error
	page := 1
	if p != "" {
		if page, err = strconv.Atoi(p); err != nil {
			logrus.Errorf("error getting page querystring param: %v", err)
			return http.StatusInternalServerError, wrongParamsMsg, nil
		}
	}

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
	jid, _ := mux.Vars(r)["id"]
	logrus.Debug("Calling GetJob with id: %s", jid)

	jobID, err := strconv.Atoi(jid)
	if err != nil {
		logrus.Errorf("error getting job ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}

	j, err := s.Storage.GetJob(jobID)

	if err != nil {
		logrus.Errorf("Error retrieving job: %v", err)
		return http.StatusInternalServerError, errorRetrievingJobMsg, nil
	}

	return http.StatusOK, j, nil
}

// DeleteJob Deletes a job and its results
func (s *KhronosService) DeleteJob(r *http.Request) (int, interface{}, error) {
	// Get resul ID
	jid, _ := mux.Vars(r)["id"]
	logrus.Debug("Calling DeleteJob with id: %s", jid)

	jobID, err := strconv.Atoi(jid)
	if err != nil {
		logrus.Errorf("error getting job ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}

	j, err := s.Storage.GetJob(jobID)
	// No job, we are ok
	if err != nil {
		logrus.Errorf("Error retrieving job: %v", err)
		return http.StatusNoContent, nil, nil
	}

	if err := s.Storage.DeleteJob(j); err != nil {
		logrus.Errorf("error deleting job ID: %v", err)
		return http.StatusInternalServerError, errorDeletingJobMsg, nil
	}

	return http.StatusNoContent, nil, nil
}

// GetResults returns the jobs from an specific job
func (s *KhronosService) GetResults(r *http.Request) (int, interface{}, error) {
	// Get job ID
	jid, _ := mux.Vars(r)["jobID"]
	logrus.Debugf("Calling GetResults from jobid: %s", jid)

	// Get the job
	jobID, err := strconv.Atoi(jid)

	if err != nil {
		logrus.Errorf("error getting job ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}

	// Page will set the offset
	p := r.URL.Query().Get("page")
	page := 1
	if p != "" {
		if page, err = strconv.Atoi(p); err != nil {
			logrus.Errorf("error getting page querystring param: %v", err)
			return http.StatusInternalServerError, wrongParamsMsg, nil
		}
	}

	j, err := s.Storage.GetJob(jobID)
	if err != nil {
		return http.StatusNoContent, nil, nil
	}

	// Get job results
	results, err := s.Storage.GetResults(j, 0, 0)
	if err != nil {
		logrus.Errorf("Error retrieving job results: %v", err)
		return http.StatusInternalServerError, errorRetrievingJobResultsMsg, nil
	}

	return http.StatusOK, results, nil
}

// GetResult returns a single result by id
func (s *KhronosService) GetResult(r *http.Request) (int, interface{}, error) {
	// Get resul ID
	jid, _ := mux.Vars(r)["jobID"]
	jobID, err := strconv.Atoi(jid)
	if err != nil {
		logrus.Errorf("error getting job ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}

	rid, _ := mux.Vars(r)["resultID"]
	resultID, err := strconv.Atoi(rid)
	if err != nil {
		logrus.Errorf("error getting result ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}
	logrus.Debugf("Calling GetResult with id: %d from job '%d'", resultID, jobID)

	j, err := s.Storage.GetJob(jobID)
	if err != nil {
		logrus.Errorf("Error retrieving Job: %v", err)
		return http.StatusInternalServerError, errorRetrievingJobMsg, nil
	}
	result, err := s.Storage.GetResult(j, resultID)
	if err != nil {
		logrus.Errorf("Error retrieving job result: %v", err)
		return http.StatusInternalServerError, errorRetrievingJobResultsMsg, nil
	}

	return http.StatusOK, result, nil
}

// DeleteResult deletes a result
func (s *KhronosService) DeleteResult(r *http.Request) (int, interface{}, error) {
	// Get resul ID
	jid, _ := mux.Vars(r)["jobID"]
	jobID, err := strconv.Atoi(jid)
	if err != nil {
		logrus.Errorf("error getting job ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}

	rid, _ := mux.Vars(r)["resultID"]
	resultID, err := strconv.Atoi(rid)
	if err != nil {
		logrus.Errorf("error getting result ID: %v", err)
		return http.StatusInternalServerError, wrongParamsMsg, nil
	}
	logrus.Debugf("Calling DeleteResult with id: %d from job '%d'", resultID, jobID)

	j, err := s.Storage.GetJob(jobID)
	if err != nil {
		logrus.Errorf("Error deleting Job: %v", err)
		return http.StatusInternalServerError, errorDeletingResultMsg, nil
	}
	result, err := s.Storage.GetResult(j, resultID)
	// If no result then is ok
	if err != nil {
		logrus.Errorf("Error deleting job result: %v", err)
		return http.StatusNoContent, result, nil
	}

	if err := s.Storage.DeleteResult(result); err != nil {
		logrus.Errorf("Error deleting job result: %v", err)
		return http.StatusInternalServerError, errorDeletingResultMsg, nil
	}

	return http.StatusNoContent, result, nil
}
