package mysql

import (
	"context"
	"database/sql"
	"github.com/vithubati/go-notifier/model"
)

// getDeliverers - retrieves the list of Deliverers for the given type
func getDeliverers(ctx context.Context, db *sql.DB, dType string) ([]model.Deliverer, error) {
	const (
		queryDeliverers = "SELECT id, type, url, channelId, headers, credentials, createdAt, retry, intervalInSeconds FROM deliverer WHERE type = ? ORDER BY createdAt ASC"
	)

	deliveries := make([]model.Deliverer, 0, 0)
	rows, err := db.QueryContext(ctx, queryDeliverers, dType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d model.Deliverer
		err := rows.Scan(&d.ID, &d.Type, &d.Url, &d.ChannelID, &d.Headers, &d.Credentials, &d.CreatedAt, &d.Retry, &d.IntervalInSeconds)
		if err != nil {
			return nil, err
		}
		res, err := getDelivererTopics(ctx, db, d.ID)
		if err != nil {
			return nil, err
		}
		d.Topics = res
		deliveries = append(deliveries, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// getDelivererIDsByTopic - retrieves the list of Deliverers' IDs by a topic
func getDelivererIDsByTopic(ctx context.Context, db *sql.DB, topic string) ([]string, error) {
	const (
		queryDelivererTopics = "SELECT DISTINCT deliverer_id FROM deliverer_topic where topic = ?"
	)

	ids := make([]string, 0, 0)
	rows, err := db.QueryContext(ctx, queryDelivererTopics, topic)
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

// getDelivererTopics -  retrieves the list of Deliverers' topics for the given deliverer id
func getDelivererTopics(ctx context.Context, db *sql.DB, delivererID string) ([]model.DelivererTopic, error) {
	const (
		queryDelivererTopics = "SELECT deliverer_id, topic FROM deliverer_topic where deliverer_id = ?"
	)

	topics := make([]model.DelivererTopic, 0, 0)
	rows, err := db.QueryContext(ctx, queryDelivererTopics, delivererID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r model.DelivererTopic
		err := rows.Scan(&r.DelivererID, &r.Topic)
		if err != nil {
			return nil, err
		}
		topics = append(topics, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return topics, nil
}

// createdNotifications retrieves the list of notifications that has delivery status = 'CREATED'
func createdNotifications(ctx context.Context, db *sql.DB, delivererId string) ([]model.Notification, error) {
	const (
		queryNotifications = "SELECT n.id, n.topic, n.action, n.subject, n.message, n.createdAt, n.data, d.id FROM notification n, delivery d WHERE n.id = d.notificationId AND d.status = 'CREATED' AND d.delivererId = ? ORDER BY d.createdAt ASC"
	)

	notifications := make([]model.Notification, 0, 0)
	rows, err := db.QueryContext(ctx, queryNotifications, delivererId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(&n.ID, &n.Topic, &n.Action, &n.Subject, &n.Message, &n.CreatedAt, &n.Data, &n.NotificationDeliveryID)
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
		queryNotification = "SELECT n.id, n.topic, n.action, n.subject, n.message, n.createdAt, n.data, d.id FROM notification n, delivery d WHERE n.id = d.notificationId AND d.status = 'FAILED' AND d.delivererId = ? AND d.attempt < ? ORDER BY d.createdAt ASC"
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
		err := rows.Scan(&n.ID, &n.Topic, &n.Action, &n.Subject, &n.Message, &n.CreatedAt, &n.Data, &n.NotificationDeliveryID)
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
