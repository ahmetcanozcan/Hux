package hux

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub :
type Hub struct {
	rooms               map[string]*Room
	SocketConnection    chan *Socket
	SocketDisconnection chan *Socket
	mapMutex            *sync.Mutex
}

// GetRoom :
func (h *Hub) GetRoom(name string) *Room {
	h.mapMutex.Lock()
	r, ok := h.rooms[name]
	if !ok {
		r = NewRoom()
		h.rooms[name] = r
	}
	h.mapMutex.Unlock()
	return r
}

var upgrader = websocket.Upgrader{} // Create upgrader with default values.

// Configs :
var Configs = struct {
	URL             string
	DefaultRoomName string
}{
	URL:             "/ws/hux",
	DefaultRoomName: "main",
}

var (
	hub = &Hub{
		rooms:               make(map[string]*Room),
		SocketConnection:    make(chan *Socket),
		SocketDisconnection: make(chan *Socket),
		mapMutex:            &sync.Mutex{},
	}
)

// GetHub :
func GetHub() *Hub {
	return hub
}

// Initialize :
func Initialize() {

	// Create default Room
	hub.rooms[Configs.DefaultRoomName] = NewRoom()
	// Set handle function
	http.HandleFunc(Configs.URL, func(w http.ResponseWriter, r *http.Request) {
		// Get websocket connection 'c'
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade err:", err)
		}

		socket := newSocket(c)

		// Add socket to main room by default
		socket.Join(hub.rooms[Configs.DefaultRoomName])
		// Send socket connection signal
		hub.SocketConnection <- socket

		// Make sure closing socket
		defer socket.Disconnect()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			log.Println("recv:", string(message))
			socket.handleClientMessage(string(message))
		}

	})
}
