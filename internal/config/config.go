package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Services    []Service
	GotifyToken string
	GotifyHost  string
}

type Service struct {
	Host         string
	Path         string
	Name         string
	ExpectedJSON map[string]string
}

func NewConfig() (Config, error) {
	var c Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/health-poller")
	viper.AddConfigPath("$HOME/.health-poller")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return c, err
		} else {
			return c, err
		}
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}
