package main

import (
	"fmt"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

func main() {
	board := engine.CreateBoard()
	for i := 1; i <= 10; i++ {
		color := engine.WhiteColor
		if !board.WhiteTurn {
			color = engine.BlackColor
		}
		m, _ := engine.BestMove(board, 6, color)
		board = engine.BoardAfterMove(m, board)
		fmt.Println(engine.BoardToFEN(board))
	}
}
