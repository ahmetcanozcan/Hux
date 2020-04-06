package hux

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configs : configurations for `Hux`
var Configs = struct {
	url             string
	DefaultRoomName string
}{
	url:             "/ws/hux",
	DefaultRoomName: "main",
}

// ErrRead : when a socket reads error except connection close
var ErrRead error = errors.New("socket : read error ")

// ErrSocketConnection :  Socket connection error
var ErrSocketConnection error = errors.New("connection : socket connection error")

func init() {

	// Set handle function
	http.HandleFunc(Configs.url, func(w http.ResponseWriter, r *http.Request) {
		// Get websocket connection 'c'
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(ErrSocketConnection)
		}

		// Create a new socket
		socket := newSocket(c)

		// Add socket to main room by default
		socket.Join(hub.GetRoom(Configs.DefaultRoomName))
		// Send socket connection signal to socket connectin channel
		hub.SocketConnection <- socket

		// Make sure closing connection
		defer socket.Disconnect()
		//Read messages from socket
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				// if err not a connection close then log it.
				if err != websocket.ErrCloseSent {
					log.Println(ErrRead, err)
				}
				break
			}
			er := socket.handleClientMessage(string(message))
			if er != nil {
				log.Println(err)
			}
		}

	})
}
