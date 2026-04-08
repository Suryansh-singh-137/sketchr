package handlers

import (
    "fmt"
    "net/http"
    "time"

    "github.com/Suryansh-singh-137/sketchr-server/models"
    "github.com/Suryansh-singh-137/sketchr-server/room"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func ServeWs(w http.ResponseWriter, r *http.Request, rm *room.RoomManager) {
    code := r.URL.Query().Get("room")
    username := r.URL.Query().Get("username")

    if code == "" {
        http.Error(w, "room code required", http.StatusBadRequest)
        return
    }

    if username == "" {
        username = "Anonymous"
    }

    foundRoom, exists := rm.GetRoom(code)
    if !exists {
        http.Error(w, "room not found", http.StatusNotFound)
        return
    }

    if len(foundRoom.Clients) >= foundRoom.MaxPlayers {
        http.Error(w, "room is full", http.StatusForbidden)
        return
    }

    connection, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "upgrade failed", http.StatusInternalServerError)
        return
    }

    c := models.Client{
        ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
        Username: username,
        Conn:     connection,
        Send:     make(chan []byte, 256),
    }

    foundRoom.Register <- &c
    foundRoom.Players = append(foundRoom.Players, username)

    go c.ReadPump(foundRoom.Incoming, foundRoom.Unregister)
    go c.WritePump()
}
