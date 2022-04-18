package config

import "time"

const(
	// DefaultNotifierDeliveryInterval is the default (and minimum) interval for
	// the notifier's delivery interval. The notifier will attempt to deliver
	// outstanding notifications at this rate.
	DefaultNotifierDeliveryInterval = 5 * time.Second
)
