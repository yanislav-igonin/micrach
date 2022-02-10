package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type GatewayConfig struct {
	Url              string
	ApiKey           string
	BoardId          string
	BoardDescription string
}

type AppConfig struct {
	Env                  string
	Port                 int
	IsDbSeeded           bool
	IsRateLimiterEnabled bool
	ThreadsMaxCount      int
	ThreadBumpLimit      int
	IsCaptchaActive      bool
	Gateway              GatewayConfig
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

func getGatewayConfig() GatewayConfig {
	url := os.Getenv("GATEWAY_URL")
	apiKey := os.Getenv("GATEWAY_API_KEY")
	boardId := os.Getenv("GATEWAY_BOARD_ID")
	description := os.Getenv("GATEWAY_BOARD_DESCRIPTION")

	return GatewayConfig{
		Url:              url,
		ApiKey:           apiKey,
		BoardId:          boardId,
		BoardDescription: description,
	}
}

func getAppConfig() AppConfig {
	env := getValueOrDefaultString(os.Getenv("ENV"), "release")
	port := getValueOrDefaultInt(os.Getenv("PORT"), 3000)
	isDbSeeded := getValueOrDefaultBoolean(os.Getenv("IS_DB_SEEDED"), false)
	isRateLimiterEnabled := getValueOrDefaultBoolean(os.Getenv("IS_RATE_LIMITER_ENABLED"), true)
	threadsMaxCount := getValueOrDefaultInt(os.Getenv("THREADS_MAX_COUNT"), 50)
	threadBumpLimit := getValueOrDefaultInt(os.Getenv("THREAD_BUMP_LIMIT"), 500)
	isCaptchaActive := getValueOrDefaultBoolean(os.Getenv("IS_CAPTCHA_ACTIVE"), true)
	gateway := getGatewayConfig()

	return AppConfig{
		Env:                  env,
		Port:                 port,
		IsDbSeeded:           isDbSeeded,
		IsRateLimiterEnabled: isRateLimiterEnabled,
		ThreadsMaxCount:      threadsMaxCount,
		ThreadBumpLimit:      threadBumpLimit,
		IsCaptchaActive:      isCaptchaActive,
		Gateway:              gateway,
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
