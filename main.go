package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/g0g05arui/chess-engine/computed"
	"github.com/g0g05arui/chess-engine/game_state"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	r := gin.Default()
	r.GET("/best-move", func(c *gin.Context) {
		fen := c.DefaultQuery("fen", "")
		turn := c.DefaultQuery("turn", "white")
		if fen == "" {
			c.Status(400)
			return
		}

		color := game_state.WhiteColor
		if turn != "white" {
			color = game_state.BlackColor
		}
		board := game_state.FENToBoard(fen)
		cacheKey := computed.CacheKey{Fen: fen, WhiteTurn: turn == "white"}

		const defaultDepth = 4

		// Start with default cached result
		move, ok := computed.Cache[cacheKey]
		currentDepth := defaultDepth
		if !ok {
			move, _ = game_state.BestMove(board, defaultDepth, color)
			computed.Cache[cacheKey] = move
		}

		// Search for the deepest available result in DeepCache
		for d := defaultDepth + 1; d <= defaultDepth+3; d++ { // Look ahead up to 3 levels
			deepKey := computed.DeepCacheKey{Fen: fen, WhiteTurn: turn == "white", Depth: d}
			if val, exists := computed.DeepCache[deepKey]; exists {
				move = val.BestMove
				currentDepth = val.Depth
			}
		}

		// Return the best available move
		c.JSON(200, gin.H{
			"best_move": move,
			"depth":     currentDepth,
		})

		// Launch next-depth search if not already present
		go func(fen string, turn string, board game_state.Board, currentDepth int) {
			color := game_state.WhiteColor
			if turn != "white" {
				color = game_state.BlackColor
			}
			nextDepth := currentDepth + 1
			deepKey := computed.DeepCacheKey{Fen: fen, WhiteTurn: turn == "white", Depth: nextDepth}

			// If already in DeepCache, don't compute
			if _, exists := computed.DeepCache[deepKey]; exists {
				return
			}

			// If already computing, don't start again
			if _, alreadyComputing := computed.InProgress.LoadOrStore(deepKey, struct{}{}); alreadyComputing {
				return
			}

			// Defer removing from in-progress after we're done
			defer computed.InProgress.Delete(deepKey)

			fmt.Printf("Computing deeper best move for depth %d...\n", nextDepth)
			move, _ := game_state.BestMove(board, nextDepth, color)
			computed.DeepCache[deepKey] = computed.DeepCacheValue{
				BestMove: move,
				Depth:    nextDepth,
			}
		}(fen, turn, board, currentDepth)

	})

	r.Run(":" + PORT)
}
