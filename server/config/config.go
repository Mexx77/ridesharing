package config

import (
	"github.com/Mexx77/ridesharing/logging"
	"os"
)

var config *Config

// env vars
const environment = "ENVIRONMENT"
const port 		  = "PORT"
const mongoPw 	  = "MONGO_PW"
const jwtSecret   = "JWT_SECRET"

// build constants
const DevEnvironment = "dev"

type Config struct {
	Environment string
	Port        string
	MongoPw     string
	JwtSecret   []byte
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Environment: getConfigString(environment, true),
			Port:        getConfigString(port, true),
			MongoPw:     getConfigString(mongoPw, false),
			JwtSecret:   []byte(getConfigString(jwtSecret, false)),
		}
		return config
	} else {
		return config
	}
}

func getConfigString(envVar string, print bool) string {
	entry := os.Getenv(envVar)
	if entry == "" {
		logging.Error.Print("Please configure " + envVar)
		os.Exit(1)
	}
	if print {
		logging.Info.Print(envVar+": ", entry)
	} else {
		logging.Info.Print(envVar + ": ******")
	}
	return entry
}
