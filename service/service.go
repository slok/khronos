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
	return "/"
}

// Middleware registers the middlewares to execute in the request flow
func (s *KhronosService) Middleware(h http.Handler) http.Handler {
	// Add authentication
	h = s.AuthenticationHandler(h)
	return h
}

// JSONMiddleware wraps all the requests around this middlewares
func (s *KhronosService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return j
}

// JSONEndpoints maps the routes to the API endpoints
func (s *KhronosService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{

		"/api/v1/ping": map[string]server.JSONEndpoint{
			// ping is used to check the service is alive
			"GET": s.Ping,
		},

		"/api/v1/jobs": map[string]server.JSONEndpoint{
			// Returns all the registered jobs
			"GET": s.GetJobs,
			// Register a new job
			"POST": s.CreateNewJob,
		},

		"/api/v1/jobs/{id}": map[string]server.JSONEndpoint{
			"GET":    s.GetJob,
			"DELETE": s.DeleteJob,
		},

		"/api/v1/jobs/{jobID}/results": map[string]server.JSONEndpoint{
			"GET": s.GetResults,
		},

		"/api/v1/jobs/{jobID}/results/{resultID}": map[string]server.JSONEndpoint{
			"GET":    s.GetResult,
			"DELETE": s.DeleteResult,
		},
	}
}

// Endpoints maps the rountes to the regular endpoints
func (s *KhronosService) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/": map[string]http.HandlerFunc{
			"GET": s.JobsList,
		},
	}
}
