package service

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/schedule"
	"github.com/slok/khronos/storage"
)

func TestAuthenticationMiddleware(t *testing.T) {

	// These are the valid tokens to make the requests
	validTokens := map[string]struct{}{
		"123456789": struct{}{},
	}

	tests := []struct {
		GivenAuthHeader string
		WantCode        int
		SecurityEnabled bool
	}{
		{GivenAuthHeader: "", WantCode: http.StatusForbidden, SecurityEnabled: true},
		{GivenAuthHeader: "", WantCode: http.StatusOK, SecurityEnabled: false},
		{GivenAuthHeader: "Bearer ", WantCode: http.StatusForbidden, SecurityEnabled: true},
		{GivenAuthHeader: "Bearer 987654321", WantCode: http.StatusForbidden, SecurityEnabled: true},
		{GivenAuthHeader: "Bearer 123456789", WantCode: http.StatusOK, SecurityEnabled: true},
		{GivenAuthHeader: "123456789", WantCode: http.StatusForbidden, SecurityEnabled: true},
		{GivenAuthHeader: "Bearer  123456789", WantCode: http.StatusForbidden, SecurityEnabled: true},
	}

	for _, test := range tests {
		testConfig := config.NewAppConfig(os.Getenv(config.KhronosConfigFileKey))
		// Disable or enable security based on the test
		testConfig.APIDisableSecurity = !test.SecurityEnabled
		testStorageClient := storage.NewDummy()
		// mock custom auth tokens on database
		testStorageClient.Tokens = validTokens
		testSchedulerClient := schedule.NewDummyCron(testConfig, testStorageClient, 0, "OK")

		s := &KhronosService{
			Config:  testConfig,
			Storage: testStorageClient,
			Cron:    testSchedulerClient,
		}

		// Create a testing request
		r, _ := http.NewRequest("GET", "/", nil)
		r.Header.Add("Authorization", test.GivenAuthHeader)
		w := httptest.NewRecorder()

		// Apply middleware to test
		s.AuthenticationHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).ServeHTTP(w, r)

		// Check testing ok
		if w.Code != test.WantCode {
			t.Errorf("Authorization middleware error on response status code; expected %d, got %d instead", test.WantCode, w.Code)
		}
	}

}
