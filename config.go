package skit

import (
	"errors"
	"strings"
)

// Config stores configuration options for skit.
type Config struct {
	Token string `json:"token" yaml:"token"`
}

// Validate performs basic validation of configuration.
func (cfg Config) Validate() error {
	if len(strings.TrimSpace(cfg.Token)) == 0 {
		return errors.New("empty token")
	}

	return nil
}
