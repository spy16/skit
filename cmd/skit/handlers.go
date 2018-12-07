package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spy16/skit"
	"github.com/spy16/skit/handlers"
	"github.com/spy16/skit/handlers/command"
)

var hs = map[string]makeFunc{
	"echo": func(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
		var echoCfg struct {
			Match []string
		}
		if err := mapstructure.Decode(cfg, &echoCfg); err != nil {
			return nil, err
		}
		return handlers.Echo(echoCfg.Match...)
	},
	"command": func(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
		cc := command.Config{}
		if err := mapstructure.Decode(cfg, &cc); err != nil {
			return nil, err
		}
		return command.New(cc, lg)
	},
}
