package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func initDBConn() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@/notifier?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}
	return db, nil
}

func Test_getDeliverer(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	res, err := getDeliverers(ctx, db, "WEBHOOK")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_getDelivererIDsByTopic(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	res, err := getDelivererIDsByTopic(ctx, db, "VPC")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_getDelivererTopics(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	res, err := getDelivererTopics(ctx, db, "111")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_createdNotifications(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	res, err := createdNotifications(ctx, db, "111")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_failedNotifications(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	res, err := createdNotifications(ctx, db, "111")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}
