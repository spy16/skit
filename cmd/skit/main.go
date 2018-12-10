package main

import (
	"context"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/spy16/skit"
)

func main() {
	cmd := &cobra.Command{
		Use:   "skit",
		Short: "skit is a sick slack bot",
		Long:  skitHelp,
		Run:   runSkit,
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

func runSkit(cmd *cobra.Command, args []string) {
	cfg := loadConfig(cmd)
	logger := makeLogger(cfg.LogLevel, cfg.LogFormat)
	sl := skit.New(cfg.Token, logger)

	if len(cfg.NoHandler) > 0 {
		tpl, err := template.New("simple").Parse(cfg.NoHandler)
		if err != nil {
			logger.Fatalf("invalid no_handler template: %v", err)
		}
		sl.NoHandler = *tpl
	}
	registerHandlers(cfg.Handlers, sl, logger)

	if err := sl.Listen(context.Background()); err != nil {
		logger.Fatalf("err: %s", err)
	}
}

const skitHelp = `
Build slack bots quickly and easily!

For more info, visit https://github.com/spy16/skit
`
