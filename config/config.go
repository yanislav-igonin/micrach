package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type AppConfig struct {
	Env                  string
	Port                 int
	SeedDb               bool
	IsRateLimiterEnabled bool
	ThreadsMaxCount      int
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
		log.Panicln(fmt.Sprintf("Could not parse %s to int", portString))
	}

	seedDbString := os.Getenv("SEED_DB")
	seedDb := seedDbString == "true"

	isRateLimiterEnabledString := os.Getenv("IS_RATE_LIMITER_ENABLED")
	isRateLimiterEnabled := isRateLimiterEnabledString == "true"

	threadsMaxCountString := os.Getenv("THREADS_MAX_COUNT")
	threadsMaxCount, err := strconv.Atoi(threadsMaxCountString)
	if err != nil {
		log.Panicln(fmt.Sprintf("Could not parse %s to int", threadsMaxCountString))
	}

	return AppConfig{
		Env:                  env,
		Port:                 port,
		SeedDb:               seedDb,
		IsRateLimiterEnabled: isRateLimiterEnabled,
		ThreadsMaxCount:      threadsMaxCount,
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
