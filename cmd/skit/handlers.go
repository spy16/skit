package main

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spy16/skit"
	"github.com/spy16/skit/handlers/lua"
)

// Mapping of handler-type to handler setup function.
var handlerMap = map[string]makeFunc{
	"simple":  makeSimpleHandler,
	"command": makeCommandHandler,
	"lua":     makeLuaHandler,
}

func makeLuaHandler(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
	var conf struct {
		Paths      []string
		Handler    string
		Source     string
		ConfigPath string
	}

	if err := mapstructure.Decode(cfg, &conf); err != nil {
		return nil, err
	}

	conf.Paths = append(conf.Paths, fmt.Sprintf("%s/?.lua", conf.ConfigPath))
	lh, err := lua.New(conf.Source, conf.Handler, conf.Paths)
	if err != nil {
		return nil, err
	}
	return lh, nil
}

func makeSimpleHandler(lg *logrus.Logger, cfg map[string]interface{}) (skit.Handler, error) {
	var conf struct {
		Match   []string
		Message string
	}
	if err := mapstructure.Decode(cfg, &conf); err != nil {
		return nil, err
	}
	return skit.SimpleHandler(conf.Message, conf.Match...)
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

	cmdH, err := skit.CommandHandler(cc.Cmd, cc.Args, cc.Match)
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
