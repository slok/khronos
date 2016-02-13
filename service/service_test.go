package service

import (
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
	testConfig        = &config.AppConfig{}
	testStorageClient = storage.NewDummy()
)

func TestPing(t *testing.T) {

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

	// Testing data
	tests := []struct {
		givenURI       string
		givenHTTPJobs  map[string]*job.HTTPJob
		wantCode       int
		wantBodyLength int
	}{
		{
			givenURI:       "/api/v1/jobs",
			givenHTTPJobs:  make(map[string]*job.HTTPJob),
			wantCode:       http.StatusOK,
			wantBodyLength: 0,
		},
		{
			givenURI: "/api/v1/jobs",
			givenHTTPJobs: map[string]*job.HTTPJob{
				"job:test1": &job.HTTPJob{Job: job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true}, URL: &url.URL{}},
				"job:test2": &job.HTTPJob{Job: job.Job{ID: 2, Name: "test2", Description: "test2", When: "0 30 * * * *", Active: true}, URL: &url.URL{}},
			},
			wantCode:       http.StatusOK,
			wantBodyLength: 2,
		},
	}

	// Tests
	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.HTTPJobs = test.givenHTTPJobs
		testStorageClient.HTTPJobCounter = len(test.givenHTTPJobs)

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

		var got []*job.HTTPJob
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error(err)
		}

		if len(got) != test.wantBodyLength {
			t.Errorf("Expected length '%d'. Got '%d' instead ", test.wantBodyLength, len(got))
		}
	}

}
