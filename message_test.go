package hux

import (
	"strings"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
}

func TestParsing(t *testing.T) {
	strVal := "Test String"
	msg := newSocketMessage("Test", strVal)
	s, _ := msg.stringify()

	if v := strings.Split(s, messageSeperator)[1]; v != "\""+strVal+"\"" {
		t.Log("Wanted :", strVal, "Found:", v)
		t.FailNow()
	}

	objVal := testStruct{
		Name: "ahmet",
	}
	s, _ = newSocketMessage("Test", objVal).stringify()
	if v := strings.Split(s, messageSeperator)[1]; v != "{\"name\":\"ahmet\"}" {
		t.Log("Wanted :", strVal, "Found:", v)
		t.FailNow()
	}

}
