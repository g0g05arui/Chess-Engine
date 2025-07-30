package game_state

import "strings"

type Board = struct {
	PiecesSlice  []Piece
	PiecesMatrix [9][9]Piece
	whiteTurn    bool
}

func CreateBoard() Board {
	var pieces []Piece
	var matrix [9][9]Piece

	// Initialize matrix with empty pieces
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			matrix[i][j] = Piece{} // Empty piece
		}
	}

	backRank := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}

	// White pieces
	for i := int8(1); i <= 8; i++ {
		// White pawns
		pawn := Piece{
			pType:    Pawn,
			Color:    WhiteColor,
			Pos:      Position{Line: 2, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, pawn)
		matrix[2][i] = pawn

		// White back rank
		backPiece := Piece{
			pType:    backRank[i-1],
			Color:    WhiteColor,
			Pos:      Position{Line: 1, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, backPiece)
		matrix[1][i] = backPiece
	}

	// Black pieces
	for i := int8(1); i <= 8; i++ {
		// Black pawns
		pawn := Piece{
			pType:    Pawn,
			Color:    BlackColor,
			Pos:      Position{Line: 7, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, pawn)
		matrix[7][i] = pawn

		// Black back rank
		backPiece := Piece{
			pType:    backRank[i-1],
			Color:    BlackColor,
			Pos:      Position{Line: 8, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, backPiece)
		matrix[8][i] = backPiece
	}

	return Board{
		PiecesSlice:  pieces,
		PiecesMatrix: matrix,
		whiteTurn:    true,
	}
}

func BoardToString(board Board) string {
	grid := [9][9]string{}
	for i := int8(1); i <= 8; i++ {
		for j := int8(1); j <= 8; j++ {
			grid[i][j] = "."
		}
	}

	for _, p := range board.PiecesSlice {
		grid[p.Pos.Line][p.Pos.Column] = PieceToString(p)
	}

	var sb strings.Builder
	for rank := int8(8); rank >= 1; rank-- {
		for file := int8(1); file <= 8; file++ {
			sb.WriteString(grid[rank][file])
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
