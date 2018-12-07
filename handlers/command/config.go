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
	RedirectErr bool
	Timeout     string

	timeout time.Duration
}

// Validate required configuration options.
func (cfg *Config) Validate() error {
	if len(cfg.Cmd) == 0 {
		return errors.New("cmd cannot be empty")
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return err
	}
	cfg.timeout = timeout
	return nil
}

func (cfg *Config) setDefaults() {

}
