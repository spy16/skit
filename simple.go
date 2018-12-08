package skit

import (
	"context"
	"fmt"
	"text/template"
)

// SimpleHandler responds with a simple message if the input message matches one
// of the regular expressions. The message (argument tplStr) can be a golang text
// template. Named captures from the regex that matches the input message will be
// used as data for rendering the message template.
func SimpleHandler(lg Logger, exps []string, tplStr string) (Handler, error) {
	rexps, err := ParseExprs(exps)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("simple").Parse(tplStr)
	if err != nil {
		return nil, err
	}

	handler := HandlerFunc(func(ctx context.Context, sk *Skit, ev *MessageEvent) bool {
		for _, rexp := range rexps {
			matches := CaptureAll(rexp, ev.Text)
			if matches == nil {
				continue
			}

			msg, err := Render(*tpl, matches)
			if err != nil {
				sk.SendText(ctx, fmt.Sprintf(":face_with_symbols: Something is not right: %v", err), ev.Channel)
				return true
			}

			sk.SendText(ctx, msg, ev.Channel)
			return true
		}

		return false
	})
	return handler, nil
}
