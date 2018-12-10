package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spy16/skit"
)

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
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err)
		os.Exit(1)
	}

	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)
	viper.ReadInConfig()

	cfg := config{
		Token:              viper.GetString("TOKEN"),
		NoHandler:          viper.GetString("NO_HANDLER"),
		RouteGroupMessages: viper.GetBool("ROUTE_GROUP_MESSAGES"),
		LogLevel:           viper.GetString("LOG_LEVEL"),
		LogFormat:          viper.GetString("LOG_FORMAT"),
		Handlers:           loadHandlerConfigs(configFile),
	}

	return cfg
}

type config struct {
	Token              string
	NoHandler          string
	RouteGroupMessages bool
	LogLevel           string
	LogFormat          string
	Handlers           []map[string]interface{}
}

func loadHandlerConfigs(configFile string) []map[string]interface{} {
	var hCfg struct {
		Handlers []map[string]interface{} `toml:"handlers,omitempty"`
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("failed to load config file: %s\n", err)
		os.Exit(1)
	}

	if err := toml.Unmarshal(data, &hCfg); err != nil {
		fmt.Printf("failed to read toml config file: %s\n", err)
		os.Exit(1)
	}

	configDir, err := filepath.Abs(filepath.Dir(configFile))
	if err != nil {
		fmt.Printf("failed to get absolute path of config file parent: %v", err)
		os.Exit(1)
	}

	for i := range hCfg.Handlers {
		if hCfg.Handlers[i] == nil {
			hCfg.Handlers[i] = map[string]interface{}{}
		}

		hCfg.Handlers[i]["configFile"] = configFile
		hCfg.Handlers[i]["configPath"] = configDir
	}

	return hCfg.Handlers
}
