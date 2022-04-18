package config

import (
	"time"
)

// Config holds the required configs for the Notifier service
type Config struct {
	Notifier Notifier
	Trace    bool
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
	DeliveryInterval time.Duration
	// A "true" or "false" value
	//
	// Whether Notifier nodes handle migrations to their database.
	Migrations bool
}

func (n *Notifier) validate() error {
	if n.DeliveryInterval < 1*time.Second {
		n.DeliveryInterval = DefaultNotifierDeliveryInterval
	}
	return nil
}
