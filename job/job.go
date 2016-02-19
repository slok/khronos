package job

import (
	"encoding/json"
	"net/url"
)

// Job is the unit of job to be executed periodically making an HTTP call
type Job struct {
	ID          int
	Name        string
	Description string
	When        string
	Active      bool
	URL         *url.URL
}

// MarshalJSON is a custom json marshaller for Job
func (j *Job) MarshalJSON() ([]byte, error) {
	// Alias is a custom type to inherint all the properties of Job but not the methods
	type Alias Job

	return json.Marshal(struct {
		*Alias
		URL string
	}{
		Alias: (*Alias)(j),
		URL:   j.URL.String(),
	})
}

// UnmarshalJSON is a custom json unmarshaller for Job
func (j *Job) UnmarshalJSON(data []byte) error {
	// Alias is a custom type to inherint all the properties of Job but not the methods
	type Alias Job

	aux := struct {
		*Alias
		URL string
	}{
		Alias: (*Alias)(j),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	sURL, err := url.Parse(aux.URL)
	if err != nil {
		return err
	}

	j.URL = sURL
	return nil
}
