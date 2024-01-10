package configs

import (
	"flag"
	"os"
	"time"
)

// AuthConfig Конфигурация авторизации.
type AuthConfig struct {
	SecretKey string
}

// Config Общая конфигурация сервиса.
type Config struct {
	HostServer           string `env:"SERVER_ADDRESS"`
	BaseURLShortener     string `env:"BASE_URL"`
	LogLevel             string
	JSONFileDB           string
	DatabaseDSN          string `env:"DATABASE_DSN"`
	Auth                 AuthConfig
	TimebackupToJSONFile time.Duration
}

var newConfig Config

// GetConfig() выводит не импортируемую переменную newConfig.
func GetConfig() Config {
	return newConfig
}

// InitConfig() Инициализирует локальную не импортируемую переменную newConfig.
// Вызывает все доступные методы получения конфигов.
// Вызывается один раз при старте проекта.
func InitConfig() Config {

	initDefaultConfig()
	parseFlags()
	setConfigFromEnv()

	return newConfig
}

// initDefaultConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения..
// Вызывается один раз при старте проекта.
func initDefaultConfig() {
	newConfig = Config{
		HostServer:           ":8080",
		BaseURLShortener:     "",
		JSONFileDB:           "/tmp/short-url-db.json",
		LogLevel:             "info",
		DatabaseDSN:          "",
		Auth:                 AuthConfig{SecretKey: "supersecretkey"},
		TimebackupToJSONFile: time.Minute * 10,
	}

}

// parseFlags обрабатывает аргументы командной строки.
// и сохраняет их значения в соответствующих переменных.
func parseFlags() {

	flag.StringVar(&newConfig.HostServer, "a", newConfig.HostServer, "адрес и порт для запуска сервера")
	flag.StringVar(&newConfig.DatabaseDSN, "d", newConfig.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&newConfig.JSONFileDB, "f", newConfig.JSONFileDB, "полное имя файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&newConfig.BaseURLShortener, "b", newConfig.BaseURLShortener, "базовый URl  который будет использоваться для короткой ссылки")
	flag.StringVar(&newConfig.LogLevel, "l", "info", "уровень логирования")
	flag.Parse()
}

// setConfigFromEnv() Прсваевает полям значения из ENV.
// Вызывается один раз при старте проекта.
func setConfigFromEnv() Config {

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
