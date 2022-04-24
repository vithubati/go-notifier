package mysql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vithubati/go-notifier/model"
	"testing"
)

func Test_createNotification(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	n := model.Notification{
		Resource: "VPC",
		Action:   "UPDATED",
	}
	err = createNotification(ctx, db, n)
	assert.Nil(t, err)
}

func Test_createDeliverer(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	d := model.Deliverer{
		Type:              "WEBHOOK",
		Url:               "https://www.stackoverflow.com/ttest",
		Retry:             3,
		IntervalInSeconds: 10,
		Resources:         []model.DelivererResource{{Resource: "SERVER"}},
	}
	headers := make(map[string][]string)
	headers["X-Request-id"] = []string{"456456"}
	d.Headers = headers
	err = createDeliverer(ctx, db, d)
	assert.Nil(t, err)
}
