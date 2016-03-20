package service

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

var (
	dummyKey = "123456789"
)

// AuthenticationHandler Checks the application security, Authorization header
// and let it pass if correct.
func (s *KhronosService) AuthenticationHandler(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check security only if enabled
		if !s.Config.APIDisableSecurity {
			// Check header
			auth := r.Header.Get("Authorization")
			// TODO: Remove dummy key
			key := dummyKey
			if auth != fmt.Sprintf("Bearer %s", key) {
				logrus.Debugf("Forbidden access with header '%s'", auth)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		f.ServeHTTP(w, r)
	})
}
