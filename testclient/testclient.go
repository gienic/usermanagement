package main

import (
	"fmt"
	"io"

	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"golang.org/x/net/websocket"
)

type Message struct {
	Type string
	AuthData
	EncryptData
	Error string
}

type EncryptData struct {
	Public []byte
}

type AuthData struct {
	Token string
}

var public rsa.PublicKey

func main() {

	origin := "http://localhost/"
	url := "ws://localhost:1234/"
	ws, _ := websocket.Dial(url, "", origin)

	ch := make(chan bool)

	go read(ws)

	err := websocket.JSON.Send(ws, Message{
		Type: "AUTH",
		AuthData: AuthData{
			Token: "12345678",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	pk, _ := json.Marshal(privKey.PublicKey)

	err = websocket.JSON.Send(ws, Message{
		Type: "ENCRYPT",
		EncryptData: EncryptData{
			Public: pk,
		},
	})

	if err != nil {
		fmt.Println(err)
	}

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

		json.Unmarshal(message.Public, public)
	}
}
