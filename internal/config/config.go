package config

import (
	"errors"
	"os"

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

	err := viper.ReadInConfig()
	if err != nil {
		return c, err
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return c, err
	}

	if c.GotifyToken == "" {
		if gotifyToken := os.Getenv("GOTIFY_TOKEN"); gotifyToken != "" {
			c.GotifyToken = gotifyToken
		} else {
			return c, errors.New("Error retrieving the gotify token from the config file or environment variables.")
		}
	}

	return c, nil
}
