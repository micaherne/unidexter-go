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

// The piece square tables are just nicked from crafty for the time being.
var pieceSquareKnight = [2][2][64]int{
	{
		{-41, -29, -27, -15, -15, -27, -29, -41,
			-9, 4, 14, 20, 20, 14, 4, -9,
			-7, 10, 23, 29, 29, 23, 10, -7,
			-5, 12, 25, 32, 32, 25, 12, -5, /* [mg][black][sq] */
			-5, 10, 23, 28, 28, 23, 10, -5,
			-7, -2, 19, 19, 19, 19, -2, -7,
			-9, -6, -2, 0, 0, -2, -6, -9,
			-31, -29, -27, -25, -25, -27, -29, -31},

		{-31, -29, -27, -25, -25, -27, -29, -31,
			-9, -6, -2, 0, 0, -2, -6, -9,
			-7, -2, 19, 19, 19, 19, -2, -7,
			-5, 10, 23, 28, 28, 23, 10, -5, /* [mg][white][sq] */
			-5, 12, 25, 32, 32, 25, 12, -5,
			-7, 10, 23, 29, 29, 23, 10, -7,
			-9, 4, 14, 20, 20, 14, 4, -9,
			-41, -29, -27, -15, -15, -27, -29, -41}},

	{
		{-41, -29, -27, -15, -15, -27, -29, -41,
			-9, 4, 14, 20, 20, 14, 4, -9,
			-7, 10, 23, 29, 29, 23, 10, -7,
			-5, 12, 25, 32, 32, 25, 12, -5, /* [eg][black][sq] */
			-5, 10, 23, 28, 28, 23, 10, -5,
			-7, -2, 19, 19, 19, 19, -2, -7,
			-9, -6, -2, 0, 0, -2, -6, -9,
			-31, -29, -27, -25, -25, -27, -29, -31},

		{-31, -29, -27, -25, -25, -27, -29, -31,
			-9, -6, -2, 0, 0, -2, -6, -9,
			-7, -2, 19, 19, 19, 19, -2, -7,
			-5, 10, 23, 28, 28, 23, 10, -5, /* [eg][white][sq] */
			-5, 12, 25, 32, 32, 25, 12, -5,
			-7, 10, 23, 29, 29, 23, 10, -7,
			-9, 4, 14, 20, 20, 14, 4, -9,
			-41, -29, -27, -15, -15, -27, -29, -41},
	},
}

var pieceSquareBishop = [2][2][64]int{
	{{0, 0, 0, 0, 0, 0, 0, 0,
		0, 4, 4, 4, 4, 4, 4, 0,
		0, 4, 8, 8, 8, 8, 4, 0,
		0, 4, 8, 12, 12, 8, 4, 0,
		0, 4, 8, 12, 12, 8, 4, 0, /* [mg][black][sq] */
		0, 4, 8, 8, 8, 8, 4, 0,
		0, 4, 4, 4, 4, 4, 4, 0,
		-15, -15, -15, -15, -15, -15, -15, -15},

		{-15, -15, -15, -15, -15, -15, -15, -15,
			0, 4, 4, 4, 4, 4, 4, 0,
			0, 4, 8, 8, 8, 8, 4, 0,
			0, 4, 8, 12, 12, 8, 4, 0,
			0, 4, 8, 12, 12, 8, 4, 0, /* [mg][white][sq] */
			0, 4, 8, 8, 8, 8, 4, 0,
			0, 4, 4, 4, 4, 4, 4, 0,
			0, 0, 0, 0, 0, 0, 0, 0}},

	{{0, 0, 0, 0, 0, 0, 0, 0,
		0, 4, 4, 4, 4, 4, 4, 0,
		0, 4, 8, 8, 8, 8, 4, 0,
		0, 4, 8, 12, 12, 8, 4, 0,
		0, 4, 8, 12, 12, 8, 4, 0, /* [eg][black][sq] */
		0, 4, 8, 8, 8, 8, 4, 0,
		0, 4, 4, 4, 4, 4, 4, 0,
		-15, -15, -15, -15, -15, -15, -15, -15},

		{-15, -15, -15, -15, -15, -15, -15, -15,
			0, 4, 4, 4, 4, 4, 4, 0,
			0, 4, 8, 8, 8, 8, 4, 0,
			0, 4, 8, 12, 12, 8, 4, 0,
			0, 4, 8, 12, 12, 8, 4, 0, /* [eg][white][sq] */
			0, 4, 8, 8, 8, 8, 4, 0,
			0, 4, 4, 4, 4, 4, 4, 0,
			0, 0, 0, 0, 0, 0, 0, 0},
	},
}

var pieceSquareQueen = [2][2][64]int{
	{{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 4, 4, 4, 4, 0, 0,
		0, 4, 4, 6, 6, 4, 4, 0,
		0, 4, 6, 8, 8, 6, 4, 0,
		0, 4, 6, 8, 8, 6, 4, 0, /* [mg][black][sq] */
		0, 4, 4, 6, 6, 4, 4, 0,
		0, 0, 4, 4, 4, 4, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0},

		{0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 4, 4, 4, 4, 0, 0,
			0, 4, 4, 6, 6, 4, 4, 0,
			0, 4, 6, 8, 8, 6, 4, 0,
			0, 4, 6, 8, 8, 6, 4, 0, /* [mg][white][sq] */
			0, 4, 4, 6, 6, 4, 4, 0,
			0, 0, 4, 4, 4, 4, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0},
	},

	{{0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 4, 4, 4, 4, 0, 0,
		0, 4, 4, 6, 6, 4, 4, 0,
		0, 4, 6, 8, 8, 6, 4, 0,
		0, 4, 6, 8, 8, 6, 4, 0, /* [eg][black][sq] */
		0, 4, 4, 6, 6, 4, 4, 0,
		0, 0, 4, 4, 4, 4, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0},

		{0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 4, 4, 4, 4, 0, 0,
			0, 4, 4, 6, 6, 4, 4, 0,
			0, 4, 6, 8, 8, 6, 4, 0,
			0, 4, 6, 8, 8, 6, 4, 0, /* [eg][white][sq] */
			0, 4, 4, 6, 6, 4, 4, 0,
			0, 0, 4, 4, 4, 4, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0},
	},
}

var pieceSquarePawn = [2][64]int{
	{0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 3, 5, 5, 3, 0, -5, /* [mg][black][sq] */
		-5, 0, 5, 10, 10, 5, 0, -5,
		-5, 0, 3, 5, 5, 3, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 0, 0, 0, 0, 0},

	{0, 0, 0, 0, 0, 0, 0, 0,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 3, 5, 5, 3, 0, -5,
		-5, 0, 5, 10, 10, 5, 0, -5, /* [mg][white][sq] */
		-5, 0, 3, 5, 5, 3, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 0, 0, 0, 0, 0},
}

// Evaluate returns a score for a position from the point of view of the side to move.
func Evaluate(b *Board) int {
	result := 0
	blackMaterial := 0
	whiteMaterial := 0
	sideToMove := BLACK
	if b.whiteToMove {
		sideToMove = WHITE
	}
	var kingPosition int
	var opponentColour int
	if sideToMove == WHITE {
		kingPosition = b.whiteKing
		opponentColour = BLACK
	} else { // Assume calling code is correct to avoid extra if
		kingPosition = b.blackKing
		opponentColour = WHITE
	}

	// Test for checkmate.
	if IsCheck(b, sideToMove) {
		escapeSquareFound := false
		for _, offset := range DIAGONALSANDLINES {
			if square := kingPosition + offset; LegalSquareIndex(square) && !IsAttacked(b, square, opponentColour) {
				escapeSquareFound = true
				break
			}
		}
		if !escapeSquareFound {
			return -30000
		}
	}

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

			pieceType := GetPieceType(b.squares[square])
			value := pieceValues[pieceType]
			colour := GetColour(b.squares[square])
			if colour == WHITE {
				whiteMaterial += value
			} else {
				blackMaterial += value
			}

			// Add piece square bonuses.
			pieceSquareIndex := rank*8 + file
			colourIndex := colour >> 3
			switch pieceType {
			case KNIGHT:
				result += pieceSquareKnight[0][colourIndex][pieceSquareIndex]
			case BISHOP:
				result += pieceSquareBishop[0][colourIndex][pieceSquareIndex]
			case QUEEN:
				result += pieceSquareQueen[0][colourIndex][pieceSquareIndex]
			case PAWN:
				result += pieceSquarePawn[colourIndex][pieceSquareIndex]
			}

		}
	}

	if b.whiteToMove {
		result += whiteMaterial - blackMaterial
	} else {
		result += blackMaterial - whiteMaterial
	}
	return result

}
