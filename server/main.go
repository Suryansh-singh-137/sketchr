package main

import (
	"fmt"

	"github.com/Suryansh-singh-137/sketchr-server/models"
)

func main() {
    // Create two players
    p1 := &models.Player{ID: "1", Username: "Raj", Score: 0, IsDrawing: true}
    p2 := &models.Player{ID: "2", Username: "Priya", Score: 0, IsDrawing: false}
    r1 := &models.Room{ID: "room-1", IsGameActive: false}
    r1.AddPlayer(p1)
    r1.AddPlayer(p2)

    p1.AddScore(10)

    fmt.Printf("Room: %+v\n", r1)
    fmt.Printf("Player 1: %+v\n", p1)
    fmt.Printf("Player 2: %+v\n", p2)
}
