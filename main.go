package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/sonlis/health-poller/internal/alerting"
	"github.com/sonlis/health-poller/internal/config"
	"github.com/sonlis/health-poller/internal/healthcheck"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}
	var wg sync.WaitGroup
	for index := range conf.Services {
		wg.Add(1)
		go func(service config.Service) {
			defer wg.Done()
			unealthyMessage, err := healthcheck.GetHealth(service.Host, service.Path, service.ExpectedJSON)
			if err != nil {
				log.Printf("Checking the health of %s failed: %v", service.Name, err)
				err := alerting.AlertServiceUnealthy(service.Name, fmt.Sprintf("impossible to reach %s: %v", service.Name, err), conf.GotifyToken, conf.GotifyHost)
				if err != nil {
					log.Printf("Sending the alert to gotify failed: %v", err)
				}
			}
			if unealthyMessage != "" {
				err := alerting.AlertServiceUnealthy(service.Name, fmt.Sprintf("%s is unealthy: %s", service.Name, unealthyMessage), conf.GotifyToken, conf.GotifyHost)
				if err != nil {
					log.Printf("Sending the alert to gotify failed: %v", err)
				}
			}
		}(conf.Services[index])
	}
	wg.Wait()
}
