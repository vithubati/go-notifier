package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/remind101/migrate"
	"github.com/sirupsen/logrus"
	"github.com/vithubati/go-notifier/delivery"
	"github.com/vithubati/go-notifier/migrations"
	"github.com/vithubati/go-notifier/model"
	"github.com/vithubati/go-notifier/store"
	"github.com/vithubati/go-notifier/store/mysql"
	"github.com/vithubati/go-notifier/webhook"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Webhook = "WEBHOOK"
)

// service implements a notifier service.
type service struct {
	opts  Opts
	store store.Store
}

// Opts configures the notifier service
type Opts struct {
	DeliveryInterval time.Duration
	Migrations       bool
	ConnString       string
	Client           *http.Client
	WebhookEnabled   bool
}

// New creates a New notifier service and kicks off the notification deliveries
//
// Canceling the ctx will kill any concurrent routines affiliated with
// the notifier.
func New(opts Opts) (*service, error) {

	// initialize DB
	db, err := initDBConn(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DB: %v", err)
	}
	// initialize store
	s := initStore(db)

	return &service{
		store: s,
		opts:  opts,
	}, nil
}

func (s *service) CreateDeliverer(ctx context.Context, d model.Deliverer) error {
	if d.IntervalInSeconds <= 0 {
		d.IntervalInSeconds = int(s.opts.DeliveryInterval.Seconds())
	}
	return s.store.CreateDeliverer(ctx, d)
}

func (s *service) CreateNotification(ctx context.Context, n model.Notification) error {
	return s.store.Create(ctx, n)
}

// KickOff - kick off configured deliverer type
func (s *service) KickOff(ctx context.Context) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"component": "service/service.Deliver",
	}).WithContext(ctx)
	ctxLog.Info("Kicking Off...")
	switch {
	case s.opts.WebhookEnabled:
		if err := webhookDeliveries(ctx, s.opts, s.store); err != nil {
			return err
		}
	}
	return nil
}

func initDBConn(opts Opts) (*sql.DB, error) {
	db, err := sql.Open("mysql", opts.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	// do migrations if requested
	if opts.Migrations {
		err := migrate.Exec(db, migrate.Up, migrations.Migrations...)
		if err != nil {
			return nil, fmt.Errorf("failed to perform migrations: %w", err)
		}
	}
	return db, nil
}

func initStore(db *sql.DB) store.Store {
	return mysql.NewStore(db)
}

func webhookDeliveries(ctx context.Context, opts Opts, store store.Store) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"component": "service/service.webhookDeliveries",
	}).WithContext(ctx)
	ctxLog.Infof("initiating webhook deliverers")
	deliverers, err := store.GetDeliverer(ctx, Webhook)
	if err != nil {
		return fmt.Errorf("failed to get Delivery: %v", err)
	}
	ds := make([]*delivery.Delivery, 0, len(deliverers))
	for _, d := range deliverers {
		conf := &webhook.Config{
			Headers: http.Header(d.Headers),
			Client:  opts.Client,
			Target:  d.Url,
		}
		wh, err := webhook.New(conf)
		if err != nil {
			return fmt.Errorf("failed to create webhook deliverer: %v", err)
		}
		delivry := delivery.NewDelivery(d.ID, wh, time.Duration(d.IntervalInSeconds)*time.Second, store)
		ds = append(ds, delivry)
	}
	for _, d := range ds {
		// fixme create new context so the deliverer can be canceled when required
		d.Deliver(ctx)
	}
	return nil
}
