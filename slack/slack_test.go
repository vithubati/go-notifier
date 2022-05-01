package slack

import (
	"context"
	"github.com/vithubati/go-notifier/model"
	"testing"
	"time"
)

const token = "xoxb-1336502280788-3459693313171-s9dcsGWsZyEOoxIaiM8ADWz7"

func TestDeliver(t *testing.T) {
	if token == "" {
		return
	}
	conf := Config{
		Token:     token,
		ChannelID: "C03DAJMV62J",
	}

	d, err := New(&conf)
	if err != nil {
		t.Fatalf("New failed:: %v", err)
	}
	tme, _ := time.Parse(time.RFC3339, "2022-04-18 02:58:41")
	n := model.Notification{
		Topic:               "VPC",
		Action:                 "UPDATED",
		Subject:                "VPC Updated",
		Message:                "Vpc is updated for the accountId 3290",
		CreatedAt:              tme,
		Data:                   nil,
		NotificationDeliveryID: "",
	}
	err = d.Deliver(context.Background(), n)
	if err != nil {
		t.Fatalf("New failed:: %v", err)
	}
}
