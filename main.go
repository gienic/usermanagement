package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/whitearmy/server"
	"golang.org/x/net/websocket"
)

func main() {

	s := server.NewServer("12345678")
	go s.Authorize()
	go s.Unregister()

	http.Handle("/", websocket.Handler(s.Serve))

	//http.Handle("auth", )

	fmt.Println("Listening...")
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
