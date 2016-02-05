package service

import (
	"net/http"
)

//Ping only returns that the service is alive
func (s *Service) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
