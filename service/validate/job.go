package validate

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/job"
)

// HTTPJobValidator implements the requirements of a validator in order to
// be able to create correct HTTPJobs
type HTTPJobValidator struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	When        string `json:"when"`
	Active      bool   `json:"active"`
	URL         string `json:"url"`

	// Errors after validating the instance
	Errors []error
}

// NewHTTPJobValidatorFromJSON creates a validator from a json
func NewHTTPJobValidatorFromJSON(j string) (v *HTTPJobValidator, err error) {
	v = &HTTPJobValidator{}
	err = json.Unmarshal([]byte(j), v)
	logrus.Debug("Created HTTPJob validator from json")
	return
}

// Validate validates the validator and creates teh correct instance
func (v *HTTPJobValidator) Validate() error {
	logrus.Debugf("Validating job '%s'", v.Name)

	// Flush previous errors
	v.Errors = []error{}

	// Check required fields
	if v.Name == "" {
		v.Errors = append(v.Errors, errors.New("Name is required"))
	}
	if v.When == "" {
		v.Errors = append(v.Errors, errors.New("When is required"))
	}
	if v.URL == "" {
		v.Errors = append(v.Errors, errors.New("URL is required"))
	}

	// Check valid cron
	if err := ValidCron(v.When); err != nil {
		v.Errors = append(v.Errors, errors.New("When is not a valid cron"))
	}

	if len(v.Errors) > 0 {
		return errors.New("Not valid HTTPJob")
	}

	return nil
}

// Instance returns a valid instance
func (v *HTTPJobValidator) Instance() (j *job.HTTPJob, err error) {
	if err = v.Validate(); err != nil {
		return
	}
	u, err := url.ParseRequestURI(v.URL)
	if err != nil {
		return
	}

	return &job.HTTPJob{
		Job: job.Job{
			Name:        v.Name,
			Description: v.Description,
			When:        v.When,
			Active:      v.Active,
		},
		URL: u,
	}, nil
}
