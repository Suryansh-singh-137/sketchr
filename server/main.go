package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Suryansh-singh-137/sketchr-server/handlers"

	"github.com/Suryansh-singh-137/sketchr-server/room"
)
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        // handle preflight request
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next(w, r)
    }
}
func main() {
rm := room.NewRoomManager()
http.HandleFunc("/room/create", enableCORS(func(w http.ResponseWriter, r *http.Request) {
      host := r.URL.Query().Get("username")
      if host == "" {
        host = "Anonymous"
    }
    room := rm.CreateRoom(host)
    go room.Run()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"roomId": room.ID})
}))

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handlers.ServeWs(w, r, rm)
    })
    

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
