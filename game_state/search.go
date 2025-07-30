// -------- search.go ----------
package game_state

import "math"

func BestMove(board Board, depth int, color PieceColor) (best Move, bestScore int) {
	bestScore = math.MinInt32
	for _, piece := range board.PiecesSlice {
		if piece.Color != color {
			continue
		}
		for _, to := range GenerateAllLegalMoves(piece, board) {
			mv := Move{From: piece.Pos, To: to}
			child := BoardAfterMove(mv, board)
			score := -alphaBeta(child, depth-1, -math.MaxInt32, math.MaxInt32,
				opposite(color))
			if score > bestScore {
				bestScore = score
				best = mv
			}
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
