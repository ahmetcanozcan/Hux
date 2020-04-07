package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetcanozcan/hux"
)

func main() {
	hub := hux.NewHub()
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/ws/hux", func(w http.ResponseWriter, r *http.Request) {
		socket, err := hub.InstantiateSocket(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("New Socket Connected.")
		for {
			select {
			case msg := <-socket.GetEvent("Join"):
				m := msg.String()
				fmt.Println("Join:", m)
				hub.GetRoom(m).Add(socket)
				hub.GetRoom(m).Emit("New", "NEW SOCKET CONNECTED.")

			}
		}
	})
	// Start listening port
	http.ListenAndServe(":8080", nil)
}
