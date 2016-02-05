package main

import (
	"os"

	"github.com/NYTimes/gizmo/server"

	"github.com/slok/khronos/config"
	"github.com/slok/khronos/service"
)

const khronosConfigFile = "KHRONOS_CONFIG_FILE"

func main() {
	// Get config location file
	configFile := os.Getenv(khronosConfigFile)

	// Load config
	cfg := config.NewAppConfig(configFile)
	server.Init("khronos", cfg.Server)

	// Load service
	khronosService := service.NewService(cfg)

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
