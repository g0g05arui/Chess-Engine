package game_state

import "strings"

type Board = struct {
	PiecesSlice  []Piece
	PiecesMatrix [9][9]Piece
	WhiteTurn    bool
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
			Type:     Pawn,
			Color:    WhiteColor,
			Pos:      Position{Line: 2, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, pawn)
		matrix[2][i] = pawn

		// White back rank
		backPiece := Piece{
			Type:     backRank[i-1],
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
			Type:     Pawn,
			Color:    BlackColor,
			Pos:      Position{Line: 7, Column: i},
			hasMoved: false,
		}
		pieces = append(pieces, pawn)
		matrix[7][i] = pawn

		// Black back rank
		backPiece := Piece{
			Type:     backRank[i-1],
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
		WhiteTurn:    true,
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

func BoardToFEN(board Board) string {
	var sb strings.Builder

	for rank := int8(8); rank >= 1; rank-- {
		emptyCount := 0
		for file := int8(1); file <= 8; file++ {
			p := board.PiecesMatrix[rank][file]
			if p.Type == 0 {
				emptyCount++
			} else {
				if emptyCount > 0 {
					sb.WriteString(string('0' + emptyCount))
					emptyCount = 0
				}
				sb.WriteString(PieceToFENChar(p))
			}
		}
		if emptyCount > 0 {
			sb.WriteString(string('0' + emptyCount))
		}
		if rank > 1 {
			sb.WriteString("/")
		}
	}

	// Active color
	if board.WhiteTurn {
		sb.WriteString(" w")
	} else {
		sb.WriteString(" b")
	}

	// Default values for castling, en passant, halfmove clock, and fullmove number
	sb.WriteString(" - - 0 1")

	return sb.String()
}

func PieceToFENChar(p Piece) string {
	var ch byte
	switch p.Type {
	case Pawn:
		ch = 'p'
	case Knight:
		ch = 'n'
	case Bishop:
		ch = 'b'
	case Rook:
		ch = 'r'
	case Queen:
		ch = 'q'
	case King:
		ch = 'k'
	default:
		return ""
	}

	if p.Color == WhiteColor {
		ch -= 32 // convert to uppercase
	}

	return string(ch)
}

func FENToBoard(fen string) Board {
	fields := strings.Fields(fen)
	if len(fields) < 2 {
		panic("invalid FEN string")
	}

	board := Board{
		PiecesSlice:  []Piece{},
		PiecesMatrix: [9][9]Piece{},
		WhiteTurn:    fields[1] == "w",
	}

	ranks := strings.Split(fields[0], "/")
	if len(ranks) != 8 {
		panic("invalid piece placement in FEN")
	}

	for rank := 8; rank >= 1; rank-- {
		line := ranks[8-rank]
		file := int8(1)
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch >= '1' && ch <= '8' {
				file += int8(ch - '0')
				continue
			}

			var color PieceColor
			var pType PieceType

			if ch >= 'A' && ch <= 'Z' {
				color = WhiteColor
				ch += 32 // convert to lowercase
			} else {
				color = BlackColor
			}

			switch ch {
			case 'p':
				pType = Pawn
			case 'n':
				pType = Knight
			case 'b':
				pType = Bishop
			case 'r':
				pType = Rook
			case 'q':
				pType = Queen
			case 'k':
				pType = King
			default:
				panic("invalid piece character in FEN")
			}

			p := Piece{
				Type:     pType,
				Color:    color,
				Pos:      Position{Line: int8(rank), Column: file},
				hasMoved: true,
			}

			board.PiecesSlice = append(board.PiecesSlice, p)
			board.PiecesMatrix[rank][file] = p
			file++
		}
	}

	return board
}
