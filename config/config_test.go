package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/NYTimes/gizmo/config"
)

const (
	goodConfig1 = "testdata/goodconf1.json"
	goodConfig2 = "testdata/goodconf2.json"
)

func TestCheckConfigFromFile(t *testing.T) {
	tests := []struct {
		givenConfigFile     string
		wantedConfigFile    string
		wantedLogLevel      string
		wantedHTTPPort      int
		wantedHTTPAccessLog string
	}{
		{
			givenConfigFile:     goodConfig1,
			wantedConfigFile:    goodConfig1,
			wantedHTTPPort:      12345,
			wantedLogLevel:      "debug",
			wantedHTTPAccessLog: "/tmp/access.log",
		},
		{
			givenConfigFile:     goodConfig2,
			wantedConfigFile:    goodConfig2,
			wantedHTTPPort:      98765,
			wantedLogLevel:      "warning",
			wantedHTTPAccessLog: "/var/log/khronos/access.log",
		},
	}

	for _, test := range tests {
		cfg := &AppConfig{
			ConfigFilePath: test.givenConfigFile,
			Server:         &config.Server{},
		}
		cfg.loadConfigFromFile()

		if cfg.ConfigFilePath != test.wantedConfigFile {
			t.Errorf("Config file should be %s; obtained %s", test.wantedConfigFile, cfg.ConfigFilePath)
		}

		if cfg.LogLevel != test.wantedLogLevel {
			t.Errorf("Config log level should be %s; obtained %s", test.wantedLogLevel, cfg.LogLevel)
		}

		if cfg.HTTPPort != test.wantedHTTPPort {
			t.Errorf("Config http port should be %d; obtained %d", test.wantedHTTPPort, cfg.HTTPPort)
		}

		if cfg.HTTPAccessLog != test.wantedHTTPAccessLog {
			t.Errorf("Config http access log should be %s; obtained %s", test.wantedHTTPAccessLog, cfg.HTTPAccessLog)
		}
	}
}

func TestCheckConfigFromEnv(t *testing.T) {
	tests := []struct {
		givenLogLevel  string
		givenHTTPPort  int
		wantedLogLevel string
		wantedHTTPPort int
	}{
		{
			givenLogLevel:  "debug",
			givenHTTPPort:  12345,
			wantedLogLevel: "debug",
			wantedHTTPPort: 12345,
		},
		{
			givenLogLevel:  "warning",
			givenHTTPPort:  98765,
			wantedLogLevel: "warning",
			wantedHTTPPort: 98765,
		},
	}

	for _, test := range tests {
		os.Setenv("APP_LOG_LEVEL", test.givenLogLevel)
		os.Setenv("HTTP_PORT", strconv.Itoa(test.givenHTTPPort))

		cfg := &AppConfig{Server: &config.Server{}}
		cfg.loadConfigFromEnv()

		if cfg.LogLevel != test.wantedLogLevel {
			t.Errorf("Config log level should be %s; obtained %s", test.wantedLogLevel, cfg.LogLevel)
		}

		if cfg.HTTPPort != test.wantedHTTPPort {
			t.Errorf("Config http port should be %d; obtained %d", test.wantedHTTPPort, cfg.HTTPPort)
		}
	}
}

func TestConfigureApp(t *testing.T) {
	tests := []struct {
		givenConfigFile     string
		givenHTTPPort       int
		givenLogLevel       string
		wantedConfigFile    string
		wantedLogLevel      string
		wantedHTTPPort      int
		wantedHTTPAccessLog string
	}{
		{
			givenConfigFile: goodConfig1,
			givenHTTPPort:   4433,
			givenLogLevel:   "error",

			wantedConfigFile:    goodConfig1,
			wantedHTTPPort:      4433,
			wantedLogLevel:      "error",
			wantedHTTPAccessLog: "/tmp/access.log",
		},
		{
			givenConfigFile: goodConfig2,
			givenHTTPPort:   99888,
			givenLogLevel:   "info",

			wantedConfigFile:    goodConfig2,
			wantedHTTPPort:      99888,
			wantedLogLevel:      "info",
			wantedHTTPAccessLog: "/var/log/khronos/access.log",
		},
	}

	for _, test := range tests {
		os.Setenv("APP_LOG_LEVEL", test.givenLogLevel)
		os.Setenv("HTTP_PORT", strconv.Itoa(test.givenHTTPPort))

		cfg := &AppConfig{
			ConfigFilePath: test.givenConfigFile,
			Server:         &config.Server{},
		}
		cfg.ConfigureApp()

		// Check if correct merge of settings and priority is made
		if cfg.ConfigFilePath != test.wantedConfigFile {
			t.Errorf("Config file should be %s; obtained %s", test.wantedConfigFile, cfg.ConfigFilePath)
		}

		if cfg.LogLevel != test.wantedLogLevel {
			t.Errorf("Config log level should be %s; obtained %s", test.wantedLogLevel, cfg.LogLevel)
		}

		if cfg.HTTPPort != test.wantedHTTPPort {
			t.Errorf("Config http port should be %d; obtained %d", test.wantedHTTPPort, cfg.HTTPPort)
		}

		if cfg.HTTPAccessLog != test.wantedHTTPAccessLog {
			t.Errorf("Config http access log should be %s; obtained %s", test.wantedHTTPAccessLog, cfg.HTTPAccessLog)
		}
	}

}
