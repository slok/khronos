// Khornos daemon

package main

import (
	"os"

	"github.com/NYTimes/gizmo/server"

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

	// Create the storage client
	stCli := storage.NewDummy()

	// Create scheduler and start
	cr := schedule.NewDummyCron(cfg, stCli, 0, "OK")
	cr.Start(nil)
	defer cr.Stop()

	// Load service
	khronosService := service.NewKhronosService(cfg, stCli, cr)

	// Register the service on the server
	err := server.Register(khronosService)
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	// Serve our service
	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
