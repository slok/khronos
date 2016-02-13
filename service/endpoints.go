package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

const (
	errorRetrievingAllJobsMsg = "Error retrieving all jobs"
)

//Ping informs service is alive
func (s *KhronosService) Ping(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling ping endpoint")
	return http.StatusOK, "pong", nil
}

//GetAllJobs returns a list of jobs
func (s *KhronosService) GetAllJobs(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling GetAllJobs endpoint")

	jobs, err := s.Client.GetHTTPJobs()

	if err != nil {
		logrus.Errorf("Error retrieving all jobs: %v", err)
		return http.StatusInternalServerError, errorRetrievingAllJobsMsg, nil
	}

	return http.StatusOK, jobs, nil
}
