package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port               string
	Workers            int
	OtpCleanupSeconds  int
	OtpLifeSpanSeconds int
	EmailUser          string
	EmailPass          string
	EmailHost          string
	EmailPort          string
}

func New() *Config {
	return &Config{
		Port:               GetEnv("PORT", "40700"),
		Workers:            GetEnvInt("WORKERS", 5),
		OtpCleanupSeconds:  GetEnvInt("OTP_CLEANUP_SECONDS", 180),
		OtpLifeSpanSeconds: GetEnvInt("OTP_LIFESPAN_SECONDS", 120),

		EmailUser: GetEnv("EMAIL", ""),
		EmailPass: GetEnv("APP_PASSWORD", ""),
		EmailHost: GetEnv("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort: GetEnv("EMAIL_PORT", "465"),
	}
}

func GetEnv(env string, defaultEnv string) string {
	res := os.Getenv(env)
	if res == "" {
		res = defaultEnv
	}
	return res
}

func GetEnvInt(env string, defaultEnv int) int {
	res := os.Getenv(env)
	if res == "" {
		return defaultEnv
	}
	val, err := strconv.Atoi(res)
	if err != nil {
		return defaultEnv
	}
	return val
}
