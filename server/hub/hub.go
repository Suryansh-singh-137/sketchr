package hub

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte  
}
type Hub  struct {
Clients    map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
}
//  initalises all the channel  of the  hub 
func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}
//  the  hub fxn 
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
					h.Clients[client]=true
           
        case client := <-h.Unregister:
           delete(h.Clients,client)
        case message := <-h.Broadcast:
				for client := range h.Clients {
    client.Send <- message
}
        }
    }
}
func (c *Client) ReadPump(broadcast chan []byte, unregister chan *Client) {
   
    defer func  (){
      unregister<-  c
      c.Conn.Close()
    }()
    for { 
        _, message ,err :=  c.Conn.ReadMessage()
        if err != nil {
            break ;
        }
        broadcast<-  message
    }
}

func (c *Client) WritePump() {
  for message := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
}