package sergo

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Hub :
type Hub struct {
	rooms               map[string]*Room
	SocketConnection    chan *Socket
	SocketDisconnection chan *Socket
}

var upgrader = websocket.Upgrader{} // Create upgrader with default values.

// NewHub :
func NewHub() *Hub {
	// Instantiating.
	hub := &Hub{}
	hub.rooms = make(map[string]*Room)
	hub.rooms["main"] = NewRoom()
	// Handle request.
	http.HandleFunc("/ws/sergo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
		}

		socket := newSocket(c)
		// Make sure closing socket
		defer func() {
			hub.SocketDisconnection <- socket
			c.Close()
		}()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Println("recv:", message)
			socket.handleClientMessage(string(message))
		}

	})
	return hub
}
