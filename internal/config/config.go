package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Env           string
	BotAuthToken  string
	AdminUserId   int
	FinManagerAPI FinManagerAPIConfig
}

type FinManagerAPIConfig struct {
	BaseUrl   string
	AuthToken string
}

func MustLoadConfig() *Config {
	admin_user_id := mustGetEnv("ADMIN_USER_ID")
	converted_admin_user_id, err := strconv.Atoi(admin_user_id)
	if err != nil {
		log.Fatalf("ADMIN_USER_ID must be an integer")
	}

	config := &Config{
		Env:          getEnv("ENV", "local"),
		BotAuthToken: mustGetEnv("BOT_AUTH_TOKEN"),
		AdminUserId:  converted_admin_user_id,
		FinManagerAPI: FinManagerAPIConfig{
			BaseUrl:   mustGetEnv("FIN_MANAGER_API_BASE_URL"),
			AuthToken: mustGetEnv("FIN_MANAGER_API_AUTH_TOKEN"),
		},
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return value
}
