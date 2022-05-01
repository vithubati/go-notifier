package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/remind101/migrate"
	"github.com/sirupsen/logrus"
	"github.com/vithubati/go-notifier/config"
	slack2 "github.com/vithubati/go-notifier/deliverer/slack"
	webhook2 "github.com/vithubati/go-notifier/deliverer/webhook"
	"github.com/vithubati/go-notifier/delivery"
	"github.com/vithubati/go-notifier/migrations"
	"github.com/vithubati/go-notifier/model"
	"github.com/vithubati/go-notifier/store"
	"github.com/vithubati/go-notifier/store/mysql"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	Webhook = "WEBHOOK"
	Slack   = "SLACK"
)

// service implements a notifier service.
type service struct {
	cfg   *config.Config
	store store.Store
}

// New creates a New notifier service and kicks off the notification deliveries
//
// Canceling the ctx will kill any concurrent routines affiliated with
// the notifier.
func New(db *sql.DB, cfg *config.Config) (Service, error) {
	if err := cfg.Notifier.Validate(); err != nil {
		return nil, err
	}
	if cfg.Trace {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if cfg.JsonLogFormat {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	// do migrations if requested
	if cfg.Notifier.Migrations {
		if err := migrateDB(db); err != nil {
			return nil, fmt.Errorf("failed to initialize DB: %v", err)
		}
	}
	// initialize store
	s := initStore(db)

	return &service{
		store: s,
		cfg:   cfg,
	}, nil
}

func (s *service) CreateNotification(ctx context.Context, n model.Notification) error {
	return s.store.Create(ctx, n)
}

func (s *service) CreateDeliverer(ctx context.Context, d model.Deliverer) error {
	if d.IntervalInSeconds <= 0 {
		d.IntervalInSeconds = int(s.cfg.Notifier.DeliveryInterval.Seconds())
	}
	return s.store.CreateDeliverer(ctx, d)
}

// KickOff - kick off configured deliverer type
func (s *service) KickOff(ctx context.Context) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"component": "service/service.Deliver",
	}).WithContext(ctx)
	ctxLog.Info("Kicking Off...")
	if s.cfg.Notifier.Webhook {
		if err := webhookDeliveries(ctx, s.cfg.Notifier, s.store); err != nil {
			return err
		}
	}
	if s.cfg.Notifier.Slack {
		if err := slackDeliveries(ctx, s.cfg.Notifier, s.store); err != nil {
			return err
		}
	}
	return nil
}

func migrateDB(db *sql.DB) error {
	err := migrate.Exec(db, migrate.Up, migrations.Migrations...)
	if err != nil {
		return fmt.Errorf("failed to perform migrations: %w", err)
	}
	return nil
}

func initStore(db *sql.DB) store.Store {
	return mysql.NewStore(db)
}

func webhookDeliveries(ctx context.Context, opts config.Notifier, store store.Store) error {
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
		conf := &webhook2.Config{
			Headers: http.Header(d.Headers),
			Client:  opts.Client,
			Target:  d.Url,
		}
		wh, err := webhook2.New(conf)
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

func slackDeliveries(ctx context.Context, opts config.Notifier, store store.Store) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"component": "service/service.slackDeliveries",
	}).WithContext(ctx)
	ctxLog.Infof("initiating slack deliverers")
	deliverers, err := store.GetDeliverer(ctx, Slack)
	if err != nil {
		return fmt.Errorf("failed to get Delivery: %v", err)
	}
	ds := make([]*delivery.Delivery, 0, len(deliverers))
	for _, d := range deliverers {
		conf := &slack2.Config{
			Token:     d.Credentials,
			ChannelID: d.ChannelID,
			Client:    opts.Client,
		}
		wh, err := slack2.New(conf)
		if err != nil {
			return fmt.Errorf("failed to create slack deliverer: %v", err)
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
