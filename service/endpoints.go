package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

//Ping only returns that the service is alive
func (s *Service) ping(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("Calling ping endpoint")
	w.Write([]byte("pong"))
}
