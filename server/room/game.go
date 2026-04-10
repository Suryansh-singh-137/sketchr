package room 
import (
    "encoding/json"
    "fmt"
    "math/rand"
    "strings"
    "time"
)
var words = []string{
    "elephant", "guitar", "pizza", "mountain", "bicycle",
    "umbrella", "dolphin", "rainbow", "telescope", "volcano",
    "helicopter", "butterfly", "waterfall", "diamond", "jungle",
    "penguin", "lighthouse", "saxophone", "tornado", "astronaut",
}
func getRandomeWord() string{
	return  words[rand.Intn(len(words))]
}
func (r *Room) StartNextTurn() {
    // 1. Agar sab players ki turn ho gayi → next round
    if r.TurnIndex > 0 && r.TurnIndex % len(r.Players) == 0 {
        r.CurrentRound++
    }

    // 2. Agar rounds complete → game over
    if r.CurrentRound > r.Rounds {
        r.EndGame()
        return
    }

    // 3. Next drawer choose karo
    drawer := r.Players[r.TurnIndex % len(r.Players)]
    r.CurrentDrawer = drawer
    r.CurrentWord = getRandomWord()
    r.TurnIndex++

    // 4. Sabko turn_start broadcast karo
    r.broadcastMessage(map[string]interface{}{
        "type":   "turn_start",
        "drawer": drawer,
        "round":  r.CurrentRound,
    })

    // 5. Sirf drawer ko word bhejo
    r.sendToPlayer(drawer, map[string]interface{}{
        "type": "your_word",
        "word": r.CurrentWord,
    })

    // 6. Baaki ko hint bhejo
    hint := makeHint(r.CurrentWord)
    r.broadcastMessage(map[string]interface{}{
        "type": "word_hint",
        "hint": hint,
    })

    // 7. Canvas clear karo sabka
    r.broadcastMessage(map[string]string{
        "type": "clear",
    })

    // 8. Timer shuru karo
    r.startTimer()
}
func makeHint(word string) string {
    hint := ""
    for _, ch := range word {
        if ch == ' ' {
            hint += " "
        } else {
            hint += "_ "
        }
    }
    return hint
}

func (r *Room) startTimer() {
    // purana timer cancel karo
    if r.GameTimer != nil {
        r.GameTimer.Stop()
    }

    // DrawTime seconds baad EndTurn call karo
    r.GameTimer = time.AfterFunc(
        time.Duration(r.DrawTime)*time.Second,
        func() {
            r.EndTurn()
        },
    )

    // Timer broadcast karo frontend ko
    go func() {
        for i := r.DrawTime; i >= 0; i-- {
            r.broadcastMessage(map[string]interface{}{
                "type":    "timer",
                "seconds": i,
            })
            time.Sleep(1 * time.Second)
        }
    }()
}

func (r *Room) EndTurn() {
    // Timer cancel karo
    if r.GameTimer != nil {
        r.GameTimer.Stop()
    }

    // Sabko word batao
    r.broadcastMessage(map[string]interface{}{
        "type": "turn_end",
        "word": r.CurrentWord,
    })

    // 3 seconds baad next turn
    time.Sleep(3 * time.Second)
    r.StartNextTurn()
}

func (r *Room) EndGame() {
    r.IsGameActive = false
    r.broadcastMessage(map[string]interface{}{
        "type":   "game_over",
        "scores": r.Scores,
    })
}

// points calculate karo — jaldi guess = zyada points
func (r *Room) calculatePoints(timeLeft int) int {
    return (timeLeft * 300) / r.DrawTime
}

// guess check karna — run.go ke chat case mein use hoga
func (r *Room) handleChat(message []byte) {
    var msg map[string]interface{}
    json.Unmarshal(message, &msg)

    username := fmt.Sprintf("%v", msg["username"])
    guess := fmt.Sprintf("%v", msg["message"])

    // Agar drawer hai → sirf broadcast karo, guess check mat karo
    if username == r.CurrentDrawer {
        r.broadcastToAll(message)
        return
    }

    // Sahi guess check karo
    if r.IsGameActive && strings.ToLower(guess) == strings.ToLower(r.CurrentWord) {
        // Points do
        r.Scores[username] += 100

        // Sabko batao
        r.broadcastMessage(map[string]interface{}{
            "type":     "correct_guess",
            "username": username,
        })
    } else {
        // Normal chat broadcast
        r.broadcastToAll(message)
    }
}
func  (r *Room)StartGame(){
	r.IsGameActive =  true 
	r.CurrentWord =  "1";
	r.TurnIndex  = 0
	for _,p := range r.Players{
		r.Scores[p] = 0
	}
	
    r.StartNextTurn()

}
