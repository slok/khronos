package validate

import (
	"encoding/json"

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

// NewHTTPJobFromJSON creates a validator from a json
func NewHTTPJobFromJSON(j string) (v *HTTPJobValidator, err error) {
	v = &HTTPJobValidator{}
	err = json.Unmarshal([]byte(j), v)
	logrus.Debug("Created HTTPJob validator from json")
	return
}

// Validate validates the validator and creates teh correct instance
func (v *HTTPJobValidator) Validate() (j job.HTTPJob, err error) {
	logrus.Debug("Validating job '%s'", v.Name)
	logrus.Error("Not implemented")
	return
}
