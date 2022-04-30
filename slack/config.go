package slack

import (
	"errors"
	"net/http"
)

// Config - Slack Config configures the "Slack" notification properties.
type Config struct {
	// API token
	Token string `json:"token"`

	// ChannelID - to send the notification to
	ChannelID string `json:"channelId"`

	// Client [Optional] provide a custom http client to the slack client.
	Client *http.Client
}

// Validate will validate the configuration properties
func (c *Config) validate() error {

	if c.Token == "" {
		return errors.New("API Token is not provided")
	}
	if c.ChannelID == "" {
		return errors.New("channel ID is not provided")
	}

	return nil
}
