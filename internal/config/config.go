package config

import (
	"os"
	"strconv"
)

const (
	DefaultServerPort = 8080
	DefaultDBPort     = 5432
	DefaultDBHost     = "db"
	DefaultDBUser     = "postgres"
	DefaultDBPassword = "password"
	DefaultDBName     = "pr_service"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
	ServerPort int
}

func New() *Config {
	return &Config{
		DBHost:     getEnvString("DB_HOST", DefaultDBHost),
		DBUser:     getEnvString("DB_USER", DefaultDBUser),
		DBPassword: getEnvString("DB_PASSWORD", DefaultDBPassword),
		DBName:     getEnvString("DB_NAME", DefaultDBName),
		DBPort:     getEnvInt("DB_PORT", DefaultDBPort),
		ServerPort: getEnvInt("SERVER_PORT", DefaultServerPort),
	}
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
