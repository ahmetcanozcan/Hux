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
	data := make(chan string)
	type Exam struct {
		Name string `json:"name"`
	}
	//Create a new hub
	hub := NewHub()
	testMsg := "test-data"
	http.HandleFunc("/ws/hux", func(w http.ResponseWriter, r *http.Request) {
		socket, err := hub.InstantiateSocket(w, r)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		for {
			select {
			case msg := <-socket.GetEvent("Test"):
				data <- msg.String()
			case msg := <-socket.GetEvent("Test-Obj"):
				var t Exam
				msg.ParseJSON(&t)
				data <- t.Name
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
		msg, _ := newSocketMessage("Test", testMsg).stringify()
		c.WriteMessage(websocket.BinaryMessage, []byte(msg))
		msg = <-data // Wait for Test message
		if msg != testMsg {
			t.Log("wanted:", testMsg, "got", msg)
			t.FailNow()
		}
		te := Exam{Name: "TestName"}
		sm := newSocketMessage("Test-Obj", te)
		if err != nil {
			t.Log("Cannot parse struct")
			t.FailNow()
		}
		msg, err = sm.stringify()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		c.WriteMessage(websocket.BinaryMessage, []byte(msg))
		msg = <-data // Wait for Test message
		if msg != "TestName" {
			t.Log("wanted:", testMsg, "got", msg)
			t.FailNow()
		}
		t.SkipNow()
	}()
	<-done
}
