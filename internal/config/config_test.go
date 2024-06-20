package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/viper"
)

func TestNewConfig(t *testing.T) {
	viper.AddConfigPath("../..")
	os.Setenv("GOTIFY_TOKEN", "jpp")
	got, err := NewConfig()
	if err != nil {
		t.Errorf("Error creating new config: %v", err)
	}
	expectedJSON := map[string]string{
		"database": "green",
		"health":   "green",
	}
	want := Config{
		GotifyHost:  "http://192.168.0.151:6060",
		GotifyToken: "jpp",
		Services: []Service{
			{
				Host:         "http://192.168.0.151:6060",
				Path:         "/health",
				Name:         "Gotify",
				ExpectedJSON: expectedJSON,
			},
		},
	}
	if !cmp.Equal(got, want) {
		t.Errorf("The configuration is not matching the one expected.\nGot: %v\nExpected: %v", got, want)
	}
}
