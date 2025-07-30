package game_state

import "strings"

type Position = struct {
	Line   int8
	Column int8
}

type PieceColor int8
type PieceType int8

const (
	WhiteColor PieceColor = iota
	BlackColor
)

const (
	Pawn PieceType = iota + 1
	Bishop
	Rook
	Knight
	King
	Queen
)

type Piece = struct {
	Pos      Position
	Color    PieceColor
	pType    PieceType
	hasMoved bool
}

func _uncolored_PieceToString(piece Piece) string {
	switch piece.pType {
	case Pawn:
		return "P"
	case King:
		return "K"
	case Knight:
		return "N"
	case Queen:
		return "Q"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	default:
		return "."
	}
}

func PieceToString(piece Piece) string {
	str := _uncolored_PieceToString(piece)
	if piece.Color == BlackColor {
		str = strings.ToLower(str)
	}
	return str
}
