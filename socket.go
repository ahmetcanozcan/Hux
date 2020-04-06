package hux

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ConnReaderWriter :
type ConnReaderWriter interface {
	WriteMessage(int, []byte) error
	ReadMessage() (int, []byte, error)
}

//Socket : Socket is a websocket conection handler
type Socket struct {
	events map[string]chan string
	conn   ConnReaderWriter
	emitCh chan string
	room   *Room
}

func newSocket(c ConnReaderWriter) *Socket {
	s := &Socket{
		events: make(map[string]chan string),
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
					log.Println(ErrRead, err)
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
	s.GetEvent("Disonnection") <- ""
}

func (s *Socket) readMessage() (string, error) {
	_, msg, err := s.conn.ReadMessage()
	return string(msg), err
}

// GenerateSocket :
func GenerateSocket(w http.ResponseWriter, r *http.Request) (*Socket, error) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	s := newSocket(conn)
	return s, nil

}
