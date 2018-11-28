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

func main() {
	cfg := &skit.Config{}

	var logLevel string
	makeLogger := func() *logrus.Logger {
		logger := logrus.New()
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			lvl = logrus.InfoLevel
			logger.Warnf("invalid log level '%s', defaulting to info", logLevel)
		}
		logger.SetLevel(lvl)
		return logger
	}

	cmd := &cobra.Command{
		Use:   "skit",
		Short: "skit is a sick slack bot",
		Run: func(cmd *cobra.Command, args []string) {
			logger := makeLogger()
			sl, err := skit.New(*cfg, logger, skit.WithMessageHandler(func(sl *skit.Skit, ev *slack.MessageEvent) {
				sl.SendText(context.Background(), "`hello`", ev.Channel)
			}))
			if err != nil {
				logger.Fatalf("err: %s", err)
			}

			if err := sl.Listen(context.Background()); err != nil {
				logger.Fatalf("err: %s", err)
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Logging level")
	cmd.PersistentFlags().StringP("token", "t", "", "Slack access token")

	cmd.AddCommand(newConfigCmd(cfg))
	cobra.OnInitialize(func() {
		if err := loadConfig(cmd, cfg); err != nil {
			fmt.Println("failed to load config")
			os.Exit(1)
		}
	})
	cmd.Execute()
}

func loadConfig(cmd *cobra.Command, into *skit.Config) error {
	viper.AddConfigPath(".")
	viper.SetConfigName("skit")
	viper.BindPFlags(cmd.PersistentFlags())
	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfg := skit.Config{
		Token: viper.GetString("token"),
	}
	*into = cfg
	return nil
}

func newConfigCmd(cfg *skit.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Display current configuration",
	}
	cmd.Run = func(_ *cobra.Command, args []string) {
		yaml.NewEncoder(os.Stdout).Encode(cfg)
	}
	return cmd
}
