package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/whitearmy/user/server"
	"golang.org/x/net/websocket"
)

func main() {

	d := server.NewDisposer("12345678")
	go d.Authorize()
	go d.Unregister()

	http.Handle("/", websocket.Handler(d.Serve))

	fmt.Println("Listening...")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
