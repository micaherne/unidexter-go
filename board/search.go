package board

type BestMove struct {
	Move Move
}

func negamaxInternal(b *Board, depth int, bestMove *BestMove) int {
	if depth == 0 {
		return Evaluate(b)
	}
	max := -10000000
	moves := GenerateMoves(b)
	for _, move := range moves {
		if !LegalMove(b, move) {
			continue
		}

		MakeMove(b, move)
		score := -negamaxInternal(b, depth-1, bestMove)
		UndoMove(b)

		if score > max {
			max = score
		}

	}

	return max
}

func Negamax(b *Board, depth int) Move {
	bestMove := &BestMove{}
	max := -10000000
	moves := GenerateMoves(b)
	for _, move := range moves {
		if !LegalMove(b, move) {
			continue
		}

		MakeMove(b, move)
		score := -negamaxInternal(b, depth-1, bestMove)
		UndoMove(b)

		if score > max {
			max = score
			bestMove.Move = move
		}

	}

	if bestMove.Move.to == 0 && bestMove.Move.from == 0 {
		// TODO: Better last ditch attempt choice.
		// This just returns the one with the best static evaluation.

		for _, move := range moves {
			if !LegalMove(b, move) {
				continue
			}

			MakeMove(b, move)
			score := -negamaxInternal(b, 0, bestMove)
			UndoMove(b)

			if score > max {
				max = score
				bestMove.Move = move
			}

		}
	}
	return bestMove.Move
}
