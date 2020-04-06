package hux

import (
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
