package config

import "testing"

func TestGetValueOrDefaultBoolean(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		want         bool
	}{
		{name: "default true", defaultValue: true, want: true},
		{name: "default false", defaultValue: false, want: false},
		{name: "explicit true", value: "true", want: true},
		{name: "explicit false", value: "false", defaultValue: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueOrDefaultBoolean(tt.value, tt.defaultValue); got != tt.want {
				t.Fatalf("getValueOrDefaultBoolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValueOrDefaultInt(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue int
		want         int
	}{
		{name: "default", defaultValue: 3000, want: 3000},
		{name: "configured", value: "8080", defaultValue: 3000, want: 8080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValueOrDefaultInt(tt.value, tt.defaultValue); got != tt.want {
				t.Fatalf("getValueOrDefaultInt() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDefaultConfigs(t *testing.T) {
	for _, key := range []string{
		"ENV", "PORT", "IS_DB_SEEDED", "IS_RATE_LIMITER_ENABLED",
		"THREADS_MAX_COUNT", "THREAD_BUMP_LIMIT", "IS_CAPTCHA_ACTIVE",
		"GATEWAY_URL", "GATEWAY_API_KEY", "GATEWAY_BOARD_ID",
		"GATEWAY_BOARD_URL", "GATEWAY_BOARD_DESCRIPTION", "POSTGRES_URL",
	} {
		t.Setenv(key, "")
	}

	app := getAppConfig()
	if app.Env != "release" || app.Port != 3000 || app.IsDbSeeded {
		t.Fatalf("unexpected app defaults: %+v", app)
	}
	if !app.IsRateLimiterEnabled || app.ThreadsMaxCount != 50 || app.ThreadBumpLimit != 500 || !app.IsCaptchaActive {
		t.Fatalf("unexpected feature defaults: %+v", app)
	}
	if db := getDbConfig(); db.Url != "postgresql://localhost/micrach" {
		t.Fatalf("unexpected database default: %+v", db)
	}
}
