package webhook

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

type Webhook struct {
	// a client to use for sending webhooks
	client  *http.Client
	target  *url.URL
	headers http.Header
}

// New returns a new webhook Deliverer
func New(conf *Config) (*Webhook, error) {
	if conf == nil {
		return nil, errors.New("config not provided")
	}
	if err := conf.validate(); err != nil {
		return nil, err
	}
	var d Webhook
	var err error

	d.target, err = url.Parse(conf.Target)
	if err != nil {
		return nil, err
	}
	d.headers = conf.Headers.Clone()
	if d.headers == nil {
		d.headers = make(map[string][]string)
	}
	d.headers.Set("content-type", "application/json")

	d.client = conf.Client
	return &d, nil
}
func (d *Webhook) Deliver(ctx context.Context, notificationID string) error {
	return nil
}