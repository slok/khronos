package service

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

//Ping informs service is alive
func (s *KhronosService) ping(r *http.Request) (int, interface{}, error) {
	logrus.Debug("Calling ping endpoint")
	return http.StatusOK, "pong", nil
}
