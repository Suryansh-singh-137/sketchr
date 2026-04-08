package models

import (
    "github.com/gorilla/websocket"
)

type Client struct {
    ID       string
    Username string
    Conn     *websocket.Conn
    Send     chan []byte
}

func (c *Client) ReadPump(incoming chan []byte, unregister chan *Client) {
    defer func() {
        unregister <- c
        c.Conn.Close()
    }()
    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
        incoming <- message
    }
}

func (c *Client) WritePump() {
    defer c.Conn.Close()
    for message := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}