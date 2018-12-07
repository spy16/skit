package command

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spy16/skit"
)

// New initializes the command handler with given configuration.
func New(cfg Config, lg skit.Logger) (*Command, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.setDefaults()

	exps := []*regexp.Regexp{}
	for _, exp := range cfg.Match {
		rex, err := regexp.Compile(exp)
		if err != nil {
			return nil, err
		}

		exps = append(exps, rex)
	}

	cmd := &Command{}
	cmd.cfg = cfg
	cmd.exps = exps
	cmd.Logger = lg
	return cmd, nil
}

// Command is a skit handler that runs configured commands.
type Command struct {
	skit.Logger

	exps []*regexp.Regexp
	cfg  Config
}

// Handle parses the message according to the regex and runs the command if the
// regex matches. Returns false if the message did not match the regex.
func (handler Command) Handle(ctx context.Context, sk *skit.Skit, ev *skit.MessageEvent) bool {
	for _, rexp := range handler.exps {
		matches := rexp.FindAllString(ev.Text, -1)
		if matches == nil {
			continue
		}

		handler.Debugf("expression '%s' matched", rexp.String())

		cmdCtx, cancel := context.WithTimeout(ctx, handler.cfg.timeout)
		defer cancel()

		args := renderArgs(handler.cfg.Args, matches)
		handler.Debugf("running '%s %v' with timeout '%s'", handler.cfg.Cmd, args, handler.cfg.timeout)
		cmd := exec.CommandContext(cmdCtx, handler.cfg.Cmd, args...)
		cmd.Dir, _ = filepath.Abs(handler.cfg.ConfigPath)
		out, err := cmd.CombinedOutput()
		if err != nil {

			msg := fmt.Sprintf("I fucked up:\n%s\nerr: %v", out, err)

			if cmdCtx.Err() == context.DeadlineExceeded {
				msg = fmt.Sprintf(":sob: I couldn't do it in time:\n%s\nerr: %v", out, err)
			}
			sk.SendText(ctx, msg, ev.Channel)
			return true
		}

		sk.SendText(ctx, string(out), ev.Channel)
		return true
	}

	return false
}
