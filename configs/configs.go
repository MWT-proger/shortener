package configs

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"reflect"
	"strconv"
	"time"
)

// AuthConfig Конфигурация авторизации.
type AuthConfig struct {
	SecretKey     string `default:"supersecretkey"`
	TrustedSubNet string `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
}

// Config Общая конфигурация сервиса.
type Config struct {
	HostServer           string        `json:"server_address" default:":8080" env:"SERVER_ADDRESS"`
	BaseURLShortener     string        `json:"base_url" default:"" env:"BASE_URL"`
	LogLevel             string        `default:"info"`
	JSONFileDB           string        `json:"file_storage_path" default:"/tmp/short-url-db.json"`
	DatabaseDSN          string        `json:"database_dsn" default:"" env:"DATABASE_DSN"`
	TimebackupToJSONFile time.Duration `default:"600000000000"`
	EnableHTTPS          bool          `json:"enable_https" env:"ENABLE_HTTPS"`
	ConfigJSON           string        `env:"CONFIG"`

	Auth AuthConfig
}

// NewConfig Создаёи и возвращает новый объект Config
// Вызывает все доступные методы(в порядке как указано):
// setValuesFromFlags - обрабатывает аргументы командной строки;
// setValuesFromEnv - обрабатывает переменные окружения;
// setValuesFromJSONFile - обрабатывает поля из JSONFile;
// setValuesDefaultIfNil - обрабатывает тег default.
func NewConfig() (Config, error) {
	cfg := Config{}

	cfg.setValuesFromFlags()
	cfg.setValuesFromEnv()

	if err := cfg.setValuesFromJSONFile(); err != nil {
		return cfg, err
	}

	setValuesDefaultIfNil(&cfg)

	return cfg, nil
}

// setValuesFromFlags обрабатывает аргументы командной строки,
// и сохраняет их значения в соответствующих переменных структуры.
func (cfg *Config) setValuesFromFlags() {
	flag.StringVar(&cfg.HostServer, "a", cfg.HostServer, "адрес и порт для запуска сервера")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&cfg.JSONFileDB, "f", cfg.JSONFileDB, "полное имя файла, куда сохраняются данные в формате JSON")
	flag.StringVar(&cfg.BaseURLShortener, "b", cfg.BaseURLShortener, "базовый URl  который будет использоваться для короткой ссылки")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "уровень логирования")
	flag.StringVar(&cfg.Auth.TrustedSubNet, "t", cfg.Auth.TrustedSubNet, "строковое представление бесклассовой адресации (CIDR)")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "включить HTTPS")
	flag.StringVar(&cfg.ConfigJSON, "c", cfg.ConfigJSON, "JSON конфигурации приложения")
	flag.StringVar(&cfg.ConfigJSON, "config", cfg.ConfigJSON, "JSON конфигурации приложения")

	flag.Parse()
}

// setValuesFromEnv() обрабатывает переменные окружения,
// и сохраняет их значения в соответствующих переменных структуры.
func (cfg *Config) setValuesFromEnv() {

	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		cfg.HostServer = envBaseURLShortener
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		cfg.BaseURLShortener = envBaseURL
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		cfg.LogLevel = envLogLevel
	}
	if envTrustedSubNet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubNet != "" {
		cfg.Auth.TrustedSubNet = envTrustedSubNet
	}
	if envJSONFileDB := os.Getenv("FILE_STORAGE_PATH"); envJSONFileDB != "" {
		cfg.JSONFileDB = envJSONFileDB
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		cfg.DatabaseDSN = envDatabaseDSN
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		cfg.Auth.SecretKey = envSecretKey
	}
	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS == "1" || envEnableHTTPS == "true" || envEnableHTTPS == "True" {
		cfg.EnableHTTPS = true
	} else if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS == "0" || envEnableHTTPS == "false" || envEnableHTTPS == "False" {
		cfg.EnableHTTPS = false
	}
	if envConfigJSON := os.Getenv("CONFIG"); envConfigJSON != "" {
		cfg.ConfigJSON = envConfigJSON
	}
}

// setValuesFromJSONFile() обрабатывает поля из JSONFile,
// и сохраняет их значения в соответствующих переменных структуры.
// Возвращает объект Config.
func getValuesFromJSONFile(pathConfigJSON string) (*Config, error) {

	content, err := os.ReadFile(pathConfigJSON)
	configJSON := Config{}

	if err != nil {
		return nil, errors.New(pathConfigJSON + " указанный файл конфигурации не существует.")
	}

	if err = json.Unmarshal(content, &configJSON); err != nil {
		return nil, errors.New(pathConfigJSON + " не верный формат JSON.")
	}

	return &configJSON, nil
}

// setValuesFromJSONFile() обрабатывает поля из JSONFile,
// и сохраняет их значения в соответствующих переменных структуры.
// ТОЛЬКО при условие что у полей структуры нулевые значения.
func (cfg *Config) setValuesFromJSONFile() error {

	if cfg.ConfigJSON != "" {
		configJSON, err := getValuesFromJSONFile(cfg.ConfigJSON)

		if err != nil {
			return err
		}

		valueConfig := reflect.ValueOf(cfg).Elem()
		valueConfigJSON := reflect.ValueOf(configJSON).Elem()

		for i := 0; i < valueConfig.NumField(); i++ {
			field := valueConfig.Field(i)
			fieldJSON := valueConfigJSON.Field(i)

			if field.IsZero() && !fieldJSON.IsZero() {
				field.Set(fieldJSON)
			}
		}
	}
	return nil
}

// setValuesDefaultIfNil() обрабатывает тег default,
// и сохраняет их значения в соответствующих переменных структуры.
// ТОЛЬКО при условие что у полей структуры нулевые значения.
func setValuesDefaultIfNil(cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("default")

		if tag == "" && field.Kind() != reflect.Struct {
			continue
		}
		switch field.Kind() {

		case reflect.Struct:
			setValuesDefaultIfNil(field.Addr().Interface())

		case reflect.String:
			if field.String() == "" {
				field.SetString(tag)
			}
		case reflect.Int64:
			if field.Int() == 0 {
				if intValue, err := strconv.ParseInt(tag, 10, 64); err == nil {
					field.SetInt(intValue)
				}
			}

		}
	}
}
