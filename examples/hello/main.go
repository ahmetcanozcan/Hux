package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetcanozcan/hux"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws/mux", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Socket connected.")
		sck, _ := hux.GenerateSocket(w, r)
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

	})

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)

}

func handleDisconnection(sck *hux.Socket) {
	fmt.Println("Socket disconnected")
}
