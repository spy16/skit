package skit

import (
	"regexp"
	"strings"
	"text/template"
)

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

// RenderAll executes all templates passed in with the data and returns the rendered
// strings. Returns nil and an error on first execution failure.
func RenderAll(tpls []template.Template, data interface{}) ([]string, error) {
	out := []string{}
	for _, tpl := range tpls {
		s, err := Render(tpl, data)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	return out, nil
}

// Render executes the template with data and returns the rendered string.
func Render(tpl template.Template, data interface{}) (string, error) {
	wr := &strings.Builder{}
	if err := tpl.Execute(wr, data); err != nil {
		return "", err
	}
	return wr.String(), nil
}

// CaptureAll matches s with the regex and returns all named capture values.
// Returns nil if the s was not a match.
func CaptureAll(regEx *regexp.Regexp, s string) map[string]string {
	match := regEx.FindStringSubmatch(s)
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
