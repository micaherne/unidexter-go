package board

import "fmt"

var pieceValues = map[int]int{
	PAWN:   100,
	KNIGHT: 300,
	BISHOP: 310,
	ROOK:   500,
	QUEEN:  875,
	KING:   0, // Doesn't contribute
}

// Evaluate returns a score for a position from the point of view of the side to move.
func Evaluate(b *Board) int {
	result := 0
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank<<4 | file
			if !LegalSquareIndex(square) {
				fmt.Printf("Illegal square %X\n", square)
				return result // TODO: Get rid of this test stuff
			}
			if b.squares[square] == EMPTY {
				continue
			}

			// fmt.Printf("Square %X, piece %d\n", square, GetPieceType(b.squares[square]))
			value := pieceValues[GetPieceType(b.squares[square])]
			if GetColour(b.squares[square]) == WHITE {
				result += value
			} else {
				result -= value
			}
			//fmt.Printf("Result: %d\n", result)
		}
	}
	// fmt.Printf("Final result: %d\n", result)
	if b.whiteToMove {
		return result
	}
	return -result

}
