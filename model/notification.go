package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Notification abstracts details of a change
// made on a specific Topic
type Notification struct {
	ID                     string      `json:"id"`
	Topic                  string      `json:"topic"`
	Action                 string      `json:"action"`
	Subject                string      `json:"subject"`
	Message                string      `json:"message"`
	CreatedAt              time.Time   `json:"createdAt"`
	Data                   interface{} `json:"data"`
	NotificationDeliveryID string      `json:"-"`
}

type Deliverer struct {
	ID                string           `json:"id"`
	Type              string           `json:"type"`
	Url               string           `json:"url"`
	ChannelID         string           `json:"channelId"`
	Credentials       string           `json:"credentials"`
	CreatedAt         time.Time        `json:"createdAt"`
	Retry             int              `json:"retry"`
	IntervalInSeconds int              `json:"IntervalInSeconds"`
	Headers           Headers          `json:"headers"`
	Topics            []DelivererTopic `json:"topics"`
}

type Headers map[string][]string

// Scan scans value into Json, implements sql.Scanner interface
func (h *Headers) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	byt, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	var hdrs map[string][]string
	err := json.Unmarshal(byt, &hdrs)
	if err != nil {
		return err
	}
	*h = hdrs
	return nil
}

// Value return json value, implement driver.Valuer interface
func (h Headers) Value() (driver.Value, error) {
	raw, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return driver.Value(raw), nil
}

type DelivererTopic struct {
	DelivererID string `json:"delivererId"`
	Topic       string `json:"topic"`
}

type Delivery struct {
	ID             string    `json:"id"`
	NotificationID string    `json:"notificationId"`
	DelivererID    string    `json:"delivererId"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
