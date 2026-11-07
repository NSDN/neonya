// Package config 从环境变量装载进程配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Environment 是运行环境。
type Environment string

const (
	// EnvironmentDevelopment 是本地开发环境，缺失项回落到开发默认值。
	EnvironmentDevelopment Environment = "dev"
	// EnvironmentProduction 是线上环境。
	EnvironmentProduction Environment = "prod"
)

// Config 是单个进程启动所需的全部配置。
type Config struct {
	ApplicationEnvironment Environment
	HTTPAddress            string
	DatabaseURL            string
	SessionCookieName      string
	StorageDirectory       string
	TrustedProxyCount      int
}

// Load 读取环境变量并校验。DATABASE_URL 没有默认值：连哪个库必须是显式决定。
func Load() (Config, error) {
	trustedProxyCount, err := environmentInteger("TRUSTED_PROXY_COUNT", 0)
	if err != nil {
		return Config{}, err
	}

	configuration := Config{
		ApplicationEnvironment: Environment(environmentVariable("APP_ENV", string(EnvironmentDevelopment))),
		HTTPAddress:            environmentVariable("HTTP_ADDR", ":8080"),
		DatabaseURL:            environmentVariable("DATABASE_URL", ""),
		SessionCookieName:      environmentVariable("SESSION_COOKIE_NAME", "neonya_session"),
		StorageDirectory:       environmentVariable("STORAGE_DIR", "./var/uploads"),
		TrustedProxyCount:      trustedProxyCount,
	}

	switch configuration.ApplicationEnvironment {
	case EnvironmentDevelopment, EnvironmentProduction:
	default:
		return Config{}, fmt.Errorf("unsupported APP_ENV %q", configuration.ApplicationEnvironment)
	}

	if configuration.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return configuration, nil
}

func environmentVariable(name string, fallback string) string {
	if value, found := os.LookupEnv(name); found && value != "" {
		return value
	}
	return fallback
}

func environmentInteger(name string, fallback int) (int, error) {
	value, found := os.LookupEnv(name)
	if !found || value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, err)
	}
	if parsed < 0 {
		return 0, fmt.Errorf("%s must not be negative", name)
	}

	return parsed, nil
}
