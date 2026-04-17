package room

import (
	"math/rand"
	"sync"
	"time"

	"github.com/Suryansh-singh-137/sketchr-server/models"
)

type Room struct {
    ID           string                      `json:"id"`
    Clients      map[*models.Client]bool     `json:"clients"`
    Incoming     chan []byte                
    Register     chan *models.Client
    Unregister   chan *models.Client
    MaxPlayers   int                         `json:"maxPlayers"`
    IsGameActive bool                        `json:"isGameActive"`
    DrawTime     int                         `json:"drawTime"`
    Rounds       int                         `json:"rounds"`
    Host         string                      `json:"host"`
    Players      []string                    `json:"players"`
    // Game state
    CurrentWord   string
    CurrentDrawer string
    CurrentRound  int
    Scores        map[string]int
    TurnIndex     int
     GameTimer     *time.Timer 
    
}

func generateRoomCode() string {
    chars := []string{"aus", "tra", "vis", "head", "137", "travis", "gng", "cwc"}
    code := ""
    for i := 0; i <= 1; i++ {
        code += chars[rand.Intn(len(chars))]
    }
    return code
}

func NewRoom(host string) *Room {
    return &Room{
        ID:         generateRoomCode(),
        Clients:    make(map[*models.Client]bool),
        Incoming:   make(chan []byte),
        Register:   make(chan *models.Client),
        Unregister: make(chan *models.Client),
        MaxPlayers: 8,
        IsGameActive: false,
        DrawTime:   80,
        Rounds:     3,
        Host:       host,
        Players:    []string{},
        Scores:     make(map[string]int),
        TurnIndex:  0,
    }
}

type RoomManager struct {
    Rooms map[string]*Room
    mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
    return &RoomManager{
        Rooms: make(map[string]*Room),
    }
}

func (rm *RoomManager) CreateRoom(host string) *Room {
    room := NewRoom(host)
    rm.mu.Lock()
    defer rm.mu.Unlock()
    rm.Rooms[room.ID] = room
    return room
}

func (rm *RoomManager) GetRoom(id string) (*Room, bool) {
    rm.mu.Lock()
    defer rm.mu.Unlock()
    room, exists := rm.Rooms[id]
    return room, exists
}