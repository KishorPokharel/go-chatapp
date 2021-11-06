package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/KishorPokharel/chatapp/pkg/forms"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type message struct {
	EventType string
	Name      string
	Message   string
	CreatedAt time.Time
}

type client struct {
	socket *websocket.Conn
	send   chan message
	room   room
	user   map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			log.Println("[Client] could not read message", err)
			return
		}
		name, _ := c.user["name"].(string)
		msg.Name = name
		msg.CreatedAt = time.Now()
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		name, _ := c.user["name"].(string)
		if msg.EventType == "" {
			if msg.Name == name {
				msg.EventType = "messageSent"
			} else {
				msg.EventType = "messageReceived"
			}
		}
		err := c.socket.WriteJSON(msg)
		if err != nil {
			log.Println("[Client] could not write message", err)
			return
		}
	}
}

type room struct {
	forward chan message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func newRoom() *room {
	r := &room{
		forward: make(chan message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
	return r
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			name, _ := client.user["name"].(string)
			msg := message{
				EventType: "userJoin",
				Name:      name,
			}
			for c := range r.clients {
				c.send <- msg
			}
		case client := <-r.leave:
			name, _ := client.user["name"].(string)
			close(client.send)
			delete(r.clients, client)
			msg := message{
				EventType: "userLeave",
				Name:      name,
			}
			for c := range r.clients {
				c.send <- msg
			}
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (app *application) roomHandler(w http.ResponseWriter, r *http.Request) {
	usr := app.contextGetUser(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Println("Could not upgrade connection: ", err)
		return
	}

	client := &client{
		socket: conn,
		send:   make(chan message, messageBufferSize),
		room:   *app.chatroom,
		user: map[string]interface{}{
			"name": usr.Username,
		},
	}
	client.room.join <- client
	defer func() {
		client.room.leave <- client
	}()
	go client.write()
	client.read()
}

func (app *application) chatHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "chat.html", &templateData{Form: forms.New(nil)})
}
