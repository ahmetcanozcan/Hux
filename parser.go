package hux

import (
	"fmt"
	"strings"
)

const messageSeperator = "]];|;[["

func parseRawMessage(rawMsg string) (string, string, bool) {
	sarr := strings.Split(rawMsg, messageSeperator)
	if len(sarr) != 2 {
		return "", "", false
	}
	return sarr[0], sarr[1], true
}

func stringifyMessage(name string, data string) string {
	return fmt.Sprintf("%s%s%s", name, messageSeperator, data)
}
