package config

import (
	"github.com/Mexx77/ridesharing/logging"
	"os"
)

var config *Config

// env var constants
const environment = "ENVIRONMENT"
const DevEnvironment = "dev"

type Config struct {
	Environment string
}

func GetConfig() *Config {
	if config == nil{
		return &Config{
			Environment: getConfigString(environment),
		}
	} else {
		return config
	}
}

func getConfigString(envVar string) string {
	entry := os.Getenv(envVar)
	if entry == "" {
		logging.Error.Print("Please configure " + envVar)
		os.Exit(1)
	}
	logging.Info.Print(envVar+":", entry)
	return entry
}
