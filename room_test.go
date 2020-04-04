package hux

import "testing"

func TestJoin(t *testing.T) {

	// Create a room
	room := NewRoom()
	// Try to join the room concurently
	for i := 0; i < 30; i++ {
		go func() {
			sck := newSocket(nil)
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
	for i := 0; i < 30; i++ {
		room.Add(newSocket(nil))
	}

	for sck := range room.sockets {
		go func(s *Socket) {
			room.Remove(s)
		}(sck)
	}

}
