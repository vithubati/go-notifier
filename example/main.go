package main

import (
	"context"
	"database/sql"
	"github.com/vithubati/go-notifier/config"
	"github.com/vithubati/go-notifier/service"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		db, err := sql.Open("mysql", "root:password@/notifier?parseTime=true")
		if err != nil {
			log.Fatalf("failed to open db error = %v", err)
			return
		}
		defer db.Close()
		cfg := newConfig()
		cfg.Notifier.Client = httpClient()
		s, err := service.New(db, cfg)
		if err != nil {
			log.Fatalf("Notifier() error = %v", err)
			return
		}
		if err := s.KickOff(ctx); err != nil {
			log.Fatalf("Notifier() KickOff error = %v", err)
			return
		}
		return
	}()
	wg.Wait()
}

func newConfig() *config.Config {
	return &config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	}
}

func httpClient() *http.Client {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}
