package service

import (
	"net/http"

	"github.com/slok/khronos/config"
)

const prefix = "/api/v1"

// Service is the application served
type Service struct {
	config *config.AppConfig
}

//NewService creates a service object ready to be served.
func NewService(cfg *config.AppConfig) *Service {
	return &Service{config: cfg}
}

//Prefix returns the prefix of the service (used by gizmo)
func (s *Service) Prefix() string {
	return prefix
}

// Middleware registers the middlewares to execute in the request flow
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// Endpoints maps the routes to the enpoints
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{

		// ping is used to check the service is alive
		"/ping": map[string]http.HandlerFunc{
			"GET": s.ping,
		},
	}
}
