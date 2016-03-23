package service

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/Sirupsen/logrus"
)

var (
	tokenRe = regexp.MustCompile(`Bearer (\w+)`)
)

// AuthenticationHandler Checks the application security, Authorization header
// and let it pass if correct.
func (s *KhronosService) AuthenticationHandler(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check security only if enabled
		if !s.Config.APIDisableSecurity {
			// Check header
			auth := r.Header.Get("Authorization")

			reResult := tokenRe.FindStringSubmatch(auth)
			var access bool
			fmt.Println(reResult)
			// If valid then access is true
			if len(reResult) > 0 && s.Storage.AuthenticationTokenExists(reResult[1]) {
				access = true
			}

			// Can access?
			if !access {
				logrus.Debugf("Forbidden access with header '%s'", auth)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		f.ServeHTTP(w, r)
	})
}
