package hux

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var errParse error = errors.New("parse : string cannot be parsed")

const messageSeperator = "]];|;[["

// Message :
type Message string

func newMessage(s string) Message {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		s = s[1 : len(s)-1]
	}
	return Message(s)
}

// ParseJSON :
func (m Message) ParseJSON(v interface{}) error {
	return json.Unmarshal([]byte(m), v)
}

func (m Message) String() string {
	return string(m)
}

type socketMessage struct {
	name    string
	payload interface{}
}

func newSocketMessage(name string, payload interface{}) socketMessage {
	return socketMessage{
		name:    name,
		payload: payload,
	}
}

func (m socketMessage) stringify() (string, error) {
	data, err := json.Marshal(m.payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", m.name, messageSeperator, string(data)), nil
}

func parseRawMessage(rawMsg string) (string, string, error) {
	sarr := strings.Split(rawMsg, messageSeperator)
	if len(sarr) != 2 {
		return "", "", errParse
	}
	return sarr[0], sarr[1], nil
}
