package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds incomming messages
	// that should be forwarded to the other clients
	forward chan []byte
	// join is a channel for clients wishing to join a room
	join chan *client
	// leave is a channel for clients wishing to leave a room
	leave chan *client
	// clients holds all current clients in this room
	clients map[*client]bool
}

// newRoom makes a room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	// run until terminated
	for {
		select {
		// message received on join channel
		case client := <-r.join:
			// joining
			r.clients[client] = true
		// message received on leave channel
		case client := <-r.leave:
			// leaving
			// delete the client from the map
			delete(r.clients, client)
			// close it's send channel
			close(client.send)
		// message received on forward channel
		case msg := <-r.forward:
			// forward message to all clients
			// iterate over all clients
			// add message to each clients send channel
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// cretae the websocket by calling upgrade methid
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	// create the client
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	// pass client to the join channel
	r.join <- client
	// defer the leaving operation for when client finished
	defer func() {
		r.leave <- client
	}()
	// run as a goroutine
	go client.write()
	// read method blocking operations
	// keeping the connecion alive
	client.read()
}
