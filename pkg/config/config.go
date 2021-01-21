package config

type AppConfig struct {
	AppPort string
}

func InitConf() *AppConfig {
	return &AppConfig{
		AppPort: "9249",
	}
}
