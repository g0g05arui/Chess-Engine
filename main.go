package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

func fetchAndSaveImageAsync(fen string, moveNumber int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("https://fen2image.chessvision.ai/%s", fen)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Move %d: failed to fetch image: %v\n", moveNumber, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Move %d: bad response: %s\n", moveNumber, resp.Status)
		return
	}

	filePath := filepath.Join("game_status", fmt.Sprintf("%d.png", moveNumber))
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Move %d: failed to create file: %v\n", moveNumber, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Move %d: failed to save image: %v\n", moveNumber, err)
	}
}

func main() {
	// Ensure the output directory exists
	if err := os.MkdirAll("game_status", os.ModePerm); err != nil {
		fmt.Printf("Failed to create game_status directory: %v\n", err)
		return
	}

	board := engine.CreateBoard()
	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		color := engine.WhiteColor
		if !board.WhiteTurn {
			color = engine.BlackColor
		}

		m, _ := engine.BestMove(board, 5, color)
		board = engine.BoardAfterMove(m, board)

		fen := engine.BoardToFEN(board)
		fmt.Println(fen)

		wg.Add(1)
		go fetchAndSaveImageAsync(fen, i, &wg)
	}

	wg.Wait()
	fmt.Println("All images downloaded.")
}
