package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Suryansh-singh-137/sketchr-server/handlers"
	"github.com/Suryansh-singh-137/sketchr-server/hub"
)

func main() {
    // Create two players
    // p1 := &models.Player{ID: "1", Username: "Raj", Score: 0, IsDrawing: true}
    // p2 := &models.Player{ID: "2", Username: "Priya", Score: 0, IsDrawing: false}
    // r1 := &models.Room{ID: "room-1", IsGameActive: false}
    // r1.AddPlayer(p1)
    // r1.AddPlayer(p2)

    // p1.AddScore(10)

    // fmt.Printf("Room: %+v\n", r1)
    // fmt.Printf("Player 1: %+v\n", p1)
    // fmt.Printf("Player 2: %+v\n", p2)
		   h := hub.NewHub()
    go h.Run()
      http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handlers.ServeWs(w, r, h)
    })

    fmt.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
