package handlers

import (
	"context"
	"regexp"

	"github.com/nlopes/slack"
	"github.com/spy16/skit"
)

// Echo echo's back messages recieved from slack if the message
// matches one of the regular expressions.
func Echo(exps ...string) (skit.Handler, error) {
	rexps := []*regexp.Regexp{}
	for _, exp := range exps {
		rexp, err := regexp.Compile(exp)
		if err != nil {
			return nil, err
		}
		rexps = append(rexps, rexp)
	}

	handler := skit.HandlerFunc(func(sk *skit.Skit, ev *slack.MessageEvent) bool {
		for _, rexp := range rexps {
			if rexp.Match([]byte(ev.Text)) {
				sk.SendText(context.Background(), ev.Text, ev.Channel)
				return true
			}
		}

		return false
	})
	return handler, nil
}
