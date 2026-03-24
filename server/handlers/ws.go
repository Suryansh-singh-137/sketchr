package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Suryansh-singh-137/sketchr-server/hub"
	"github.com/Suryansh-singh-137/sketchr-server/room"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func ServeWs(w http.ResponseWriter, r *http.Request,  rm *room.RoomManager) {
    code:= r.URL.Query().Get("room");
    if code == ""{
  http.Error(w, "code cannot be empty", http.StatusInternalServerError)
  return 
    }
  foundRoom  ,  exists := rm.GetRoom(code)
  if !exists {
    http.Error(w, "room not found", http.StatusNotFound)
    return 
  }
  if len(foundRoom.Clients)>= foundRoom.MaxPlayers { 
     http.Error(w, "room is full", http.StatusForbidden)
    return
  }
  

    connection, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "not able to upgrade the conn", http.StatusInternalServerError)
        return
    }
    
c := hub.Client{
    ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
    Conn: connection,
    Send: make(chan []byte, 236),
}
foundRoom.Register <- &c
 go c.ReadPump(foundRoom.Broadcast,foundRoom.Unregister)
go c.WritePump()  
}