package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/storage"
)

const prefix = "/api/v1"

// KhronosService is the application served
type KhronosService struct {
	Config *config.AppConfig
	Client storage.Client
}

//NewKhronosService creates a service object ready to be served.
func NewKhronosService(cfg *config.AppConfig, client storage.Client) *KhronosService {
	logrus.Debug("New Khronos service created")
	return &KhronosService{
		Config: cfg,
		Client: client,
	}
}

//Prefix returns the prefix of the service (used by gizmo)
func (s *KhronosService) Prefix() string {
	return prefix
}

// Middleware registers the middlewares to execute in the request flow
func (s *KhronosService) Middleware(h http.Handler) http.Handler {
	return h
}

// JSONMiddleware wraps all the requests around this middlewares
func (s *KhronosService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return j
}

// JSONEndpoints maps the routes to the enpoints
func (s *KhronosService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{

		"/ping": map[string]server.JSONEndpoint{
			// ping is used to check the service is alive
			"GET": s.Ping,
		},

		"/jobs": map[string]server.JSONEndpoint{
			// Returns all the registered jobs
			"GET": s.GetAllJobs,
			// Register a new job
			"POST": s.CreateNewJob,
		},
	}
}
