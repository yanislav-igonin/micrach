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
	ThreadBumpLimit      int
	IsCaptchaActive      bool
}

type DbConfig struct {
	Url string
}

func getValueOrDefaultBoolean(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	return value == "true"
}

func getValueOrDefaultInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicln(fmt.Sprintf("Could not parse %s to int", value))
	}
	return intValue
}

func getValueOrDefaultString(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getAppConfig() AppConfig {
	env := getValueOrDefaultString(os.Getenv("ENV"), "release")
	port := getValueOrDefaultInt(os.Getenv("PORT"), 3000)
	seedDb := getValueOrDefaultBoolean(os.Getenv("SEED_DB"), false)
	isRateLimiterEnabled := getValueOrDefaultBoolean(os.Getenv("IS_RATE_LIMITER_ENABLED"), true)
	threadsMaxCount := getValueOrDefaultInt(os.Getenv("THREADS_MAX_COUNT"), 50)
	threadBumpLimit := getValueOrDefaultInt(os.Getenv("THREAD_BUMP_LIMIT"), 500)
	isCaptchaActive := getValueOrDefaultBoolean(os.Getenv("IS_CAPTCHA_ACTIVE"), true)

	return AppConfig{
		Env:                  env,
		Port:                 port,
		SeedDb:               seedDb,
		IsRateLimiterEnabled: isRateLimiterEnabled,
		ThreadsMaxCount:      threadsMaxCount,
		ThreadBumpLimit:      threadBumpLimit,
		IsCaptchaActive:      isCaptchaActive,
	}
}

func getDbConfig() DbConfig {
	url := getValueOrDefaultString(os.Getenv("POSTGRES_URL"), "postgresql://localhost/micrach")

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
