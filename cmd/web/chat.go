package main

import (
	"log"
	"net/http"
	"time"

	"github.com/KishorPokharel/chatapp/pkg/forms"
	"github.com/KishorPokharel/chatapp/pkg/models"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type eventType string

const (
	eventChat     eventType = "message"
	eventJoin     eventType = "join"
	eventLeave    eventType = "leave"
	eventError    eventType = "error"
	eventUserList eventType = "users"
)

type message struct {
	Type    eventType `json:"type"`
	Sender  string    `json:"sender,omitempty"`
	Content string    `json:"content,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Clients []string  `json:"clients,omitempty"`
}

type client struct {
	socket *websocket.Conn
	send   chan message
	room   room
	user   map[string]interface{}
	app    *application
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
		id, _ := c.user["id"].(int64)
		msg.Sender = name
		m := &models.Message{UserID: id, Body: msg.Content}
		err = c.app.models.Messages.Insert(m)
		if err != nil {
			c.app.logger.Println("Couldnot insert a message: ", err)
			continue
		}
		msg.Date = m.CreatedAt
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		//name, _ := c.user["name"].(string)
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
			currentClientNames := make(map[string]bool)
			for key := range r.clients {
				name, _ := key.user["name"].(string)
				currentClientNames[name] = true
			}
			keys := []string{}
			for key := range currentClientNames {
				keys = append(keys, key)
			}
			msg := message{
				Type:    eventJoin,
				Clients: keys,
			}
			for c := range r.clients {
				c.send <- msg
			}
		case client := <-r.leave:
			name, _ := client.user["name"].(string)
			delete(r.clients, client)
			close(client.send)
			msg := message{
				Type:   eventLeave,
				Sender: name,
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
			"id":   usr.ID,
		},
		app: app,
	}
	client.room.join <- client
	defer func() {
		client.room.leave <- client
	}()
	go client.write()
	client.read()
}

func (app *application) chatHandler(w http.ResponseWriter, r *http.Request) {
	usr := app.contextGetUser(r)
	messages, err := app.models.Messages.GetAll()
	if err != nil {
		app.serverError(w, r, err)
	}
	app.render(w, r, "chat.html", &templateData{
		Form:     forms.New(nil),
		Messages: messages,
		User:     usr,
	})
}
