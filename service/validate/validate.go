// Package validate implements utilities to validate properties and objects
package validate

import (
	"errors"

	"github.com/robfig/cron"
)

const (
	notValidCron = "Invalid cron syntax"
	required     = "Required value"
)

// Validator validates properties of a validator object
type Validator interface {
	Validate() error
}

// ValidCron checks if the syntax of value is valid for a cron job
func ValidCron(value string) error {
	_, err := cron.Parse(value)
	if err != nil {
		return errors.New(notValidCron)
	}
	return nil
}
