package main

import (
	"fmt"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

func main() {
	board := engine.CreateBoard()
	m, score := engine.BestMove(board, 5, engine.WhiteColor)
	fmt.Printf("Best move: %v → %v  (score %d cp)\n", m.From, m.To, score)
}
