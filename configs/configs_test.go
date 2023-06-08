package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:       ":1234",
			BaseURLShortener: "http://example.ru",
			JSONFileDB:       "../../db.json",
		}},
		{name: "Тест 1", want: Config{
			HostServer:       ":7777",
			BaseURLShortener: "",
			JSONFileDB:       "../../dbExample.json",
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
			HostServer:       ":8080",
			BaseURLShortener: "",
			JSONFileDB:       "../../db.json",
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
			JSONFileDB:       "../../db.json",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", tt.want.HostServer)
			os.Setenv("BASE_URL", tt.want.BaseURLShortener)
			SetConfigFromEnv()
			assert.Equal(t, newConfig, tt.want, "newConfig не совпадает с ожидаемым")
		})
	}
}
