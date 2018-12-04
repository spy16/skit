package handlers

import (
	"context"
	"regexp"

	"github.com/nlopes/slack"
	"github.com/spy16/skit"
)

// Command runs the configured command when the input matches one of the
// regular expression patterns.
func Command(exprs []string, args []string) (skit.Handler, error) {
	rexps, err := parseExprs(exprs)
	if err != nil {
		return nil, err
	}

	handler := skit.HandlerFunc(func(sk *skit.Skit, ev *slack.MessageEvent) bool {
		for _, rexp := range rexps {
			if rexp.Match([]byte(ev.Text)) {
				sk.SendText(context.Background(), "I can't do this yet.", ev.Channel)
				return true
			}
		}

		return false
	})
	return handler, nil
}

func parseExprs(exprs []string) ([]*regexp.Regexp, error) {
	rexps := []*regexp.Regexp{}
	for _, exp := range exprs {
		rexp, err := regexp.Compile(exp)
		if err != nil {
			return nil, err
		}
		rexps = append(rexps, rexp)
	}

	return rexps, nil
}
