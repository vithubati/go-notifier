package mysql

import (
	"context"
	"database/sql"
	"github.com/vithubati/go-notifier/model"
)

// getDeliverers - retrieves the list of Deliverers for the given type
func getDeliverers(ctx context.Context, db *sql.DB, dType string) ([]model.Deliverer, error) {
	const (
		queryDeliverers = "SELECT id, type, url, credentials, createdAt, retry, intervalInSeconds FROM deliverer WHERE type = ? ORDER BY createdAt ASC"
	)

	deliveries := make([]model.Deliverer, 0, 0)
	rows, err := db.QueryContext(ctx, queryDeliverers, dType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d model.Deliverer
		err := rows.Scan(&d.ID, &d.Type, &d.Url, &d.Credentials, &d.CreatedAt, &d.Retry, &d.Interval)
		if err != nil {
			return nil, err
		}
		res, err := getDelivererResources(ctx, db, d.ID)
		if err != nil {
			return nil, err
		}
		d.Resources = res
		deliveries = append(deliveries, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// getDelivererIDsByResource - retrieves the list of Deliverers' IDs by a resource
func getDelivererIDsByResource(ctx context.Context, db *sql.DB, resource string) ([]string, error) {
	const (
		queryDelivererResources = "SELECT DISTINCT deliverer_id FROM deliverer_resource where resource = ?"
	)

	ids := make([]string, 0, 0)
	rows, err := db.QueryContext(ctx, queryDelivererResources, resource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

// getDelivererResources -  retrieves the list of Deliverers' resources for the given deliverer id
func getDelivererResources(ctx context.Context, db *sql.DB, delivererID string) ([]model.DelivererResource, error) {
	const (
		queryDelivererResources = "SELECT deliverer_id, resource FROM deliverer_resource where deliverer_id = ?"
	)

	resources := make([]model.DelivererResource, 0, 0)
	rows, err := db.QueryContext(ctx, queryDelivererResources, delivererID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r model.DelivererResource
		err := rows.Scan(&r.DelivererID, &r.Resource)
		if err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return resources, nil
}

// createdNotifications retrieves the list of notifications that has delivery status = 'CREATED'
func createdNotifications(ctx context.Context, db *sql.DB, delivererId string) ([]model.Notification, error) {
	const (
		queryNotifications = "SELECT n.id, n.resource, n.action, d.status, n.createdAt, n.data, d.id FROM notification n, delivery d WHERE n.id = d.notificationId AND d.status = 'CREATED' AND d.delivererId = ? ORDER BY d.createdAt ASC"
	)

	notifications := make([]model.Notification, 0, 0)
	rows, err := db.QueryContext(ctx, queryNotifications, delivererId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(&n.ID, &n.Resource, &n.Action, &n.Action, &n.CreatedAt, &n.Data, &n.NotificationDeliveryID)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

// failedNotifications retrieves the list of notifications that has delivery status = 'FAILED'
func failedNotifications(ctx context.Context, db *sql.DB, delivererId string) ([]model.Notification, error) {
	const (
		queryDeliverer    = "SELECT retry FROM deliverer where id = ?"
		queryNotification = "SELECT n.id, n.resource, n.action, d.status, n.createdAt, n.data, d.id FROM notification n, delivery d WHERE n.id = d.notificationId AND d.status = 'FAILED' AND d.delivererId = ? AND d.attempt < ? ORDER BY d.createdAt ASC"
	)
	var retry int
	if err := db.QueryRowContext(ctx, queryDeliverer, delivererId).Scan(&retry); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	notifications := make([]model.Notification, 0, 0)
	rows, err := db.QueryContext(ctx, queryNotification, delivererId, retry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(&n.ID, &n.Resource, &n.Action, &n.Action, &n.CreatedAt, &n.Data, &n.NotificationDeliveryID)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}
