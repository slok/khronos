package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/NYTimes/gizmo/server"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/job"
	"github.com/slok/khronos/storage"
)

var (
	testConfig = &config.AppConfig{}
)

func TestPing(t *testing.T) {
	testStorageClient := storage.NewDummy()

	// Testing data
	tests := []struct {
		givenURI string
		wantCode int
		wantBody interface{}
	}{
		{
			givenURI: "/api/v1/ping",
			wantCode: http.StatusOK,
			wantBody: "\"pong\"\n",
		},
	}

	// Tests
	for _, test := range tests {

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config: testConfig,
			Client: testStorageClient,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		got, _ := ioutil.ReadAll(w.Body)
		if string(got) != test.wantBody {
			t.Errorf("Expected body '%s'. Got '%s' instead ", test.wantBody, string(got))
		}
	}

}

func TestGetAllJobs(t *testing.T) {
	testStorageClient := storage.NewDummy()

	// Testing data
	tests := []struct {
		givenURI       string
		givenJobs      map[string]*job.Job
		wantCode       int
		wantBodyLength int
	}{
		{
			givenURI:       "/api/v1/jobs",
			givenJobs:      make(map[string]*job.Job),
			wantCode:       http.StatusOK,
			wantBodyLength: 0,
		},
		{
			givenURI: "/api/v1/jobs",
			givenJobs: map[string]*job.Job{
				"job:test1": &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}},
				"job:test2": &job.Job{ID: 2, Name: "test2", Description: "test2", When: "0 30 * * * *", Active: true, URL: &url.URL{}},
			},
			wantCode:       http.StatusOK,
			wantBodyLength: 2,
		},
	}

	// Tests
	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Jobs = test.givenJobs
		testStorageClient.JobCounter = len(test.givenJobs)

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config: testConfig,
			Client: testStorageClient,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		var got []*job.Job
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error(err)
		}

		if len(got) != test.wantBodyLength {
			t.Errorf("Expected length '%d'. Got '%d' instead ", test.wantBodyLength, len(got))
		}
	}
}

func TestCreateNewJob(t *testing.T) {
	testStorageClient := storage.NewDummy()

	// Testing data
	tests := []struct {
		givenURI    string
		givenBody   string
		givenJobs   map[string]*job.Job
		wantCode    int
		wantBody    string
		wantJobslen int
	}{
		{
			givenURI:    "/api/v1/jobs",
			givenBody:   `{"active": true, "description": "Simple hello world", "url": "http://crons.test.com/hello-world", "when": "@daily", "name": "hello-world"}`,
			givenJobs:   make(map[string]*job.Job),
			wantCode:    http.StatusCreated,
			wantJobslen: 1,
		},

		{
			givenURI:    "/api/v1/jobs",
			givenBody:   `{"active": true, "description": "Simple hello world", "url": "http://crons.test.com/hello-world", "when": "@daily"}`,
			givenJobs:   make(map[string]*job.Job),
			wantCode:    http.StatusBadRequest,
			wantJobslen: 0,
		},
	}

	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Jobs = test.givenJobs
		testStorageClient.JobCounter = len(test.givenJobs)

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config: testConfig,
			Client: testStorageClient,
		})

		b := bytes.NewReader([]byte(test.givenBody))
		r, _ := http.NewRequest("POST", test.givenURI, b)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)

		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}
		if len(testStorageClient.Jobs) != test.wantJobslen {
			t.Errorf("Expected len '%d'. Got '%d' instead ", len(testStorageClient.Jobs), test.wantJobslen)
		}
	}
}
