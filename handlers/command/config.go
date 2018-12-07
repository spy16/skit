package command

import (
	"errors"
	"time"
)

// Config represents config options for running external command.
type Config struct {
	ConfigPath  string
	Match       []string
	Cmd         string
	Args        []string
	Timeout     time.Duration
	RedirectErr bool
}

// Validate required configuration options.
func (cfg Config) Validate() error {
	if len(cfg.Cmd) == 0 {
		return errors.New("cmd cannot be empty")
	}
	return nil
}

func (cfg *Config) setDefaults() {
	if cfg.Timeout.Seconds() == 0 {
		cfg.Timeout = 1 * time.Minute
	}
}
