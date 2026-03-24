package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Suryansh-singh-137/sketchr-server/handlers"
	
	"github.com/Suryansh-singh-137/sketchr-server/room"
)

func main() {
   rm:= room.NewRoomManager()
   r := rm.CreateRoom()
   fmt.Println("roomhaas been created with id:",r.ID)
   go r.Run()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handlers.ServeWs(w, r, rm)
    })

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
