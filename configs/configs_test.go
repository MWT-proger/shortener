package configs

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:           ":1234",
			BaseURLShortener:     "http://example.ru",
			JSONFileDB:           "../../db.json",
			LogLevel:             "info",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			TimebackupToJSONFile: time.Minute * 10,
		}},
		{name: "Тест 2", want: Config{
			HostServer:           ":7777",
			BaseURLShortener:     "",
			JSONFileDB:           "../../dbExample.json",
			LogLevel:             "debug",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			TimebackupToJSONFile: time.Minute * 10,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newConfig = tt.want

			got := GetConfig()

			assert.Equal(t, got, tt.want, "GetConfig() не совпадает с ожидаемым")
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:           ":8080",
			BaseURLShortener:     "",
			JSONFileDB:           "/tmp/short-url-db.json",
			LogLevel:             "info",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			TimebackupToJSONFile: time.Minute * 10,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitConfig()
			assert.Equal(t, newConfig, tt.want, "newConfig не совпадает с ожидаемым")
		})
	}
}

func TestGetConfigFromEnv(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:       ":7777",
			BaseURLShortener: "http://site.ru",
			JSONFileDB:       "/tmp/db.json",
			LogLevel:         "info",
			DatabaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				`localhost`, `postgres`, `postgres`, `testDB`),
			Auth:                 AuthConfig{SecretKey: "NewSuperSecretKeyTEEEEEEEEEEST"},
			TimebackupToJSONFile: 0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", tt.want.HostServer)
			os.Setenv("BASE_URL", tt.want.BaseURLShortener)
			os.Setenv("LOG_LEVEL", tt.want.LogLevel)
			os.Setenv("FILE_STORAGE_PATH", tt.want.JSONFileDB)
			os.Setenv("DATABASE_DSN", tt.want.DatabaseDSN)
			os.Setenv("SECRET_KEY", tt.want.Auth.SecretKey)
			setConfigFromEnv()
			assert.Equal(t, newConfig, tt.want, "newConfig не совпадает с ожидаемым")
		})
	}
}
