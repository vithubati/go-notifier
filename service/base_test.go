package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vithubati/go-notifier/config"
	"github.com/vithubati/go-notifier/model"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	nService, err := New(&config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			ConnString:       "root:password@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, nService)
}

func TestCreateNotification(t *testing.T) {
	nService, err := New(&config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			ConnString:       "root:password@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, nService)

	n := model.Notification{
		Topic:   "SERVER",
		Action:  "CREATE",
		Subject: "SERVER Created",
		Message: "Server is created for the accountId G445",
	}
	err = nService.CreateNotification(context.Background(), n)
	assert.Nil(t, err)
}
func TestCreateDeliverer(t *testing.T) {
	nService, err := New(&config.Config{
		Notifier: config.Notifier{
			Webhook:          true,
			Slack:            true,
			ConnString:       "root:password@/notifier?parseTime=true",
			DeliveryInterval: 5 * time.Second,
			Migrations:       true,
		},
		Trace:         false,
		JsonLogFormat: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, nService)

	d := model.Deliverer{
		Type:              "WEBHOOK",
		Url:               "https://www.stackoverflow.com/ttest",
		Retry:             3,
		IntervalInSeconds: 10,
		Topics:            []model.DelivererTopic{{Topic: "SERVER"}},
	}
	headers := make(map[string][]string)
	headers["X-Request-id"] = []string{"456456"}
	d.Headers = headers
	err = nService.CreateDeliverer(context.Background(), d)
	assert.Nil(t, err)
}
