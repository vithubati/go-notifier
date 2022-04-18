package store

import (
	"context"
	"github.com/vithubati/go-notifier/model"
)

// Store is an interface contains store persistence layer signature
type Store interface {
	NotificationStore
	RecipientStore
}

type NotificationStore interface {
	// Create persists the provided notification - used by client
	Create(ctx context.Context, n model.Notification) error
}

type RecipientStore interface {
	// GetCreated retrieves the list of notifications that are not delivered yet
	GetCreated(ctx context.Context, delivererId string) ([]model.Notification, error)

	// GetFailed retrieves the list of notifications that are failed to deliver
	GetFailed(ctx context.Context, delivererId string) ([]model.Notification, error)

	// Delivered marks the provided notificationDelivery id as DELIVERED
	Delivered(ctx context.Context, id string) error

	// Failed marks the provided notificationDelivery id as FAILED
	Failed(ctx context.Context, id string) error

	// CreateDeliverer persists the provided Deliverer - used by client
	// Eg:- webhooks and other Deliverer
	CreateDeliverer(ctx context.Context, n model.Deliverer) error

	// GetDeliverer retrieves the list of Deliverer for the given Deliverer Type
	GetDeliverer(ctx context.Context, dType string) ([]model.Deliverer, error)
}
