package slack

import (
	"errors"
)

// Config - Slack Config configures the "Slack" notification properties.
type Config struct {
	// API token
	Token string `json:"token"`

	ChannelID string `json:"channelId"`
}

// Validate will validate the configuration properties
func (c *Config) validate() error {

	if c.Token == "" {
		return errors.New("API Token is not provided")
	}

	return nil
}
