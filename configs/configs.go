package configs

type Config struct {
	HostServer       string
	BaseURLShortener string
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
