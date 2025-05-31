package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server
	Redis    Redis
	Database Database
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Database struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	SSLMode  string `mapstructure:"sslmode"`
}

func LoadConfig(env string) (*Config, error) {
	v := viper.New()
	if env == "" {
		env = "local"
	}
	configPath := fmt.Sprintf("config/config.%s.yaml", env)
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
