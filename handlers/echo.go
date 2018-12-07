package handlers

import (
	"context"

	"github.com/spy16/skit"
)

// Echo echo's back messages recieved from slack if the message
// matches one of the regular expressions.
func Echo(exps ...string) (skit.Handler, error) {
	rexps, err := skit.ParseExprs(exps)
	if err != nil {
		return nil, err
	}

	handler := skit.HandlerFunc(func(ctx context.Context, sk *skit.Skit, ev *skit.MessageEvent) bool {
		for _, rexp := range rexps {
			if rexp.Match([]byte(ev.Text)) {
				sk.SendText(ctx, ev.Text, ev.Channel)
				return true
			}
		}

		return false
	})
	return handler, nil
}
