package game_state

var pieceValue = map[PieceType]int{
	Pawn:   100,
	Knight: 320,
	Bishop: 330,
	Rook:   500,
	Queen:  900,
	King:   0,
}

func Evaluate(board Board, sideToMove PieceColor) int {
	score := 0
	const CENTER_VALUE_MULTIPLIER float32 = 0.05
	for _, p := range board.PiecesSlice {
		v := pieceValue[p.pType]
		pos := p.Pos
		if pos.Line >= 4 && pos.Line <= 5 && pos.Column >= 3 && pos.Column <= 6 && (p.pType == Pawn || p.pType == Knight) {
			score += int(CENTER_VALUE_MULTIPLIER * float32(pieceValue[p.pType]))
			if p.Color == WhiteColor {
				score += v
			} else {
				score -= v
			}
		}
		if p.Color == WhiteColor {
			score += v
		} else {
			score -= v
		}
	}

	if sideToMove == BlackColor {
		score = -score
	}
	return score
}
