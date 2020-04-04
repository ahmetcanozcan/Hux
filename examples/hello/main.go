package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetcanozcan/hux"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	hux.Initialize()
	go handleHub()
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}

func handleHub() {
	h := hux.GetHub()
	for {
		select {
		case sck := <-h.SocketConnection:
			go handleSocket(sck)
		case sck := <-h.SocketDisconnection:
			go handleDisconnection(sck)
		}
	}
}

func handleDisconnection(sck *hux.Socket) {
	fmt.Println("Socket disconnected")
}

func handleSocket(sck *hux.Socket) {
	fmt.Println("Socket connected.")
	for {
		select {
		case data := <-sck.GetEvent("Hello"):
			fmt.Println(data)
			sck.Emit("Hello", "Hello There!")
		case data := <-sck.GetEvent("Sum"):
			fmt.Println("Get Sum event", data)
			sck.Emit("Sum", "idk")
		}
	}
}
