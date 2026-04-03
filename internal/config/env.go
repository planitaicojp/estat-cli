package config

import "os"

const (
	EnvAppID     = "ESTAT_APP_ID"
	EnvFormat    = "ESTAT_FORMAT"
	EnvLang      = "ESTAT_LANG"
	EnvConfigDir = "ESTAT_CONFIG_DIR"
)

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
