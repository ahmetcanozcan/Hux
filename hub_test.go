package hux

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSocketConnection(t *testing.T) {

	done := make(chan string)
	//Create a new hub
	hub := NewHub()
	testMsg := "test-data"
	http.HandleFunc("/ws/hux", func(w http.ResponseWriter, r *http.Request) {
		socket, err := GenerateSocket(w, r)
		hub.GetRoom("main").Add(socket)
		if err != nil {
			t.FailNow()
		}
		for {
			select {
			case msg := <-socket.GetEvent("Test"):
				done <- msg
			}

		}
	})

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
		c.WriteMessage(websocket.BinaryMessage, []byte(stringifyMessage("Test", testMsg)))
	}()

	msg := <-done // Wait for Test message
	if msg != testMsg {
		t.FailNow()
	}
}
