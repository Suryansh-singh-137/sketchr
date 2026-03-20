package models

type Player struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    Score     int    `json:"score"`
    IsDrawing bool   `json:"isDrawing"`
}

func (p *Player) AddScore(points int) {
    p.Score += points
}

type Room struct {
    ID           string    `json:"id"`
    Players      []*Player `json:"players"`
    IsGameActive bool      `json:"isGameActive"`
}

func (r *Room) AddPlayer(p *Player) {
    r.Players = append(r.Players, p)
}