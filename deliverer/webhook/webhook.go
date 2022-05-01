package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/vithubati/go-notifier/delivery"
	"github.com/vithubati/go-notifier/model"
	"io"
	"net/http"
	"net/url"
)

// Compile-time check to ensure Webhook implements delivery.Deliverer.
var _ delivery.Deliverer = (*Webhook)(nil)

type Webhook struct {
	// a client to use for sending webhooks
	client  *http.Client
	target  *url.URL
	headers http.Header
}

// New returns a new webhook Deliverer
func New(conf *Config) (delivery.Deliverer, error) {
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

func (d *Webhook) Deliver(ctx context.Context, notification model.Notification) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"notificationId": notification.ID,
		"component":      "webhook/deliverer.Deliver",
	}).WithContext(ctx)
	jsonBytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req := &http.Request{
		URL:    d.target,
		Header: d.headers,
		Body:   io.NopCloser(bytes.NewBuffer(jsonBytes)),
		Method: http.MethodPost,
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("delivery was not successful. " + resp.Status)
	}
	ctxLog.Infof("successfully deliverd to %s. Status code is: %v", req.URL, resp.StatusCode)
	return nil
}
