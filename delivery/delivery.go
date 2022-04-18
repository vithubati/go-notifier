package delivery

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vithubati/go-notifier/model"
	"github.com/vithubati/go-notifier/store"
	"time"
)

// Delivery handles the business logic of delivering notifications.
type Delivery struct {
	// ID of the deliverer who is delivering the deliveries
	DelivererID string
	// a Deliverer implementation to invoke.
	Deliverer Deliverer
	// the interval at which we will attempt delivery of notifications.
	interval time.Duration
	// a store to retrieve notifications and update their receipts
	store store.Store
}

func NewDelivery(delivererID string, d Deliverer, interval time.Duration, store store.Store) *Delivery {
	return &Delivery{
		DelivererID: delivererID,
		Deliverer:   d,
		interval:    interval,
		store:       store,
	}
}

// Deliver begins delivering notifications.
//
// Canceling the ctx will end delivery.
func (d *Delivery) Deliver(ctx context.Context) {
	ctxLog := logrus.WithFields(logrus.Fields{
		"delivererId": d.DelivererID,
		"component":   "delivery/delivery.Deliver",
	}).WithContext(ctx)

	ctxLog.Info("delivering notifications")
	go d.deliver(ctx)
}

// deliver is intended to be ran as a go routine.
//
// implements a blocking event loop via a time.Ticker
func (d *Delivery) deliver(ctx context.Context) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"delivererId": d.DelivererID,
		"component":   "delivery/delivery.deliver",
	}).WithContext(ctx)
	ctxLog.Infof("d.interval: %v", d.interval)
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ctxLog.Infof("delivery tick")
			err := d.runDelivery(ctx)
			if err != nil {
				ctxLog.WithContext(ctx).
					WithError(err).Error("encountered error on tick")
			}
		}
	}
}

// runDelivery determines notifications to deliver and
// calls the implemented Deliverer to perform the actions.
func (d *Delivery) runDelivery(ctx context.Context) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"delivererId": d.DelivererID,
		"component":   "delivery/delivery.RunDelivery",
	}).WithContext(ctx)

	var toDeliver []model.Notification
	// get created
	created, err := d.store.GetCreated(ctx, d.DelivererID)
	if err != nil {
		return err
	}
	if ln := len(created); ln != 0 {
		ctxLog.WithFields(logrus.Fields{"created:": ln}).
			Info("Found notifications in created status")
		toDeliver = append(toDeliver, created...)
	}

	// get failed
	failed, err := d.store.GetFailed(ctx, d.DelivererID)
	if err != nil {
		return err
	}
	if ln := len(failed); ln != 0 {
		ctxLog.WithFields(logrus.Fields{"failed:": ln}).
			Info("Found notifications in failed status")
		toDeliver = append(toDeliver, failed...)
	}

	for _, n := range toDeliver {
		err = d.do(ctx, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// do performs the delivery of notifications via the composed deliverer
//
// TODO do's actions should be performed under a distributed lock.
func (d *Delivery) do(ctx context.Context, n model.Notification) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"notificationId": n.ID,
		"component":      "delivery/delivery.do",
	}).WithContext(ctx)

	// deliver the notification
	if err := d.Deliverer.Deliver(ctx, n); err != nil {
		// OK for this to fail, notification will stay in Created status.
		// notifier is failing, lets back off it until next tick.
		ctxLog.Error("failed to deliver notifications", err)
		if err := d.store.Failed(ctx, n.NotificationDeliveryID); err != nil {
			return err
		}
		return nil
	}
	// mark as delivered
	if err := d.store.Delivered(ctx, n.NotificationDeliveryID); err != nil {
		return err
	}

	ctxLog.Infof("successfully delivered notification %s", n.ID)
	return nil
}
