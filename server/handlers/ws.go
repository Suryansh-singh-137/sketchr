package handlers

import (
    "net/http"
    "github.com/Suryansh-singh-137/sketchr-server/hub"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func ServeWs(w http.ResponseWriter, r *http.Request, h *hub.Hub) {
    connection, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "not able to upgrade the conn", http.StatusInternalServerError)
        return
    }

    c := hub.Client{
        ID:   "2323r",
        Conn: connection,
        Send: make(chan []byte, 256),
    }

    h.Register <- &c
		go c.ReadPump(h)
go c.WritePump()
}