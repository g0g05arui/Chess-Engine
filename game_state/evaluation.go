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
	const ATTACK_VISIBILITY_MULTIPLIER int = 15 // percent-based scaling

	for _, p := range board.PiecesSlice {
		v := pieceValue[p.pType]

		// Central control bonus for pawns and knights
		if p.Pos.Line >= 4 && p.Pos.Line <= 5 && p.Pos.Column >= 3 && p.Pos.Column <= 6 &&
			(p.pType == Pawn || p.pType == Knight) {
			score += int(CENTER_VALUE_MULTIPLIER*float32(v)) * getSign(p.Color)
		}

		// Base material score
		score += v * getSign(p.Color)

		// Attack value bonus for Rook, Bishop, Knight
		if p.pType == Rook || p.pType == Bishop || p.pType == Knight {
			visible := GenerateAllVisiblePositions(p, board)
			for _, pos := range visible {
				target, found := _Find_Piece_By_Pos(pos, board)
				if found && target.Color != p.Color {
					targetValue := pieceValue[target.pType]
					bonus := (targetValue * ATTACK_VISIBILITY_MULTIPLIER) / 100
					score += bonus * getSign(p.Color)
				}
			}
		}
	}

	if sideToMove == BlackColor {
		score = -score
	}
	return score
}

func getSign(color PieceColor) int {
	if color == WhiteColor {
		return 1
	}
	return -1
}
