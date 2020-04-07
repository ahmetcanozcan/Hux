package hux

import (
	"log"

	"github.com/gorilla/websocket"
)

// ConnReaderWriter :
type connReaderWriterCloser interface {
	WriteMessage(int, []byte) error
	ReadMessage() (int, []byte, error)
	Close() error
}

// Event :
type Event chan Message

//Socket : Socket is a websocket conection handler
type Socket struct {
	events map[string]Event
	conn   connReaderWriterCloser
	emitCh chan string
	room   *Room
}

func newSocket(c connReaderWriterCloser) *Socket {
	s := &Socket{
		events: make(map[string]Event),
		conn:   c,
		emitCh: make(chan string),
	}
	go s.run()
	return s
}

func (s *Socket) run() {
	go func() {
		for {
			select {
			case msg := <-s.emitCh:
				s.conn.WriteMessage(websocket.BinaryMessage, []byte(msg))
			}
		}
	}()
	go func() {
		for {
			msg, err := s.readMessage()
			if err != nil {
				// if err not a connection close then log it.
				if err != websocket.ErrCloseSent {
					log.Println(ErrRead, msg, err)
				}
				break
			}
			er := s.handleClientMessage(msg)
			if er != nil {
				log.Println(err)
			}
		}
	}()
}

func (s *Socket) handleClientMessage(rawStr string) error {
	// Splite raw string to name message
	name, message, err := parseRawMessage(rawStr)

	if err != nil {
		return err
	}
	ch := s.GetEvent(name)
	var msg Message = newMessage(message)
	ch <- msg
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

//GetEvent : Returns an event by given name. if event is not defined then defines and returns it
func (s *Socket) GetEvent(name string) Event {
	ch, ok := s.events[name]
	if !ok {
		s.events[name] = make(chan Message)
		ch = s.events[name]
	}
	return ch
}

// Emit : Send a message to client.
func (s *Socket) Emit(name string, data interface{}) {
	msg, _ := newSocketMessage(name, data).stringify()
	s.emitCh <- msg
}

// Broadcast : Send message all clients except its client
func (s *Socket) Broadcast(name string, data interface{}) {
	sl := s.room.sockets
	for sck := range sl {
		if sck != s {
			sck.Emit(name, data)
		}
	}
}

// Disconnect : Close connection
func (s *Socket) Disconnect() {
	s.conn.Close()
	s.GetEvent("Disonnection") <- ""
}

func (s *Socket) readMessage() (string, error) {
	_, msg, err := s.conn.ReadMessage()
	return string(msg), err
}
