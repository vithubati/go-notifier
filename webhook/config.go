package webhook

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Config - Webhook Config configures the "webhook" notification properties.
type Config struct {
	// any HTTP headers necessary for the request to Target
	Headers http.Header `json:"headers"`

	// a client to use for sending webhooks
	Client *http.Client

	// the URL where our webhook will be delivered
	Target string `json:"target"`

}

// Validate will validate the configuration properties
func (c *Config) validate() error {
	if _, err := url.Parse(c.Target); err != nil {
		return fmt.Errorf("failed to parse target url: %w", err)
	}
	if c.Target == "" {
		return errors.New("target url is not provided")
	}

	if c.Client == nil {
		return errors.New("client is not provided")
	}

	return nil
}
