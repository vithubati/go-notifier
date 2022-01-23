package notifier

import "context"

// Deliverer provides the method set for delivering notifications
type Deliverer interface {
	// Deliver will push the notification ID to subscribed clients.
	Deliver(ctx context.Context, notificationID string) error
}
