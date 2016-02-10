package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NYTimes/gizmo/server"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/storage"
)

var (
	testConfig        = &config.AppConfig{}
	testStorageClient = storage.NewNil()
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
