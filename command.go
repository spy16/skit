package skit

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

// CommandHandler executes a configured command when input message matches one of
// the regex patterns. cmd and args both will be parsed as golang text templates.
// Named captures from pattern that matches the input message will be used as data
// for rendering command name and every argument.
func CommandHandler(lg Logger, cmd string, args []string, patterns []string) (*Command, error) {
	exps, err := ParseExprs(patterns)
	if err != nil {
		return nil, err
	}

	cmdTpl, err := template.New("command").Parse(cmd)
	if err != nil {
		return nil, err
	}

	argTpls := []template.Template{}
	for i, arg := range args {
		tpl, err := template.New(fmt.Sprintf("arg%d", i)).Parse(arg)
		if err != nil {
			return nil, err
		}
		argTpls = append(argTpls, *tpl)
	}

	cmdH := &Command{}
	cmdH.Logger = lg
	cmdH.cmd = *cmdTpl
	cmdH.args = argTpls
	cmdH.exprs = exps
	return cmdH, nil
}

// Command runs configured command when the message matches the regular
// expressions.
type Command struct {
	Logger
	Timeout     time.Duration
	RedirectErr bool
	WorkingDir  string

	cmd   template.Template
	args  []template.Template
	exprs []*regexp.Regexp
}

// Handle executes the command when the message matches the regular expressions.
func (cmd *Command) Handle(ctx context.Context, sk *Skit, ev *MessageEvent) bool {
	for _, expr := range cmd.exprs {
		match := CaptureAll(expr, ev.Text)
		if match == nil {
			continue
		}
		match["event"] = *ev

		out, err := cmd.executeCmd(ctx, match)
		if err != nil {
			msg := fmt.Sprintf("I fucked up :sob: (%v):\n%s", err, string(out))
			if err == context.Canceled {
				msg = fmt.Sprintf("I was interrupted :face_with_symbols: : \n%s", string(out))
			}
			sk.SendText(ctx, msg, ev.Channel)
		}

		sk.SendText(ctx, string(out), ev.Channel)
		return true
	}
	return false
}

func (cmd *Command) executeCmd(ctx context.Context, match map[string]interface{}) ([]byte, error) {
	if cmd.Timeout.Seconds() == 0 {
		cmd.Timeout = 1 * time.Minute
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, cmd.Timeout)
	defer cancel()

	execCmd, err := makeCmd(timeoutCtx, cmd.cmd, cmd.args, match)
	if err != nil {
		return nil, err
	}
	execCmd.Dir = cmd.WorkingDir

	var out []byte

	if cmd.RedirectErr {
		out, err = execCmd.CombinedOutput()
	} else {
		out, err = execCmd.Output()
	}

	if ctxErr := timeoutCtx.Err(); ctxErr == context.Canceled || ctxErr == context.DeadlineExceeded {
		return out, context.Canceled
	}

	return out, err
}

func makeCmd(ctx context.Context, cmd template.Template, args []template.Template, match map[string]interface{}) (*exec.Cmd, error) {
	cmdName, err := Render(cmd, match)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render command name")
	}

	cmdArgs, err := RenderAll(args, match)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render args")
	}

	execCmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	return execCmd, nil
}
