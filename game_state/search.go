package game_state

import (
	"math"
	"sort"
	"sync"
)

const maxWorkers = 11

func BestMove(board Board, depth int, color PieceColor) (best Move, bestScore int) {
	type result struct {
		move  Move
		score int
	}

	jobs := make(chan Move, 100)
	results := make(chan result, 100)
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
		for _, to := range orderedMovesByEval(color, board, piece) {
			jobs <- to
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
		for _, mv := range orderedMovesByEval(color, board, piece) {
			child := BoardAfterMove(mv, board)
			score := -alphaBeta(child, depth-1, -beta, -alpha, opposite(color))
			if score > alpha {
				alpha = score
				if alpha >= beta {
					return alpha
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

func orderedMovesByEval(color PieceColor, board Board, piece Piece) []Move {
	moves := GenerateAllLegalMoves(piece, board)
	scored := make([]struct {
		move  Move
		score int
	}, len(moves))

	for i, to := range moves {
		mv := Move{From: piece.Pos, To: to}
		child := BoardAfterMove(mv, board)
		eval := Evaluate(child, color)
		scored[i] = struct {
			move  Move
			score int
		}{mv, eval}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	ordered := make([]Move, len(moves))
	for i, s := range scored {
		ordered[i] = s.move
	}
	return ordered
}
