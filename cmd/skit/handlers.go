package main

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spy16/skit"
)

// Mapping of handler-type to handler setup function.
var handlerMap = map[string]makeFunc{
	"simple":  makeSimpleHandler,
	"command": makeCommandHandler,
}

func makeSimpleHandler(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
	var conf struct {
		Match   []string
		Message string
	}
	if err := mapstructure.Decode(cfg, &conf); err != nil {
		return nil, err
	}
	return skit.SimpleHandler(lg, conf.Match, conf.Message)
}

func makeCommandHandler(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
	var cc struct {
		Cmd         string
		Args        []string
		Match       []string
		ConfigFile  string
		ConfigPath  string
		RedirectErr bool
		Timeout     string
	}
	if err := mapstructure.Decode(cfg, &cc); err != nil {
		return nil, err
	}

	cmdH, err := skit.CommandHandler(lg, cc.Cmd, cc.Args, cc.Match)
	if err != nil {
		return nil, err
	}

	dur, err := time.ParseDuration(cc.Timeout)
	if err != nil {
		return nil, err
	}

	cmdH.RedirectErr = cc.RedirectErr
	cmdH.Timeout = dur
	cmdH.WorkingDir = cc.ConfigPath
	return cmdH, nil
}
