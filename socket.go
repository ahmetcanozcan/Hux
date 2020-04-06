package hux

import (
	"github.com/gorilla/websocket"
)

//Socket : Socket is a websocket conection handler
type Socket struct {
	events map[string]chan string
	conn   *websocket.Conn
	emitCh chan string
	room   *Room
}

func newSocket(c *websocket.Conn) *Socket {
	s := &Socket{
		events: make(map[string]chan string),
		conn:   c,
		emitCh: make(chan string),
	}
	go s.run()
	return s
}

func (s *Socket) run() {
	for {
		select {
		case msg := <-s.emitCh:
			s.conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
		}
	}

}

func (s *Socket) handleClientMessage(rawStr string) error {
	// Splite raw string to name message
	name, message, err := parseRawMessage(rawStr)

	if err != nil {
		return err
	}
	ch := s.GetEvent(name)
	ch <- message
	return nil
}

// Join : adds itself to the given room
func (s *Socket) Join(r *Room) {
	// If have a room then disconnect from it
	if s.room != nil {
		s.LeaveRoom()
	}
	s.room = r
	r.Add(s)
}

// LeaveRoom : leaves its room
func (s *Socket) LeaveRoom() {
	s.room.Remove(s)
	s.room = nil

}

//GetEvent : Return a event by given name. if event is not defined then define and return it
func (s *Socket) GetEvent(name string) chan string {
	ch, ok := s.events[name]
	if !ok {
		s.events[name] = make(chan string)
		ch = s.events[name]
	}
	return ch
}

// Emit : Send a message to client.
func (s *Socket) Emit(name string, data string) {
	msg := stringifyMessage(name, data)
	s.emitCh <- msg
}

// Broadcast : Send message all clients except its client
func (s *Socket) Broadcast(name string, data string) {
	sl := s.room.sockets
	for sck := range sl {
		if sck != s {
			sck.Emit(name, data)
		}
	}
}

// Disconnect : Close connection
func (s *Socket) Disconnect() {
	hub.SocketDisconnection <- s
}
