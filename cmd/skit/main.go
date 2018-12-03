package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spy16/skit"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Skit      skit.Config
	LogLevel  string
	LogFormat string
}

func main() {
	cfg := config{}

	cmd := &cobra.Command{
		Use:   "skit",
		Short: "skit is a sick slack bot",
	}

	cmd.PersistentFlags().StringP("log-level", "l", "info", "Logging level")
	cmd.PersistentFlags().StringP("token", "t", "", "Slack access token")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		logger := makeLogger(cfg.LogLevel)
		sl, err := skit.New(cfg.Skit, logger)
		if err != nil {
			logger.Fatalf("err: %s", err)
		}

		if err := sl.Listen(context.Background()); err != nil {
			logger.Fatalf("err: %s", err)
		}
	}

	cmd.AddCommand(newConfigCmd(cfg))

	cobra.OnInitialize(func() {
		if err := loadConfig(cmd, &cfg); err != nil {
			fmt.Println("failed to load config")
			os.Exit(1)
		}
	})
	cmd.Execute()
}

func loadConfig(cmd *cobra.Command, into *config) error {
	viper.AddConfigPath(".")
	viper.SetConfigName("skit")
	viper.BindPFlags(cmd.PersistentFlags())
	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfg := config{
		Skit: skit.Config{
			Token: viper.GetString("token"),
		},
		LogLevel:  viper.GetString("log-level"),
		LogFormat: viper.GetString("log-format"),
	}
	*into = cfg
	return nil
}

func newConfigCmd(cfg config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display current configuration",
	}
	cmd.Run = func(_ *cobra.Command, args []string) {
		yaml.NewEncoder(os.Stdout).Encode(cfg)
	}
	return cmd
}

func makeLogger(logLevel string) *logrus.Logger {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		lvl = logrus.InfoLevel
		logger.Warnf("invalid log level '%s', defaulting to info", logLevel)
	}
	logger.SetLevel(lvl)
	return logger
}

func onMessage(sl *skit.Skit, ev *slack.MessageEvent) {
	sl.SendText(context.Background(), "`hello`", ev.Channel)
}
