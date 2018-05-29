package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client represents a single chatting user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on  which messages are sent
	send chan *message
	// room is the room this client is chatting in
	room *room

	// userData holds information about the user
	userData map[string]interface{}
}

// allows client to read from the socket via ReadMessage method
func (c *client) read() {
	// defer = close the socket when the function returns
	defer c.socket.Close()
	for {
		// _, msg, err := c.socket.ReadMessage()
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		// enqueue entry
		// https://stackoverflow.com/a/15983335/3907839
		c.room.forward <- msg
	}
}

// continually accepts messages from the send channel writing
// everything out of the socket via the WriteMessage method
func (c *client) write() {
	// defer = close the socket when the function returns
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
