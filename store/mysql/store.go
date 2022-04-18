package mysql

import (
	"context"
	"database/sql"
	"github.com/vithubati/go-notifier/model"
	s "github.com/vithubati/go-notifier/store"
)

// Store implements the store.Store interface
type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) s.Store {
	return &store{db}
}

func (s store) Create(ctx context.Context, n model.Notification) error {
	return createNotification(ctx, s.db, n)
}

func (s store) GetCreated(ctx context.Context, delivererId string) ([]model.Notification, error) {
	return createdNotifications(ctx, s.db, delivererId)
}

func (s store) Delivered(ctx context.Context, id string) error {
	return delivered(ctx, s.db, id)
}

func (s store) Failed(ctx context.Context, id string) error {
	return failed(ctx, s.db, id)
}

func (s store) GetFailed(ctx context.Context, delivererId string) ([]model.Notification, error) {
	return failedNotifications(ctx, s.db, delivererId)
}

func (s store) GetDeliverer(ctx context.Context, dType string) ([]model.Deliverer, error) {
	return getDeliverers(ctx, s.db, dType)
}

func (s store) CreateDeliverer(ctx context.Context, n model.Deliverer) error {
	return createDeliverer(ctx, s.db, n)
}
