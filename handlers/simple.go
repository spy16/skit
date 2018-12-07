package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/spy16/skit"
)

// Simple handler responds with a simple message if the input message matches one
// of the regular expressions.
func Simple(exps []string, tpl string) (skit.Handler, error) {
	rexps, err := skit.ParseExprs(exps)
	if err != nil {
		return nil, err
	}

	handler := skit.HandlerFunc(func(ctx context.Context, sk *skit.Skit, ev *skit.MessageEvent) bool {
		for _, rexp := range rexps {
			matches := getParams(rexp, ev.Text)
			if matches == nil {
				continue
			}

			msg := tpl
			for key, val := range matches {
				msg = strings.Replace(msg, fmt.Sprintf("{%s}", key), val, -1)
			}

			if rexp.Match([]byte(ev.Text)) {
				sk.SendText(ctx, msg, ev.Channel)
				return true
			}
		}

		return false
	})
	return handler, nil
}

func getParams(regEx *regexp.Regexp, url string) map[string]string {
	match := regEx.FindStringSubmatch(url)
	if match == nil {
		return nil
	}

	paramsMap := map[string]string{}
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
