package game_state

import (
	"fmt"

	"github.com/g0g05arui/chess-engine/utils"
)

type Move = struct {
	From Position
	To   Position
}

func IsLegal(piece Piece, m Move, board Board) bool {
	if !(m.From.Line >= 1 && m.From.Line <= 8 && m.To.Line >= 1 && m.To.Line <= 8 &&
		m.From.Column >= 1 && m.From.Column <= 8 && m.To.Column >= 1 && m.To.Column <= 8) {
		return false
	}

	to_piece, found := _Find_Piece_By_Pos(m.To, board)

	if found && piece.Color == to_piece.Color {
		return false
	}

	newBoard := BoardAfterMove(m, board)

	if IsKingInCheck(newBoard, PieceColor(int8(WhiteColor)+utils.BoolToInt8(!board.WhiteTurn))) {
		return false
	}

	return true
}

func _Find_Piece_By_Pos(pos Position, board Board) (Piece, bool) {
	if pos.Column < 1 || pos.Column > 8 || pos.Line < 1 || pos.Line > 8 {
		return Piece{}, false
	}

	piece := board.PiecesMatrix[pos.Line][pos.Column]

	if piece.pType == 0 || (piece.Pos.Column == 0 && piece.Pos.Line == 0) {
		return Piece{}, false
	}

	return piece, true
}

func IsKingInCheck(board Board, color PieceColor) bool {
	// Find the king
	var king Piece
	found := false
	for _, p := range board.PiecesSlice {
		if p.pType == King && p.Color == color {
			king = p
			found = true
			break
		}
	}
	if !found {
		return true
	}

	// Check if any pieces of different color have visibility to the king
	for _, piece := range board.PiecesSlice {
		if piece.Color == color {
			continue
		}
		positions := GenerateAllVisiblePositions(piece, board)

		// Check if any position matches the king's position
		for _, pos := range positions {
			if pos.Column == king.Pos.Column && pos.Line == king.Pos.Line {
				return true
			}
		}
	}

	return false
}

func GenerateAllVisiblePositions(piece Piece, board Board) []Position {
	var positions []Position

	getDirections := func() []struct{ dc, dl int8 } {
		switch piece.pType {
		case Rook:
			return []struct{ dc, dl int8 }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		case Bishop:
			return []struct{ dc, dl int8 }{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}}
		case Queen:
			return []struct{ dc, dl int8 }{
				{1, 0}, {-1, 0}, {0, 1}, {0, -1},
				{1, 1}, {-1, -1}, {-1, 1}, {1, -1},
			}
		default:
			return nil
		}
	}

	addIfValid := func(c int8, l int8, allowCaptureOnly bool) {
		pos := Position{Column: c, Line: l}
		target, found := _Find_Piece_By_Pos(pos, board)
		if !found {
			if !allowCaptureOnly {
				positions = append(positions, pos)
			}
		} else if target.Color != piece.Color {
			positions = append(positions, pos)
		}
	}

	switch piece.pType {
	case Pawn:
		var dir int8 = 1
		startLine := 2
		if piece.Color != WhiteColor {
			dir = -1
			startLine = 7
		}

		// Forward 1 square
		front := Position{Column: piece.Pos.Column, Line: piece.Pos.Line + int8(dir)}
		if _, found := _Find_Piece_By_Pos(front, board); !found {
			positions = append(positions, front)

			// Forward 2 squares from starting rank
			if piece.Pos.Line == int8(startLine) {
				twoFront := Position{Column: piece.Pos.Column, Line: piece.Pos.Line + int8(2*dir)}
				if _, blocked := _Find_Piece_By_Pos(twoFront, board); !blocked {
					positions = append(positions, twoFront)
				}
			}
		}

		// Captures
		addIfValid(piece.Pos.Column-1, piece.Pos.Line+int8(dir), true)
		addIfValid(piece.Pos.Column+1, piece.Pos.Line+int8(dir), true)

	case Knight:
		moves := []struct{ dc, dl int }{
			{1, 2}, {2, 1}, {2, -1}, {1, -2},
			{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
		}
		for _, m := range moves {
			addIfValid(piece.Pos.Column+int8(m.dc), piece.Pos.Line+int8(m.dl), false)
		}

	case King:
		moves := []struct{ dc, dl int }{
			{1, 0}, {-1, 0}, {0, 1}, {0, -1},
			{1, 1}, {-1, -1}, {-1, 1}, {1, -1},
		}
		for _, m := range moves {
			addIfValid(piece.Pos.Column+int8(m.dc), piece.Pos.Line+int8(m.dl), false)
		}

	case Rook, Bishop, Queen:
		for _, d := range getDirections() {
			for i := int8(1); i < 8; i++ {
				c := piece.Pos.Column + i*d.dc
				l := piece.Pos.Line + int8(i*d.dl)
				pos := Position{Column: c, Line: l}
				target, found := _Find_Piece_By_Pos(pos, board)
				if !found {
					positions = append(positions, pos)
				} else {
					if target.Color != piece.Color {
						positions = append(positions, pos)
					}
					break
				}
			}
		}
	}

	return positions
}

func GenerateAllLegalMoves(piece Piece, board Board) []Position {
	visiblePositions := GenerateAllVisiblePositions(piece, board)
	var legalMoves []Position

	for _, pos := range visiblePositions {
		if IsLegal(piece, Move{From: piece.Pos, To: pos}, board) {
			legalMoves = append(legalMoves, pos)
		}
	}

	return legalMoves
}

func BoardAfterMove(m Move, board Board) Board {
	var updatedPieces []Piece
	var updatedMatrix [9][9]Piece

	// Initialize matrix with empty pieces
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			updatedMatrix[i][j] = Piece{} // Empty piece
		}
	}

	for _, piece := range board.PiecesSlice {
		// Skip piece at destination (captured piece)
		if piece.Pos.Column == m.To.Column && piece.Pos.Line == m.To.Line {
			continue
		}

		// Update piece if it's the one being moved
		if piece.Pos.Column == m.From.Column && piece.Pos.Line == m.From.Line {
			piece.Pos = m.To
			piece.hasMoved = true
		}

		// Add piece to updated slice and matrix
		updatedPieces = append(updatedPieces, piece)
		updatedMatrix[piece.Pos.Line][piece.Pos.Column] = piece
	}

	return Board{
		PiecesSlice:  updatedPieces,
		PiecesMatrix: updatedMatrix,
		WhiteTurn:    !board.WhiteTurn,
	}
}
func Perft(board Board, depth int, color PieceColor) int {
	if depth == 0 {
		return 1
	}

	count := 0
	for _, piece := range board.PiecesSlice {
		if piece.Color != color {
			continue
		}
		moves := GenerateAllLegalMoves(piece, board)
		for _, to := range moves {
			move := Move{From: piece.Pos, To: to}
			newBoard := BoardAfterMove(move, board)
			fmt.Println(BoardToString(newBoard))
			nextColor := BlackColor
			if color == BlackColor {
				nextColor = WhiteColor
			}
			count += Perft(newBoard, depth-1, nextColor)
		}
	}
	return count
}
