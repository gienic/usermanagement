package main

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type Message struct {
	Token   string
	Type    string
	Message string
	Error   string
}

func main() {

	origin := "http://localhost/"
	url := "ws://localhost:1234/"
	ws, _ := websocket.Dial(url, "", origin)

	ch := make(chan bool)

	err := websocket.JSON.Send(ws, Message{
		Message: "hallo",
		Type:    "LOGIN",
	})
	if err != nil {
		fmt.Println(err)
	}
	go read(ws)

	<-ch

}

func read(ws *websocket.Conn) {
	defer func() {
		ws.Close()
		fmt.Println("connection closed")
	}()

	for {
		var message Message
		// err := websocket.JSON.Receive(ws, &message)
		err := websocket.JSON.Receive(ws, &message)
		if err == io.EOF {
			return
		}

		fmt.Println(message)
	}
}
