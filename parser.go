package hux

import (
	"errors"
	"fmt"
	"strings"
)

const messageSeperator = "]];|;[["

// ErrParse :
var ErrParse error = errors.New("parse : string cannot parsed")

func parseRawMessage(rawMsg string) (string, string, error) {
	sarr := strings.Split(rawMsg, messageSeperator)
	if len(sarr) != 2 {
		return "", "", ErrParse
	}
	return sarr[0], sarr[1], nil
}

func stringifyMessage(name string, data string) string {
	return fmt.Sprintf("%s%s%s", name, messageSeperator, data)
}
