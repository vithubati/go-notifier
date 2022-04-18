package model

import (
	"time"
)

// Notification abstracts details of a change
// made on a specific resource
type Notification struct {
	ID                     string      `json:"id"`
	Resource               string      `json:"resource"`
	Action                 string      `json:"action"`
	CreatedAt              time.Time   `json:"createdAt"`
	Data                   interface{} `json:"data"`
	NotificationDeliveryID string      `json:"-"`
}

type Deliverer struct {
	ID          string              `json:"id"`
	Type        string              `json:"type"`
	Url         string              `json:"url"`
	Credentials string              `json:"credentials"`
	CreatedAt   time.Time           `json:"createdAt"`
	Retry       int                 `json:"retry"`
	Interval    int                 `json:"interval"`
	Headers     map[string][]string `json:"headers"`
	Resources   []DelivererResource `json:"resources"`
}

type DelivererResource struct {
	DelivererID string `json:"delivererId"`
	Resource    string `json:"resource"`
}

type Delivery struct {
	ID             string    `json:"id"`
	NotificationID string    `json:"notificationId"`
	DelivererID    string    `json:"delivererId"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
