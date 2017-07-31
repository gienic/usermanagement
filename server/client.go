package server

import (
	"crypto/rsa"
	"fmt"
	"io"
	"time"

	"encoding/json"
	"golang.org/x/net/websocket"
)

const (
	// AUTH : message type for authorization
	AUTH string = "AUTH"
	// AUTH : message type for authorization
	ENCRYPT string = "ENCRYPT"
	// LOGIN : message type for login
	LOGIN string = "LOGIN"
	// SIGNUP : message type for signup
	SIGNUP string = "SIGNUP"
)

// TODO: remove receive channel
type (
	// Client ...
	Client struct {
		ID         int
		Conn       *websocket.Conn
		Send       chan Message
		Receive    chan Message
		ConnTime   time.Time
		Access     map[int]map[string]func(Message)
		Private    *rsa.PrivateKey
		Public     *rsa.PublicKey
		AuthState  int
		AuthCH     chan AuthMessage
		Unregister chan int
	}

	// Message ..
	Message struct {
		Type string
		AuthData
		EncryptData
		SignupData
		LoginData
		Error string
	}
)

// Login ...
func (c *Client) Login(msg Message) {

}

// SignUp ...
func (c *Client) SignUp(msg Message) {

}

func (c *Client) send() {

	defer func() {
		c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := websocket.JSON.Send(c.Conn, msg); err != nil {
			fmt.Println("cant send")
		}
	}
}

func (c *Client) receive() {

	defer func() {
		c.Conn.Close()
		fmt.Println("connection closed from client", c.ID)
		c.Unregister <- c.ID
	}()

	for {
		var msg Message
		if err := websocket.JSON.Receive(c.Conn, &msg); err != nil {
			if err == io.EOF {
				return
			}
		}

		c.handleMessage(msg)
	}
}

func (c *Client) handleMessage(msg Message) {
	if f, ok := c.Access[c.AuthState][msg.Type]; ok {
		f(msg)
	} else {
		c.Send <- Message{
			Error: "not authorized",
		}

		fmt.Println(
			"-->",
			"[Client:", c.ID, "|",
			"AuthState:", c.AuthState, "|",
			"Requested:", msg.Type, "|",
			"Response:", "not authorized]",
		)
	}

	fmt.Println("Authorization Status: ", c.AuthState)
}

// NewClient ...
func NewClient(ws *websocket.Conn, ac chan AuthMessage, uc chan int, id int) *Client {
	client := &Client{
		ID:         id,
		Conn:       ws,
		Send:       make(chan Message),
		Receive:    make(chan Message),
		AuthState:  UNUAUTHORIZED,
		ConnTime:   time.Now(),
		AuthCH:     ac,
		Unregister: uc,
	}

	client.Access = map[int]map[string]func(Message){
		UNUAUTHORIZED: {
			AUTH: func(msg Message) {
				client.AuthCH <- AuthMessage{
					ClientID: client.ID,
					Token:    msg.AuthData.Token,
				}
			},
		},
		PROCESSING: {
			ENCRYPT: func(msg Message) {
				json.Unmarshal(msg.Public, &client.Public)
				client.AuthState = ENCRYPTED
			},
		},
		ENCRYPTED: {
			LOGIN:  client.Login,
			SIGNUP: client.SignUp,
		},
	}

	return client
}
