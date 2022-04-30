# Notifier

Notifier is a GO library for sending notifications to external messaging applications.

### To do

* [ ]  Webhook credentials support
* [ ]  Postgre support

### Usage

```go
package main

import (
	"context"
	"github.com/vithubati/go-notifier/config"
	"github.com/vithubati/go-notifier/service"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func newConfig() *config.Config {
	return &config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			ConnString:       "root:password@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s, err := Notifier(newConfig())
		if err != nil {
			log.Fatalf("Notifier() error = %v", err)
			return
		}
		if err := s.KickOff(ctx); err != nil {
			log.Fatalf("Notifier() error = %v", err)
			return
		}
		return
	}()
	wg.Wait()
}

func Notifier(cfg *config.Config) (service.Service, error) {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	c := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	cfg.Notifier.Client = c
	s, err := service.New(cfg)
	if err != nil {
		return nil, err
	}
	return s, nil
}

````
