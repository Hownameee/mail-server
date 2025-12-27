package config

import "os"

type Config struct {
	PORT      string
	EmailUser string
	EmailPass string
	EmailHost string
	EmailPort string
}

func New() *Config {
	return &Config{
		PORT: GetEnv("PORT", "40700"),

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
