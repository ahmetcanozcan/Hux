package hux

import (
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

// Configs : configurations for `Hux`
var Configs = struct {
	DefaultRoomName string
}{
	DefaultRoomName: "main",
}

// ErrRead : when a socket reads error except connection close
var ErrRead error = errors.New("socket : read error ")

// ErrSocketConnection :  Socket connection error
var ErrSocketConnection error = errors.New("connection : socket connection error")

// Upgrader :
var Upgrader = websocket.Upgrader{} // Create upgrader with default values.

// Hub :
type Hub struct {
	rooms    map[string]*Room
	mapMutex *sync.Mutex
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

// NewHub :
func NewHub() *Hub {
	return &Hub{
		rooms:    make(map[string]*Room),
		mapMutex: &sync.Mutex{},
	}
}
