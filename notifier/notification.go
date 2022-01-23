package notifier

import "time"

// Action indicates the action performed on the Resource
type Action string

const (
	Create Action = "CREATE"
	Update Action = "UPDATE"
	Delete Action = "DELETE"
)

// Notification abstracts details of a change
// made on a specific resource
type Notification struct {
	ID       string      `json:"id"`
	Resource string      `json:"resource"`
	Action   Action      `json:"action"`
	Time     time.Time   `json:"time"`
	Data     interface{} `json:"data"`
}

