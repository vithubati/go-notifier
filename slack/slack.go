package slack

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/vithubati/go-notifier/delivery"
	"github.com/vithubati/go-notifier/model"
)

// Compile-time check to ensure Slack implements delivery.Deliverer.
var _ delivery.Deliverer = (*Slack)(nil)

type Slack struct {
	client    *slack.Client
	channelID string
}

// New returns a new Slack Deliverer
func New(conf *Config) (delivery.Deliverer, error) {
	if conf == nil {
		return nil, errors.New("config not provided")
	}
	if err := conf.validate(); err != nil {
		return nil, err
	}
	client := slack.New(conf.Token)
	return &Slack{client: client, channelID: conf.ChannelID}, nil
}

func (s *Slack) Deliver(ctx context.Context, notification model.Notification) error {
	ctxLog := logrus.WithFields(logrus.Fields{
		"notificationId": notification.ID,
		"component":      "slack/deliverer.Deliver",
	}).WithContext(ctx)
	if len(s.channelID) == 0 {
		ctxLog.Infof("No channel id is registerered. Returning...")
		return nil
	}
	attachment := slack.Attachment{
		Pretext: notification.Message,
		Fields: []slack.AttachmentField{
			{
				Title: "Resource",
				Value: notification.Resource,
			},
			{
				Title: "Action",
				Value: notification.Action,
			},
			{
				Title: "Created at",
				Value: notification.CreatedAt.String(),
			},
		},
	}

	channelID, timestamp, err := s.client.PostMessageContext(
		ctx,
		s.channelID,
		slack.MsgOptionText(notification.Subject, false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return fmt.Errorf("slack post message was nit successful: %s", err.Error())
	}
	ctxLog.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}
