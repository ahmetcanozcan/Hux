package hux

import (
	"testing"
)

const sckCount = 50

func TestJoin(t *testing.T) {

	// Create a room
	room := NewRoom()
	// Try to join the room concurently
	for i := 0; i < sckCount; i++ {
		go func() {
			sck := newSocket(newTestConnection(""))
			sck.Join(room)
		}()
	}
	// If socket map do not generate a nil pointer error,
	// test is passed.
}

func TestLeave(t *testing.T) {
	// Create a room
	room := NewRoom()
	// Add sockets to the room
	for i := 0; i < sckCount; i++ {
		room.Add(newSocket(newTestConnection("")))
	}

	for sck := range room.sockets {
		go func(s *Socket) {
			room.Remove(s)
		}(sck)
	}

}
