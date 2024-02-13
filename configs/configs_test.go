package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:       ":8080",
			BaseURLShortener: "",
			JSONFileDB:       "/tmp/short-url-db.json",
			LogLevel:         "debug",
			DatabaseDSN:      "",

			TimebackupToJSONFile: time.Minute * 10,
			EnableHTTPS:          false,
			Auth:                 AuthConfig{SecretKey: "supersecretkey", TrustedSubNet: "0.0.0.0"},
		}},
	}
	os.Args = []string{"test", "-l", "debug", "-t", "0.0.0.0", "-s", "true"}
	os.Setenv("ENABLE_HTTPS", "false")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _ := NewConfig()
			assert.Equal(t, tt.want, cfg, "newConfig не совпадает с ожидаемым")
		})
	}
}
