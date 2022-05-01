# Notifier

Notifier is a GO library for sending notifications to external messaging applications on demand. For example, when an
event occurs on a topic in an application, if you want that event to be sent as a notification to external services
(such as slack or webhooks) which have subscribed to, this library will take care of it.

### Model

<img alt="db_model.png" height="350" src="./assets/db_model.png" width="400"/>

### Usage

```shell
go get -u github.com/vithubati/go-notifier
```

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

func newConfig() *config.Config {
	return &config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			ConnString:       "<usernam>>:<password>@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	}
}

````

### To do

* [ ]  Webhook credentials support
* [ ]  Postgre support
