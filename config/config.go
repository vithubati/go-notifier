package config

import (
	"errors"
	"net/http"
	"time"
)

// Config holds the required configs for the Notifier service
type Config struct {
	Notifier      Notifier
	Trace         bool
	JsonLogFormat bool
}

// Notifier provides Clair Notifier node configuration
type Notifier struct {
	// Configures the notifier for webhook delivery
	Webhook bool
	// A Postgres connection string.
	//
	// Formats:
	// url: "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
	// or
	// string: "user=pqgotest dbname=pqgotest sslmode=verify-full"
	ConnString string
	// A time.ParseDuration parsable string
	//
	// The frequency at which the notifier attempt delivery of created or previously failed
	// notifications
	// If a value smaller then 1 second is provided it will be replaced with the
	// default 5 second delivery interval.
	// IMPORTANT - this value will be overridden if each deliverer's interval value is > 0
	DeliveryInterval time.Duration
	// A "true" or "false" value
	//
	// Whether Notifier nodes handle migrations to their database.
	Migrations bool
	Client     *http.Client
}

func (n *Notifier) Validate() error {
	if n.DeliveryInterval < 1*time.Second {
		n.DeliveryInterval = DefaultNotifierDeliveryInterval
	}
	if len(n.ConnString) == 0 {
		return errors.New("database connection url is required")
	}
	return nil
}
