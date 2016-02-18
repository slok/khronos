package job

import (
	"encoding/json"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

// Scheduler interface implements the Schedule method wich will schedule a job
// to be executed in the future
type Scheduler interface {
	Schedule(c cron.Cron) error
}

// Job is the most basic unit of job without execution
type Job struct {
	ID          int
	Name        string
	Description string
	When        string
	Active      bool
}

// HTTPJob is the unit of job to be executed periodically making an HTTP call
type HTTPJob struct {
	Job
	URL *url.URL
	// TODO: method
	// TODO: body
	// TODO: headers
}

// Schedule schedules an HTTP call
func (h *HTTPJob) Schedule(c cron.Cron) error {
	logrus.Debug("Schedule %s job on %s with the url $v", h.Name, h.When, *h.URL)
	return nil
}

// MarshalJSON is a custom json marshaller for HTTPJob
func (h *HTTPJob) MarshalJSON() ([]byte, error) {
	// Alias is a custom type to inherint all the properties of HTTPJob but not the methods
	type Alias HTTPJob

	return json.Marshal(struct {
		*Alias
		URL string
	}{
		Alias: (*Alias)(h),
		URL:   h.URL.String(),
	})
}

// UnmarshalJSON is a custom json unmarshaller for HTTPJob
func (h *HTTPJob) UnmarshalJSON(data []byte) error {
	// Alias is a custom type to inherint all the properties of HTTPJob but not the methods
	type Alias HTTPJob

	aux := struct {
		*Alias
		URL string
	}{
		Alias: (*Alias)(h),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	sURL, err := url.Parse(aux.URL)
	if err != nil {
		return err
	}

	h.URL = sURL
	return nil
}
