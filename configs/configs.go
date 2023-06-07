package configs

import "os"

type Config struct {
	HostServer       string `env:"SERVER_ADDRESS"`
	BaseURLShortener string `env:"BASE_URL"`
}

var newConfig Config

func InitConfig() *Config {
	newConfig = Config{
		HostServer:       ":8080",
		BaseURLShortener: "",
	}
	return &newConfig
}

func GetConfig() Config {
	return newConfig
}

func GetConfigFromEnv() {
	if envBaseURLShortener := os.Getenv("SERVER_ADDRESS"); envBaseURLShortener != "" {
		newConfig.HostServer = envBaseURLShortener
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		newConfig.BaseURLShortener = envBaseURL
	}
}
