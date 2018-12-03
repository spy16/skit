package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spy16/skit"
	"github.com/spy16/skit/handlers"
	yaml "gopkg.in/yaml.v2"
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
	cmd.PersistentFlags().StringP("config", "c", "skit.yaml", "Configuration file path")

	cmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Display current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := loadConfig(cmd)
			yaml.NewEncoder(os.Stdout).Encode(cfg)
		},
	})

	cmd.Execute()
}

func registerHandlers(handlers []handler, sk *skit.Skit, logger *logrus.Logger) {
	for _, cfg := range handlers {
		ht := strings.TrimSpace(strings.ToLower(cfg.Type))
		maker, found := hs[ht]
		if !found {
			logger.Fatalf("handler of type '%s' not found, exiting", ht)
		}

		h, err := maker(cfg)
		if err != nil {
			logger.Fatalf("failed to init handler '%s': %s\n", ht, err)
		}

		sk.Register(h)
	}
}

var hs = map[string]makeFunc{
	"echo": func(cfg handler) (skit.Handler, error) {
		return handlers.Echo(cfg.Match...)
	},
}

type makeFunc func(cfg handler) (skit.Handler, error)

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

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("failed to read yaml config file: %s\n", err)
		os.Exit(1)
	}

	token := os.Getenv("TOKEN")
	if len(token) > 0 {
		cfg.Token = token
	}

	return cfg
}

type config struct {
	Token     string    `json:"token" yaml:"token"`
	LogLevel  string    `json:"log_level" yaml:"log_level"`
	LogFormat string    `json:"log_format" yaml:"log_format"`
	Handlers  []handler `json:"handlers" yaml:"handlers"`
}

type handler struct {
	Type  string   `json:"type" yaml:"type"`
	Match []string `json:"match" yaml:"match"`
	Cmd   []string `json:"cmd" yaml:"cmd"`
}
