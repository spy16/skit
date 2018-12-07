package skit

import "regexp"

// ParseExprs parses given list of regular expressions and returns the
// compiled objects.
func ParseExprs(exprs []string) ([]*regexp.Regexp, error) {
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
