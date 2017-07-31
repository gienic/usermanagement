package server

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

const (
	// UNUAUTHORIZED ...
	UNUAUTHORIZED int = iota
	// PROCESSING ...
	PROCESSING
	// ENCRYPTED ...
	ENCRYPTED
)

type (
	// Server ...
	Disposer struct {
		Clients    map[int]*Client
		AuthCH     chan AuthMessage
		unregister chan int
		Token      string
	}

	// AuthMessage ...
	AuthMessage struct {
		ClientID int
		Token    string
	}
)

func (d *Disposer) register(c *Client) {
	id := len(d.Clients)
	d.Clients[id] = c
}

func (d *Disposer) Unregister() {
	for {
		select {
		case id, ok := <-d.unregister:
			if ok {
				delete(d.Clients, id)
				fmt.Println("unregistered")
			}
		}
	}
}

// Authorize ...
func (d *Disposer) Authorize() {
	for {
		select {
		case authMSG, ok := <-d.AuthCH:
			if !ok {
				if authMSG.Token == d.Token {
					d.initializeE2E(authMSG.ClientID)
				}
			}
		}
	}
}

func (d *Disposer) initializeE2E(clientid int) {
	d.Clients[clientid].AuthState = PROCESSING

	key := d.generateKey()
	d.Clients[clientid].Private = key

	msg, _ := json.Marshal(key.Public())

	d.Clients[clientid].Send <- Message{
		EncryptData: EncryptData{
			Public: string(msg),
		},
	}

}

func (d *Disposer) generateKey() *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return privKey
	}

	return &rsa.PrivateKey{}
}

// Disposer ...
func (d *Disposer) Serve(ws *websocket.Conn) {
	client := NewClient(ws, d.AuthCH, d.unregister, len(d.Clients))

	d.register(client)

	go client.send()
	client.receive()
}

// NewDisposer ...
func NewDisposer(token string) *Disposer {
	return &Disposer{
		Clients:    make(map[int]*Client, 0),
		Token:      token,
		AuthCH:     make(chan AuthMessage, 0),
		unregister: make(chan int, 0),
	}
}
