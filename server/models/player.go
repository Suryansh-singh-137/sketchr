package models

type InGamePlayer struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    Score     int    `json:"score"`
    IsDrawing bool   `json:"isDrawing"`
}

func (p *InGamePlayer) AddScore(points int) {
    p.Score += points
}

type Room struct {
    ID           string    `json:"id"`
    InGamePlayer      []*InGamePlayer `json:"ingameplayers"`
    IsGameActive bool      `json:"isGameActive"`
}

func (r *Room) AddPlayer(p *InGamePlayer) {
    r.InGamePlayer = append(r.InGamePlayer, p)
}