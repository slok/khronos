package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/schedule"
	"github.com/slok/khronos/storage"
)

const prefix = "/api/v1"

// KhronosService is the application served
type KhronosService struct {
	Config  *config.AppConfig
	Storage storage.Client
	Cron    *schedule.Cron
}

//NewKhronosService creates a service object ready to be served.
func NewKhronosService(cfg *config.AppConfig, storage storage.Client, cron *schedule.Cron) *KhronosService {
	logrus.Debug("New Khronos service created")
	return &KhronosService{
		Config:  cfg,
		Storage: storage,
		Cron:    cron,
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
			"GET": s.GetJobs,
			// Register a new job
			"POST": s.CreateNewJob,
		},

		"/jobs/{id}": map[string]server.JSONEndpoint{
			"GET":    s.GetJob,
			"DELETE": s.DeleteJob,
		},

		"/jobs/{jobID}/results": map[string]server.JSONEndpoint{
			"GET": s.GetResults,
		},

		"/jobs/{jobID}/results/{resultID}": map[string]server.JSONEndpoint{
			"GET":    s.GetResult,
			"DELETE": s.DeleteResult,
		},
	}
}
