// Khornos daemon

package main

import (
	"os"
	"time"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/schedule"
	"github.com/slok/khronos/service"
	"github.com/slok/khronos/storage"
)

func main() {
	// Get config location file
	configFile := os.Getenv(config.KhronosConfigFileKey)

	// Load config
	cfg := config.NewAppConfig(configFile)
	server.Init("khronos", cfg.Server)

	var stCli storage.Client
	var err error

	// Create the storage client
	switch cfg.StorageEngine {
	case "dummy":
		stCli = storage.NewDummy()
	case "boltdb":
		to := time.Duration(cfg.BoltDBTimeoutSeconds)
		stCli, err = storage.NewBoltDB(cfg.BoltDBPath, to)
		if err != nil {
			logrus.Fatalf("Error opening boltdb database: %v", err)
		}
	default:
		logrus.Fatal("Wrong Storage engine")
	}

	// Create scheduler and start
	cr := schedule.NewDummyCron(cfg, stCli, 0, "OK")
	cr.Start(nil)
	defer cr.Stop()

	// Load service
	khronosService := service.NewKhronosService(cfg, stCli, cr)

	// Register the service on the server
	err = server.Register(khronosService)
	if err != nil {
		logrus.Fatalf("unable to register service: %v", err)
	}

	// Serve our service
	err = server.Run()
	if err != nil {
		logrus.Fatalf("server encountered a fatal error: %v", err)
	}
}
