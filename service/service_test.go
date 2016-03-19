package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/NYTimes/gizmo/server"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/job"
	"github.com/slok/khronos/schedule"
	"github.com/slok/khronos/storage"
)

var (
	testConfig = config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
)

func TestPing(t *testing.T) {
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

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
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
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

func TestGetJobs(t *testing.T) {
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

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
				"job:1": &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}},
				"job:2": &job.Job{ID: 2, Name: "test2", Description: "test2", When: "0 30 * * * *", Active: true, URL: &url.URL{}},
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
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
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

func TestGetJobsPaginated(t *testing.T) {
	jobs := make(map[string]*job.Job)
	totalJobs := 27
	pageSize := 5
	testStorageClient := storage.NewDummy()
	// Custom pagination
	paginationTestConfig := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	paginationTestConfig.APIResourcesPerPage = pageSize
	testCronEngine := schedule.NewDummyCron(paginationTestConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Create our custom dummy job database
	for i := 1; i <= totalJobs; i++ {
		k := fmt.Sprintf("job:%d", i)
		v := &job.Job{ID: i, Name: fmt.Sprintf("test%d", i), When: "@daily", URL: &url.URL{}}
		jobs[k] = v
	}

	// Testing data
	tests := []struct {
		givenURI   string
		wantJobIDs []int
	}{
		{
			givenURI:   "/api/v1/jobs",
			wantJobIDs: []int{1, 2, 3, 4, 5},
		},
		{
			givenURI:   "/api/v1/jobs?page=0",
			wantJobIDs: []int{1, 2, 3, 4, 5},
		},
		{
			givenURI:   "/api/v1/jobs?page=1",
			wantJobIDs: []int{1, 2, 3, 4, 5},
		},
		{
			givenURI:   "/api/v1/jobs?page=2",
			wantJobIDs: []int{6, 7, 8, 9, 10},
		},
		{
			givenURI:   "/api/v1/jobs?page=3",
			wantJobIDs: []int{11, 12, 13, 14, 15},
		},
		{
			givenURI:   "/api/v1/jobs?page=4",
			wantJobIDs: []int{16, 17, 18, 19, 20},
		},
		{
			givenURI:   "/api/v1/jobs?page=5",
			wantJobIDs: []int{21, 22, 23, 24, 25},
		},
		{
			givenURI:   "/api/v1/jobs?page=6",
			wantJobIDs: []int{26, 27},
		},
		{
			givenURI:   "/api/v1/jobs?page=7",
			wantJobIDs: []int{},
		},
	}

	// Tests
	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Jobs = jobs
		testStorageClient.JobCounter = totalJobs

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  paginationTestConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)

		var got []*job.Job
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error(err)
		}

		// Check length
		if len(got) != len(test.wantJobIDs) {
			t.Errorf("Expected length '%d'. Got '%d' instead ", len(test.wantJobIDs), len(got))
		}

		// Check IDs ok (should be in order)
		for k, i := range test.wantJobIDs {
			if got[k].ID != i {
				t.Errorf("Expected job id '%d'. Got '%d' instead ", i, got[k].ID)
			}
		}
	}
}

func TestGetJob(t *testing.T) {
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	j2 := &job.Job{ID: 2, Name: "test2", Description: "test2", When: "0 30 * * * *", Active: false, URL: &url.URL{}}
	// Testing data
	tests := []struct {
		givenURI  string
		givenJobs map[string]*job.Job
		wantCode  int
		wantJob   *job.Job
	}{
		{
			givenURI:  "/api/v1/jobs/1",
			givenJobs: make(map[string]*job.Job),
			wantCode:  http.StatusInternalServerError,
		},
		{
			givenURI: "/api/v1/jobs/2",
			givenJobs: map[string]*job.Job{
				"job:1": &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}},
				"job:2": j2, // The one to check,
				"job:3": &job.Job{ID: 3, Name: "test3", Description: "test3", When: "*/10 30 * 4 * 1", Active: true, URL: &url.URL{}},
			},
			wantCode: http.StatusOK,
			wantJob:  j2,
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
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		// Only check when ok
		if w.Code == http.StatusOK {
			var got *job.Job
			err := json.NewDecoder(w.Body).Decode(&got)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, test.wantJob) {
				t.Errorf("Expected job '%#v'. Got '%#v' instead ", test.wantJob, got)
			}
		}
	}
}

func TestDeleteJob(t *testing.T) {
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Testing data
	j := &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}}
	tests := []struct {
		givenURI     string
		givenJobs    map[string]*job.Job
		givenResults map[string]map[string]*job.Result
		wantCode     int
		wantJob      *job.Job
	}{
		{
			givenURI:     "/api/v1/jobs/1",
			givenJobs:    make(map[string]*job.Job),
			givenResults: make(map[string]map[string]*job.Result),
			wantCode:     http.StatusNoContent,
		},
		{
			givenURI: "/api/v1/jobs/1",
			givenJobs: map[string]*job.Job{
				"job:1": j,
			},
			givenResults: map[string]map[string]*job.Result{
				"job:1:results": map[string]*job.Result{
					"result:1": &job.Result{ID: 1, Job: j, Out: "test1", Status: job.ResultOK, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:2": &job.Result{ID: 2, Job: j, Out: "test1", Status: job.ResultError, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:3": &job.Result{ID: 3, Job: j, Out: "test1", Status: job.ResultInternalError, Start: time.Now().UTC(), Finish: time.Now().UTC()},
				},
			},
			wantCode: http.StatusNoContent,
		},
	}

	// Tests
	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Jobs = test.givenJobs
		testStorageClient.JobCounter = len(test.givenJobs)
		testStorageClient.Results = test.givenResults

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("DELETE", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		// Check no job
		if _, ok := test.givenJobs["job:1"]; ok {
			t.Error("Job should be deleted, job present on database")
		}

		// Check no job results
		if _, ok := test.givenResults["job:1:results"]; ok {
			t.Error("Job results should be deleted, job results present on database")
		}

	}
}

func TestCreateNewJob(t *testing.T) {
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

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
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
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

func TestGetResults(t *testing.T) {
	j := &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}}
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Testing data
	tests := []struct {
		givenURI     string
		givenResults map[string]map[string]*job.Result
		givenJobs    map[string]*job.Job
		wantCode     int
	}{
		{
			givenURI:     "/api/v1/jobs/1/results",
			givenResults: make(map[string]map[string]*job.Result),
			givenJobs:    map[string]*job.Job{"job:1": j},
			wantCode:     http.StatusInternalServerError,
		},
		{
			givenURI: "/api/v1/jobs/1/results",
			givenResults: map[string]map[string]*job.Result{
				"job:1:results": map[string]*job.Result{
					"result:1": &job.Result{ID: 1, Job: j, Out: "test1", Status: job.ResultOK, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:2": &job.Result{ID: 2, Job: j, Out: "test1", Status: job.ResultError, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:3": &job.Result{ID: 3, Job: j, Out: "test1", Status: job.ResultInternalError, Start: time.Now().UTC(), Finish: time.Now().UTC()},
				},
			},
			givenJobs: map[string]*job.Job{"job:1": j},
			wantCode:  http.StatusOK,
		},
	}

	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Results = test.givenResults
		testStorageClient.Jobs = test.givenJobs

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		// Only check in good results
		if w.Code == http.StatusOK {
			b, err := ioutil.ReadAll(w.Body)

			if err != nil {
				t.Errorf("Error reading result body: %v", err)
			}
			var gotRes []*job.Result
			if err := json.Unmarshal(b, &gotRes); err != nil {
				t.Errorf("Error unmarshaling: %v", err)
			}
			rs, ok := test.givenResults["job:1:results"]
			if !ok {
				t.Errorf("Error getting results: %v", err)
			}
			if len(gotRes) != len(rs) {
				t.Errorf("Expected len '%d'. Got '%d' instead ", len(rs), len(gotRes))
			}
		}
	}
}

func TestGetResult(t *testing.T) {
	j := &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}}
	res3 := &job.Result{ID: 3, Job: nil, Out: "test1", Status: job.ResultInternalError, Start: time.Now().UTC(), Finish: time.Now().UTC()}
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Testing data
	tests := []struct {
		givenURI     string
		givenResults map[string]map[string]*job.Result
		givenJobs    map[string]*job.Job
		wantCode     int
		wantResult   *job.Result
	}{
		{
			givenURI:     "/api/v1/jobs/1/results/1",
			givenResults: make(map[string]map[string]*job.Result),
			givenJobs:    make(map[string]*job.Job),
			wantCode:     http.StatusInternalServerError,
		},
		{
			givenURI:     "/api/v1/jobs/1/results/1",
			givenResults: make(map[string]map[string]*job.Result),
			givenJobs:    map[string]*job.Job{"job:1": j},
			wantCode:     http.StatusInternalServerError,
		},
		{
			givenURI: "/api/v1/jobs/1/results/3",
			givenResults: map[string]map[string]*job.Result{
				"job:1:results": map[string]*job.Result{
					"result:1": &job.Result{ID: 1, Job: j, Out: "test1", Status: job.ResultOK, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:2": &job.Result{ID: 2, Job: j, Out: "test1", Status: job.ResultError, Start: time.Now().UTC(), Finish: time.Now().UTC()},
					"result:3": res3, // the one to check
				},
			},
			givenJobs:  map[string]*job.Job{"job:1": j},
			wantCode:   http.StatusOK,
			wantResult: res3,
		},
	}

	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Results = test.givenResults
		testStorageClient.Jobs = test.givenJobs

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		// Only check in good results
		if w.Code == http.StatusOK {
			b, err := ioutil.ReadAll(w.Body)

			if err != nil {
				t.Errorf("Error reading result body: %v", err)
			}
			var gotRes *job.Result
			if err := json.Unmarshal(b, &gotRes); err != nil {
				t.Errorf("Error unmarshaling: %v", err)
			}
			// Fix jobs for the deep equal
			gotRes.Job = nil

			if !reflect.DeepEqual(gotRes, test.wantResult) {
				t.Errorf("Expected Result '%#v'. Got '%#v' instead ", test.wantResult, gotRes)
			}
		}

	}
}

func TestDeleteResult(t *testing.T) {
	j := &job.Job{ID: 1, Name: "test1", Description: "test1", When: "@daily", Active: true, URL: &url.URL{}}
	testStorageClient := storage.NewDummy()
	testCronEngine := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Testing data
	tests := []struct {
		givenURI          string
		givenResults      map[string]map[string]*job.Result
		givenJobs         map[string]*job.Job
		wantCode          int
		wantDeletedResult string
	}{
		{
			givenURI:     "/api/v1/jobs/1/results/1",
			givenResults: make(map[string]map[string]*job.Result),
			givenJobs:    make(map[string]*job.Job),
			wantCode:     http.StatusInternalServerError,
		},
		{
			givenURI:     "/api/v1/jobs/1/results/1",
			givenResults: make(map[string]map[string]*job.Result),
			givenJobs:    map[string]*job.Job{"job:1": j},
			wantCode:     http.StatusNoContent,
		},
		{
			givenURI: "/api/v1/jobs/1/results/1",
			givenResults: map[string]map[string]*job.Result{
				"job:1:results": map[string]*job.Result{
					"result:1": &job.Result{ID: 1, Job: j, Out: "test1", Status: job.ResultOK, Start: time.Now().UTC(), Finish: time.Now().UTC()},
				},
			},
			givenJobs:         map[string]*job.Job{"job:1": j},
			wantCode:          http.StatusNoContent,
			wantDeletedResult: "result:1",
		},
	}

	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Results = test.givenResults
		testStorageClient.Jobs = test.givenJobs

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("DELETE", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)
		if w.Code != test.wantCode {
			t.Errorf("Expected response code '%d'. Got '%d' instead ", test.wantCode, w.Code)
		}

		rs, _ := testStorageClient.Results["job:1:results"]
		if _, ok := rs[test.wantDeletedResult]; ok {
			t.Error("Result should be deleted; present on database")
		}
	}
}

func TestGetResultsPaginated(t *testing.T) {
	results := map[string]map[string]*job.Result{
		"job:1:results": map[string]*job.Result{},
	}
	totalResults := 54
	pageSize := 7
	testStorageClient := storage.NewDummy()
	// Custom pagination
	paginationTestConfig := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
	paginationTestConfig.APIResourcesPerPage = pageSize
	testCronEngine := schedule.NewDummyCron(paginationTestConfig, testStorageClient, 0, "OK")
	testCronEngine.Start(nil)

	// Create our custom dummy results database
	j := &job.Job{ID: 1, Name: "test1", When: "@daily", Active: true, URL: &url.URL{}}
	for i := 1; i <= totalResults; i++ {
		v := &job.Result{ID: i, Job: j, Out: fmt.Sprintf("test%d", i), Status: job.ResultInternalError, Start: time.Now().UTC(), Finish: time.Now().UTC()}
		k := fmt.Sprintf("result:%d", i)
		results["job:1:results"][k] = v
	}

	// Testing data
	tests := []struct {
		givenURI      string
		wantResultIDs []int
	}{
		{
			givenURI:      "/api/v1/jobs/1/results",
			wantResultIDs: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=0",
			wantResultIDs: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=1",
			wantResultIDs: []int{1, 2, 3, 4, 5, 6, 7},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=2",
			wantResultIDs: []int{8, 9, 10, 11, 12, 13, 14},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=3",
			wantResultIDs: []int{15, 16, 17, 18, 19, 20, 21},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=4",
			wantResultIDs: []int{22, 23, 24, 25, 26, 27, 28},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=5",
			wantResultIDs: []int{29, 30, 31, 32, 33, 34, 35},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=6",
			wantResultIDs: []int{36, 37, 38, 39, 40, 41, 42},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=7",
			wantResultIDs: []int{43, 44, 45, 46, 47, 48, 49},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=8",
			wantResultIDs: []int{50, 51, 52, 53, 54},
		},
		{
			givenURI:      "/api/v1/jobs/1/results?page=9",
			wantResultIDs: []int{},
		},
	}

	// Tests
	for _, test := range tests {
		// Set our dummy 'database' on the storage client
		testStorageClient.Results = results
		testStorageClient.Jobs = map[string]*job.Job{"job:1": j}

		// Create a testing server
		testServer := server.NewSimpleServer(nil)

		// Register our service on the server (we don't need configuration for this service)
		testServer.Register(&KhronosService{
			Config:  paginationTestConfig,
			Storage: testStorageClient,
			Cron:    testCronEngine,
		})

		// Create request and a test recorder
		r, _ := http.NewRequest("GET", test.givenURI, nil)
		w := httptest.NewRecorder()
		testServer.ServeHTTP(w, r)

		var got []*job.Result
		err := json.NewDecoder(w.Body).Decode(&got)
		if err != nil {
			t.Error(err)
		}

		// Check length
		if len(got) != len(test.wantResultIDs) {
			t.Errorf("Expected length '%d'. Got '%d' instead ", len(test.wantResultIDs), len(got))
		}

		// Check IDs ok (should be in order)
		for k, i := range test.wantResultIDs {
			if got[k].ID != i {
				t.Errorf("Expected result id '%d'. Got '%d' instead ", i, got[k].ID)
			}
		}
	}
}
