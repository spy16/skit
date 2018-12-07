package command

import (
	"fmt"
	"strings"
)

func renderArgs(tpls []string, data []string) []string {
	args := []string{}
	for _, tpl := range tpls {
		for i := 0; i < len(data); i++ {
			arg := strings.Replace(tpl, fmt.Sprintf("$%d", i), data[i], -1)
			args = append(args, arg)
		}
	}
	return args
}
