package sergo

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

//Socket :
type Socket struct {
	events map[string]chan string
	conn   *websocket.Conn
	room   *Room
}

func newSocket(c *websocket.Conn) *Socket {
	s := &Socket{
		events: make(map[string]chan string),
		conn:   c,
	}
	return s
}

func (s *Socket) handleClientMessage(rawStr string) {
	name, message, ok := parseRawMessage(rawStr)
	if !ok {
		log.Println("Invalid Text", rawStr)
		return
	}
	ch := s.GetEvent(name)
	fmt.Println("name: ", name, "message: ", message, "ch", ch)
	ch <- message
}

// Join :
func (s *Socket) Join(r *Room) {
	// If have a room then disconnect from it
	if s.room != nil {
		delete(s.room.sockets, s)
	}
	r.Add(s)
}

//GetEvent :
func (s *Socket) GetEvent(name string) chan string {
	ch, ok := s.events[name]
	if !ok {
		s.events[name] = make(chan string)
		ch = s.events[name]
	}
	return ch
}

// Emit :
func (s *Socket) Emit(name string, data string) {
	msg := stringifyMessage(name, data)
	s.conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
}

// Broadcast :
func (s *Socket) Broadcast(name string, data string) {
	sl := s.room.sockets
	for sck := range sl {
		if sck != s {
			sck.Emit(name, data)
		}
	}

}
