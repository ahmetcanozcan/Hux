package sergo

import (
	"fmt"
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
	hub := &Hub{
		SocketConnection:    make(chan *Socket),
		SocketDisconnection: make(chan *Socket),
	}
	hub.rooms = make(map[string]*Room)
	hub.rooms["main"] = NewRoom()
	// Handle request.
	http.HandleFunc("/ws/sergo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
		}

		// Make sure closing socket
		defer func() {
			fmt.Println("Connectino closed")
			c.Close()
		}()

		socket := newSocket(c)
		hub.SocketConnection <- socket

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Println("recv:", string(message))
			socket.handleClientMessage(string(message))
		}

	})
	return hub
}
