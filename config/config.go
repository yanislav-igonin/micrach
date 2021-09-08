package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	Env    string
	Port   int
	SeedDb bool
}

type DbConfig struct {
	Url string
}

func getAppConfig() AppConfig {
	env := os.Getenv("ENV")
	if env == "" {
		env = "release"
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "3000"
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(fmt.Sprintf("Could not parse %s to int", portString))
	}

	seedDbString := os.Getenv("SEED_DB")
	seedDb := seedDbString == "true"

	return AppConfig{
		Env:    env,
		Port:   port,
		SeedDb: seedDb,
	}
}

func getDbConfig() DbConfig {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		url = "postgresql://localhost/micrach"
	}
	return DbConfig{
		Url: url,
	}
}

var App AppConfig
var Db DbConfig

func Init() {
	App = getAppConfig()
	Db = getDbConfig()
}
