package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vithubati/go-notifier/model"
	"time"
)

// createNotification inserts the provided notification
// and creates records for notification delivery based on subscribed deliverers
//
// these operations occur under a transaction to preserve an atomic operation.
func createNotification(ctx context.Context, db *sql.DB, n model.Notification) error {
	const (
		insertNotification = `INSERT INTO notification (id, action, resource, subject, message, createdAt, data) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?);`
		insertDelivery     = `INSERT INTO delivery (id, notificationId, delivererId, status, createdAt, updatedAt) VALUES (?, ?, ?, 'CREATED', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`
	)
	deliverers, err := getDelivererIDsByResource(ctx, db, n.Resource)
	if err != nil {
		return err
	}
	n.ID = uuid.NewString()
	n.CreatedAt = time.Now()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// insert into Notification table
	rslt, err := tx.ExecContext(ctx, insertNotification, n.ID, n.Action, n.Resource, n.Subject, n.Message, n.Data)
	if err != nil {
		return err
	}
	r, err := rslt.RowsAffected()
	if err != nil {
		return errors.New(fmt.Sprintf("when inserting notification, %s", err.Error()))
	}
	if r <= 0 {
		return fmt.Errorf("no rows affected when inserting notification")
	}

	// insert into delivery
	for _, delivererId := range deliverers {
		rslt, err = tx.ExecContext(ctx, insertDelivery, uuid.New().String(), n.ID, delivererId)
		if err != nil {
			return err
		}
		r, err = rslt.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("when inserting delivery, %s", err.Error()))
		}
		if r <= 0 {
			return fmt.Errorf("no rows affected when inserting delivery")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("when commiting notification insert, %s", err.Error()))
	}
	return nil
}

// createDeliverer inserts the provided Deliverer eg:= webhook
//
// these operations occur under a transaction to preserve an atomic operation.
func createDeliverer(ctx context.Context, db *sql.DB, d model.Deliverer) error {
	const (
		insertDeliverer = `INSERT INTO deliverer (id, type, url, headers, credentials, createdAt, retry, intervalInSeconds) 
VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?, ?);`
		insertDelivererResource = "INSERT INTO deliverer_resource (deliverer_id, resource) VALUES (?, ?); "
	)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// insert into Notification table
	d.ID = uuid.NewString()
	rslt, err := tx.ExecContext(ctx, insertDeliverer, uuid.New().String(), d.Type, d.Url, d.Headers, d.Credentials, d.Retry, d.IntervalInSeconds)
	if err != nil {
		return err
	}
	r, err := rslt.RowsAffected()
	if err != nil {
		return errors.New(fmt.Sprintf("when inserting Deliverer, %s", err.Error()))
	}
	if r <= 0 {
		return fmt.Errorf("no rows affected when inserting Deliverer")
	}

	// insert into deliverer resource
	for _, res := range d.Resources {
		rslt, err = tx.ExecContext(ctx, insertDelivererResource, d.ID, res.Resource)
		if err != nil {
			return err
		}
		r, err = rslt.RowsAffected()
		if err != nil {
			return errors.New(fmt.Sprintf("when inserting Deliverer resource, %s", err.Error()))
		}
		if r <= 0 {
			return fmt.Errorf("no rows affected when inserting Deliverer resource")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New(fmt.Sprintf("when commiting Deliverer insert, %s", err.Error()))
	}
	return nil
}
