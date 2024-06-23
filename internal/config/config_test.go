package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
)

func TestNewConfig(t *testing.T) {
	viper.AddConfigPath("../..")
	t.Setenv("GOTIFY_TOKEN", "jpp")
	got, err := NewConfig()
	if err != nil {
		t.Errorf("Error creating new config: %v", err)
	}

	want := Config{
		GotifyHost:  "http://127.0.0.1:8080",
		GotifyToken: "jpp",
		Services: []Service{
			{
				Host: "http://127.0.0.1:9090",
				Path: "/-/healthy",
				Name: "Prometheus",
			},
			{
				Host: "http://127.0.0.1:9090",
				Path: "/health",
				Name: "Prometheus-fail",
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("The configuration is not matching the one expected.\nGot: %v\nExpected: %v", got, want)
	}
}
