package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spy16/skit"
)

func main() {
	cmd := &cobra.Command{
		Use:   "skit",
		Short: "skit is a sick slack bot",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := loadConfig(cmd)
			logger := makeLogger(cfg.LogLevel, cfg.LogFormat)
			skCfg := skit.Config{
				Token: cfg.Token,
			}

			sl, err := skit.New(skCfg, logger)
			if err != nil {
				logger.Fatalf("err: %s", err)
			}
			registerHandlers(cfg.Handlers, sl, logger)

			if err := sl.Listen(context.Background()); err != nil {
				logger.Fatalf("err: %s", err)
			}
		},
	}
	cmd.PersistentFlags().StringP("config", "c", "skit.toml", "Configuration file path")

	cmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Display current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := loadConfig(cmd)
			toml.NewEncoder(os.Stdout).Encode(cfg)
		},
	})

	cmd.Execute()
}

func registerHandlers(handlers []map[string]interface{}, sk *skit.Skit, logger *logrus.Logger) {
	for index, cfg := range handlers {
		typ, ok := cfg["type"].(string)
		if !ok {
			logger.Fatalf("all handlers need 'type'")
		}

		name, ok := cfg["name"].(string)
		if !ok {
			logger.Fatalf("all handlers need 'name'")
		}

		name = strings.TrimSpace(name)
		typ = strings.TrimSpace(strings.ToLower(typ))

		maker, found := handlerMap[typ]
		if !found {
			logger.Fatalf("handler of type '%s' not found, exiting", typ)
		}

		h, err := maker(logger, cfg)
		if err != nil {
			logger.Fatalf("failed to init handler '%s', index %d: %s\n", typ, index, err)
		}

		sk.Register(name, h)
	}
}

type makeFunc func(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error)

func makeLogger(logLevel, logFormat string) *logrus.Logger {
	logger := logrus.New()

	if logLevel == "" {
		logLevel = "info"
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lvl = logrus.InfoLevel
		logger.Warnf("invalid log level '%s', defaulting to info", logLevel)
	}
	logger.SetLevel(lvl)

	if logFormat == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

func loadConfig(cmd *cobra.Command) config {
	cfg := config{}
	configFile, err := cmd.PersistentFlags().GetString("config")
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err)
		os.Exit(1)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("failed to read toml config file: %s\n", err)
		os.Exit(1)
	}

	configDir, err := filepath.Abs(filepath.Dir(configFile))
	if err != nil {
		fmt.Printf("failed to get absolute path of config file parent: %v", err)
		os.Exit(1)
	}

	for i := range cfg.Handlers {
		if cfg.Handlers[i] == nil {
			cfg.Handlers[i] = map[string]interface{}{}
		}

		cfg.Handlers[i]["configFile"] = configFile
		cfg.Handlers[i]["configPath"] = configDir
	}

	token := os.Getenv("TOKEN")
	if len(token) > 0 {
		cfg.Token = token
	}

	return cfg
}

type config struct {
	Token     string                   `toml:"token,omitempty"`
	LogLevel  string                   `toml:"log_level,omitempty"`
	LogFormat string                   `toml:"log_format,omitempty"`
	Handlers  []map[string]interface{} `toml:"handlers,omitempty"`
}
