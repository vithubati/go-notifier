package delivery

import (
	"context"
	"github.com/vithubati/go-notifier/model"
)

// Deliverer provides the method set for delivering notifications
type Deliverer interface {
	// Deliver will push the notification to subscribed clients via
	// configured deliverers
	Deliver(ctx context.Context, notification model.Notification) error
}
