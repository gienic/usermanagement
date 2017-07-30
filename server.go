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
	Server struct {
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

func (s *Server) register(c *Client) {
	id := len(s.Clients)
	s.Clients[id] = c
}

func (s *Server) Unregister() {
	for {
		select {
		case id, ok := <-s.unregister:
			if ok {
				delete(s.Clients, id)
				fmt.Println("unregistered")
			}
		}
	}
}

// Authorize ...
func (s *Server) Authorize() {
	for {
		select {
		case authMSG, ok := <-s.AuthCH:
			if !ok {
				if authMSG.Token == s.Token {
					s.initializeE2E(authMSG.ClientID)
				}
			}
		}
	}
}

func (s *Server) initializeE2E(clientid int) {
	s.Clients[clientid].AuthState = PROCESSING

	key := s.generateKey()
	s.Clients[clientid].Private = key

	msg, _ := json.Marshal(key.Public())

	s.Clients[clientid].Send <- Message{
		EncryptData: EncryptData{
			Public: string(msg),
		},
	}

}

func (s *Server) generateKey() *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return privKey
	}

	return &rsa.PrivateKey{}
}

// Serve ...
func (s *Server) Serve(ws *websocket.Conn) {
	client := NewClient(ws, &s.AuthCH, &s.unregister, len(s.Clients))

	s.register(client)

	go client.send()
	client.receive()
}

// NewServer ...
func NewServer(token string) *Server {
	return &Server{
		Clients:    make(map[int]*Client, 0),
		Token:      token,
		AuthCH:     make(chan AuthMessage, 0),
		unregister: make(chan int, 0),
	}
}
