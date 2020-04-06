package hux

import (
	"testing"
)

type TestConnection struct {
	msg string
	w   chan bool
}

func newTestConnection(m string) *TestConnection {
	return &TestConnection{
		msg: m,
		w:   make(chan bool),
	}
}

func (te *TestConnection) WriteMessage(int, []byte) error {
	<-te.w
	return nil
}

func (te *TestConnection) ReadMessage() (int, []byte, error) {
	<-te.w
	return 0, []byte(te.msg), nil
}

func TestBasicEvent(t *testing.T) {
	// Create socket
	s := &Socket{
		events: make(map[string]chan string),
		conn:   nil,
		room:   nil,
		emitCh: make(chan string),
	}

	done := make(chan bool)

	go func() {
		msg := <-s.GetEvent("foo")
		t.Log(msg)
		done <- true
	}()

	s.GetEvent("foo") <- "Hi"
	<-done

}

func TestEventsWithSelect(t *testing.T) {
	done := make(chan bool)
	s := &Socket{
		events: make(map[string]chan string),
		conn:   nil,
		room:   nil,
		emitCh: make(chan string),
	}

	go func() {
		for i := 0; i < 3; {
			select {
			case msg := <-s.GetEvent("foo"):
				t.Log("foo:", msg)
				i++
			case msg := <-s.GetEvent("bar"):
				t.Log("bar:", msg)
				i++
			case msg := <-s.GetEvent("joe"):
				t.Log("joe:", msg)
				i++
			default:
				continue
			}
		}
		done <- true
	}()
	s.GetEvent("foo") <- "hi"
	s.GetEvent("bar") <- "hi"
	s.GetEvent("joe") <- "hi"

	<-done
}

func TestSocketJoin(t *testing.T) {
	// Create room
	room := NewRoom()
	s := &Socket{
		events: make(map[string]chan string),
		conn:   nil,
		room:   nil,
		emitCh: make(chan string),
	}

	s.Join(room)

	if room != s.room {
		t.Error("expect:", room, "found: ", s.room)
	}

	haveSck := false
	for sck := range room.sockets {
		if sck == s {
			haveSck = true
		}
	}
	// if socket is not in the room.sockets
	if !haveSck {
		t.Error("Socket not in room.sockets")
	}
}
