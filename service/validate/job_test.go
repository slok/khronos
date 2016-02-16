package validate

import (
	"errors"
	"reflect"
	"testing"
)

func TestHTTPJobValidatorJSON(t *testing.T) {
	tests := []struct {
		givenJSON     string
		wantValidator HTTPJobValidator
	}{
		{
			givenJSON: `{"active": true, "description": "Simple hello world", "url": "http://crons.test.com/hello-world", "when": "@daily", "name": "hello-world"}`,
			wantValidator: HTTPJobValidator{
				Name:        "hello-world",
				Description: "Simple hello world",
				When:        "@daily",
				Active:      true,
				URL:         "http://crons.test.com/hello-world",
				Errors:      nil,
			},
		},
	}

	for _, test := range tests {
		v, err := NewHTTPJobValidatorFromJSON(test.givenJSON)

		if err != nil {
			t.Error(err)
		}

		v.Errors = nil // Handy for tests
		if !reflect.DeepEqual(*v, test.wantValidator) {
			t.Errorf("Validators are not equal; expected %v; got %v", test.wantValidator, *v)
		}
	}

}

func TestHTTPJobValidatorValidation(t *testing.T) {
	tests := []struct {
		givenValidator *HTTPJobValidator
		wantError      bool
		wantErrors     []error
	}{
		{
			givenValidator: &HTTPJobValidator{
				Name:        "hello-world",
				Description: "Simple hello world",
				When:        "@daily",
				Active:      true,
				URL:         "http://crons.test.com/hello-world",
			},
			wantError:  false,
			wantErrors: []error{},
		},
		{
			givenValidator: &HTTPJobValidator{},
			wantError:      true,
			wantErrors: []error{
				errors.New("Name is required"),
				errors.New("When is required"),
				errors.New("URL is required"),
				errors.New("When is not a valid cron"),
			},
		}, {
			givenValidator: &HTTPJobValidator{
				Name: "hello-world",
				When: "@daily",
				URL:  "http://crons.test.com/hello-world",
			},
			wantError:  false,
			wantErrors: []error{},
		}, {
			givenValidator: &HTTPJobValidator{
				Description: "Simple hello world",
				When:        "@daily",
				Active:      true,
			},
			wantError: true,
			wantErrors: []error{
				errors.New("Name is required"),
				errors.New("URL is required"),
			},
		},
	}

	for _, test := range tests {
		err := test.givenValidator.Validate()
		if !test.wantError && err != nil {
			t.Error("Didn't expect error")
		} else if test.wantError && err == nil {
			t.Error("Excepted error")
		}

		if !reflect.DeepEqual(test.givenValidator.Errors, test.wantErrors) {
			t.Errorf("Errors are not equal; expected %v; got %v", test.wantErrors, test.givenValidator.Errors)
		}
	}

}
