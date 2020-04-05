package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetcanozcan/hux"
)

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/",fs)
	//Initialize hux
	hux.Initialize()
	// Get hub
	hub := hux.GetHub()
	// Handle hub
	go func(){
		for {
			select{
			case sck :=<-hub.SocketConnection:
				fmt.Println("Socket connected.")
				go handleSocket(sck)
			case <-hub.SocketDisconnection:
				fmt.Println("Socket disconnected.")
			}
		}
	}()

	// Start listening port
	http.ListenAndServe(":8080", nil)
}

func handleSocket(s *hux.Socket){
	for {
		select {
		case roomName :=<-s.GetEvent("Join"):
			r := hux.GetHub().GetRoom(roomName)
			s.Join(r)
			r.Emit("New","New Socket Connected")
			
		}
	}
}