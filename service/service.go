package service

import (
	"context"
	"github.com/vithubati/go-notifier/model"
)

type Service interface {
	// CreateNotification stores the given notification and creates delivery records
	// based on the subscribed deliverers
	//
	// delivery records are created with CREATED state and will be picked up by the
	// worker to tell the deliverer to kick off the delivery
	CreateNotification(ctx context.Context, n model.Notification) error

	// CreateDeliverer creates a new Deliverer in the DB
	//
	// Deliverer delivers notification deliveries which are in CREATED status
	CreateDeliverer(ctx context.Context, n model.Deliverer) error

	// KickOff initiates the Notifier worker.
	//
	// Worker starts all the enabled deliverer's delivery.
	KickOff(ctx context.Context) error
}
