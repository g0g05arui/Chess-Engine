// -------- search.go ----------
package game_state

import (
	"math"
	"sync"
)

const maxWorkers = 11

func BestMove(board Board, depth int, color PieceColor) (best Move, bestScore int) {
	type result struct {
		move  Move
		score int
	}

	jobs := make(chan Move, 100)      // Moves to evaluate
	results := make(chan result, 100) // Scored results
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for mv := range jobs {
				child := BoardAfterMove(mv, board)
				score := -alphaBeta(child, depth-1, -math.MaxInt32, math.MaxInt32, opposite(color))
				results <- result{mv, score}
			}
		}()
	}

	// Generate all moves and send them to the jobs channel
	for _, piece := range board.PiecesSlice {
		if piece.Color != color {
			continue
		}
		for _, to := range GenerateAllLegalMoves(piece, board) {
			mv := Move{From: piece.Pos, To: to}
			jobs <- mv
		}
	}
	close(jobs)

	// Wait for all workers to finish, then close the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	bestScore = math.MinInt32
	for res := range results {
		if res.score > bestScore {
			bestScore = res.score
			best = res.move
		}
	}
	return
}

func alphaBeta(board Board, depth int, alpha, beta int, color PieceColor) int {
	if depth == 0 {
		return Evaluate(board, color)
	}

	for _, piece := range board.PiecesSlice {
		if piece.Color != color {
			continue
		}
		for _, to := range GenerateAllLegalMoves(piece, board) {
			mv := Move{From: piece.Pos, To: to}
			child := BoardAfterMove(mv, board)
			score := -alphaBeta(child, depth-1, -beta, -alpha, opposite(color))
			if score > alpha {
				alpha = score
				if alpha >= beta {
					return alpha // β‑cutoff
				}
			}
		}
	}
	return alpha
}

func opposite(c PieceColor) PieceColor {
	if c == WhiteColor {
		return BlackColor
	}
	return WhiteColor
}
