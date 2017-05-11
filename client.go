package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) openReader() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

func (c *Client) openWriter() {
	defer func() { 
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <- c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<- c.send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
