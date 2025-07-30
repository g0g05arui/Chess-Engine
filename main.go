package main

import (
	"fmt"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

func main() {
	board := engine.CreateBoard()
	nodes := engine.Perft(board, 5, engine.WhiteColor)
	fmt.Println("Total positions at depth 4:", nodes)
}
