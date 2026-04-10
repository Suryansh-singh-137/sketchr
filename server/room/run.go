package room

import (
    "encoding/json"
)

func (r *Room) Run() {
    for {
        select {
        case client := <-r.Register:
            r.Clients[client] = true

        case client := <-r.Unregister:
            delete(r.Clients, client)
            r.removePlayer(client.Username)

        case message := <-r.Incoming:
            r.routeMessage(message)
        }
    }
}

func (r *Room) routeMessage(message []byte) {
    var msg map[string]interface{}
    json.Unmarshal(message, &msg)

    msgType, _ := msg["type"].(string)

    switch msgType {
    case "draw", "clear":
        r.broadcastToAll(message)
    case "chat":
        r.broadcastToAll(message)
    case "game_start":
        go r.StartGame()  
        r.broadcastToAll(message)
    case "player_joined":
        r.broadcastToAll(message)
    }
}

func (r *Room) broadcastToAll(message []byte) {
    for client := range r.Clients {
        client.Send <- message
    }
}

func (r *Room) broadcastMessage(msg interface{}) {
    data, _ := json.Marshal(msg)
    r.broadcastToAll(data)
}

func (r *Room) sendToPlayer(username string, msg interface{}) {
    data, _ := json.Marshal(msg)
    for client := range r.Clients {
        if client.Username == username {
            client.Send <- data
            return
        }
    }
}

func (r *Room) removePlayer(username string) {
    for i, p := range r.Players {
        if p == username {
            r.Players = append(r.Players[:i], r.Players[i+1:]...)
            break
        }
    }
}