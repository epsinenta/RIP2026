package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		cfg := &Config{
			ServiceHost: "0.0.0.0",
			ServicePort: 8080,
		}
		logrus.Info("config loaded (defaults)")
		return cfg, nil
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	if cfg.ServicePort == 0 {
		cfg.ServicePort = 8080
	}

	logrus.Info("config loaded")
	return cfg, nil
}
