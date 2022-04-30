package mysql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vithubati/go-notifier/model"
	"testing"
)

const data = `"{\n  \"type\": \"SERVER\",\n  \"id\": \"1\",\n  \"attributes\": {\n    \"Name\": \"ProdServer\",\n    \"region\": \"Canadacentral\",\n    \"cidr\": \"10.0.0.0/16\"\n  },\n  \"subnets\": [\n    {\n      \"subnet\": {\n        \"cidr\": \"10.0.1.0/24\"\n      }\n    }\n  ]\n}"`

func Test_createNotification(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	n := model.Notification{
		Resource: "SERVER",
		Action:   "CREATE",
		Subject:  "SERVER Created",
		Message:  "Server is created for the accountId G445",
		Data:     data,
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

func Test_createDeliverer2(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	d := model.Deliverer{
		Type:              "SLACK",
		ChannelID:         "C03DAJMV62J",
		Credentials:       "xoxb-1336502280788-3459693313171-cP9Kc3UQ93gfDiBThv5ubNvK",
		Retry:             3,
		IntervalInSeconds: 10,
		Resources:         []model.DelivererResource{{Resource: "SERVER"}},
	}
	err = createDeliverer(ctx, db, d)
	assert.Nil(t, err)
}
