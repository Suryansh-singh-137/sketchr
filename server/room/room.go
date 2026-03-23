package room

import (
	"math/rand"
	"sync"

	"github.com/Suryansh-singh-137/sketchr-server/hub"
)

type Room struct {
	ID           string     `json:"id"`
	Clients      map[*hub.Client]bool `json:"clients"`
	Broadcast    chan []byte `json:"broadcast"`
	Register     chan *hub.Client  `json:"register"`
	Unregister   chan *hub.Client `json:"unregister"`
	MaxPlayers   int `json:"maxplayer"` 
	IsGameActive bool `json:"isgameactive"` 
}
// genrate a uniue  room code of 6 charcter 

func generateRoomCode()string {
 chars := []string {"aus","tra","vis","head","137","travis","gng","cwc"}
 code :="" 

for i := 0; i <= 1; i++ {
    code += chars[rand.Intn(len(chars))]
}
 
 return  code
}
//  new room is getting created with fresh id  ,maxplayer  , and channel intialising for tht particular room 
func NewRoom() *Room{
 return   &Room{
	ID:  generateRoomCode(),
	Clients:  make(map[*hub.Client]bool),
	Broadcast:  make(chan []byte),
	Register:  make(chan *hub.Client),
	Unregister: make(chan *hub.Client),
	MaxPlayers:  8,
	IsGameActive:  false,
 }
} 
//  room  manager struct  code -> room struct
type RoomManager struct{
	Rooms map[string]*Room
	mu sync.Mutex
}
// newroommanager  
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms: make(map[string]*Room),
	}
}
// creatoing a new room  
func (rm *RoomManager) CreateRoom() *Room{
	
	room:=NewRoom()
	rm.mu.Lock()
	rm.Rooms[room.ID] =room
		defer rm.mu.Unlock()
return  room
}