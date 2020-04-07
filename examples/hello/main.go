package main

import (
	"fmt"
	"net/http"

	"github.com/ahmetcanozcan/hux"
)

type summer struct {
	Nums []float32 `json:"nums"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	hub := hux.NewHub()
	http.HandleFunc("/ws/hux", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Socket connected.")
		sck, _ := hub.InstantiateSocket(w, r)
		for {
			select {
			case data := <-sck.GetEvent("Hello"):
				s := data.String()
				fmt.Println(s)
				sck.Emit("Hello", "Hello There!")
			case data := <-sck.GetEvent("Sum"):
				var s summer
				data.ParseJSON(&s)
				var sum float32 = 0
				for _, num := range s.Nums {
					sum += num
				}
				sck.Emit("Sum", sum)
			}
		}

	})

	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)

}

func handleDisconnection(sck *hux.Socket) {
	fmt.Println("Socket disconnected")
}
