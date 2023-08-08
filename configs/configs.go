package configs

import (
	"os"
)

type AuthConfig struct {
	SecretKey string
}

type Config struct {
	HostServer       string `env:"SERVER_ADDRESS"`
	BaseURLShortener string `env:"BASE_URL"`
	LogLevel         string
	JSONFileDB       string
	DatabaseDSN      string `env:"DATABASE_DSN"`
	Auth             AuthConfig
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		HostServer:       ":8080",
		BaseURLShortener: "",
		JSONFileDB:       "/tmp/short-url-db.json",
		LogLevel:         "info",
		DatabaseDSN:      "",
		Auth:             AuthConfig{SecretKey: "supersecretkey"},
	}
	return &newConfig
}

// GetConfig() выводит не импортируемую переменную newConfig
func GetConfig() Config {
	return newConfig
}

// SetConfigFromEnv() Прсваевает полям значения из ENV
// Вызывается один раз при старте проекта
func SetConfigFromEnv() Config {

	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		newConfig.HostServer = envBaseURLShortener
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		newConfig.BaseURLShortener = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		newConfig.LogLevel = envLogLevel
	}
	if envJSONFileDB := os.Getenv("FILE_STORAGE_PATH"); envJSONFileDB != "" {
		newConfig.JSONFileDB = envJSONFileDB
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		newConfig.DatabaseDSN = envDatabaseDSN
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		newConfig.Auth.SecretKey = envSecretKey
	}
	return newConfig
}
