package board

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Piece int

const (
	EMPTY  = 0
	PAWN   = 1
	KNIGHT = 2
	BISHOP = 3
	ROOK   = 4
	QUEEN  = 5
	KING   = 6
)

const (
	BLACK = 0
	WHITE = 1
)

type Board struct {
	squares [128]Piece
}

func (p Piece) GetType() int {
	return int(p) & 7
}

func (p Piece) GetColour() int {
	return int(p) >> 3 & 1
}

func PieceFromNotation(symbol rune) Piece {
	var result Piece
	switch unicode.ToLower(symbol) {
	case 'p':
		result = Piece(PAWN)
	case 'n':
		result = Piece(KNIGHT)
	case 'b':
		result = Piece(BISHOP)
	case 'r':
		result = Piece(ROOK)
	case 'q':
		result = Piece(QUEEN)
	case 'k':
		result = Piece(KING)
	}
	if unicode.ToUpper(symbol) == symbol {
		result |= 8
	}
	return result
}

func FromFEN(fen string) Board {
	b := Board{}
	fenParts := strings.Split(fen, " ")
	boardParts := strings.Split(fenParts[0], "/")
	for rank, line := range boardParts {
		file := 0
		for _, char := range line {
			if n, err := strconv.Atoi(string(char)); err == nil {
				file += n
			} else {
				b.squares[(7-rank)*16+file] = PieceFromNotation(char)
				fmt.Printf("%d %d ", (7-rank)*16+file, PieceFromNotation(char))
				file++
			}
		}
		fmt.Println("")
	}
	return b
}
