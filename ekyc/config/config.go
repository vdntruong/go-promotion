package config

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Database `yaml:"database"`
		JWT      `yaml:"jwt"`
		Redis    `yaml:"redis"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}

	HTTP struct {
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read-timeout"`
		WriteTimeout time.Duration `yaml:"write-timeout"`
	}

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}

	JWT struct {
		Issuer string `yaml:"issuer"`
		Secret string `yaml:"secret"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
