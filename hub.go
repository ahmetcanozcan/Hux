package hux

import (
	"errors"
	"net/http"
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

// Hub :
type Hub struct {
	rooms       map[string]*Room
	mapMutex    *sync.Mutex
	Upgrader    websocket.Upgrader
	DefaultRoom *Room
}

// InstantiateSocket :
func (h *Hub) InstantiateSocket(w http.ResponseWriter, r *http.Request) (*Socket, error) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	s := newSocket(conn)
	return s, nil

}

// GetRoom : Returns a existing room. if it's not exist, creates and returns a new room
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

// Emit :
func (h *Hub) Emit(name string, data interface{}) {
	h.DefaultRoom.Emit(name, data)
}

// NewHub :
func NewHub() *Hub {
	h := &Hub{
		rooms:    make(map[string]*Room),
		mapMutex: &sync.Mutex{},
	}
	// Set default room
	h.DefaultRoom = h.GetRoom("main")
	return h
}
