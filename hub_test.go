package hux

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketConnection(t *testing.T) {
	done := make(chan bool)
	Initialize() // Initialize the app
	// Listen a socket connected or not
	go func() {
		for {
			select {
			case sck := <-hub.SocketConnection:
				done <- true
				_ = <-sck.GetEvent("Test")
				done <- true
			case <-hub.SocketDisconnection:
				done <- true
			}
		}
	}()
	// listen port
	go http.ListenAndServe(":8080", nil)

	// Deal and send hello message
	func() {
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/hux"}

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()
		c.WriteMessage(websocket.BinaryMessage, []byte(stringifyMessage("Test", "test-data")))
	}()

	<-done // Wait for connection
	<-done // Wait for Test message
	<-done // Wait for disconnection

	// if done channel gets signal three times in 30 seconds, test is passed
}
