package board

import (
	"strconv"
	"strings"
	"unicode"
)

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
	WHITE = 8 // Flag
)

const (
	N   = 16
	E   = 1
	S   = -16
	W   = -1
	NE  = N + E // 17
	SE  = S + E // -15
	SW  = S + W // -17
	NW  = N + W //15
	NNE = N + NE
	ENE = E + NE
	ESE = E + SE
	SSE = S + SE
	SSW = S + SW
	WSW = W + SW
	WNW = W + NW
	NNW = N + NW
)

var (
	DIAGONALS         = []int{NE, SE, SW, NW}
	LINES             = []int{N, S, E, W}
	DIAGONALSANDLINES = append(DIAGONALS, LINES...)
	SDIAGONALS        = []int{SE, SW}
	NDIAGONALS        = []int{NE, NW}
	KNIGHTMOVES       = []int{NNE, ENE, ESE, SSE, SSW, WSW, WNW, NNW}
)

type Board struct {
	squares     [128]int
	whiteToMove bool
	castling    int // KQkq
	ep          int
}

func GetPieceType(p int) int {
	return p & 7
}

func GetColour(p int) int {
	return p & WHITE
}

func PieceFromNotation(symbol rune) int {
	var result int
	switch unicode.ToLower(symbol) {
	case 'p':
		result = PAWN
	case 'n':
		result = KNIGHT
	case 'b':
		result = BISHOP
	case 'r':
		result = ROOK
	case 'q':
		result = QUEEN
	case 'k':
		result = KING
	}
	if unicode.ToUpper(symbol) == symbol {
		result |= WHITE
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
				file++
			}
		}
	}

	b.whiteToMove = (fenParts[1] == "w")

	castlingBits := map[string]int{
		"K": 8,
		"Q": 4,
		"k": 2,
		"q": 1,
	}

	for letter, bit := range castlingBits {
		if strings.Contains(fenParts[2], letter) {
			b.castling |= bit
		}
	}

	if fenParts[3] == "-" {
		b.ep = -1
	} else {
		b.ep = NotationToSquareIndex(fenParts[3])
	}

	return b
}

type Move struct {
	from int
	to   int
}

type MoveUndo struct {
	from     int
	to       int
	captured int
	ep       int
	halfMove int
	castling int
}

func GenerateMoves(b Board) []Move {
	result := make([]Move, 0)
	for i := 0; i < 128; i++ {
		if b.squares[i] == EMPTY {
			continue
		}

		if (b.whiteToMove && GetColour(b.squares[i]) == WHITE) || (!b.whiteToMove && GetColour(b.squares[i]) == BLACK) {
			moves := GeneratePieceMoves(b, i)
			result = append(result, moves...)
		}
	}
	return result
}

func GeneratePieceMoves(b Board, i int) []Move {
	piece := b.squares[i]
	switch GetPieceType(piece) {
	case PAWN:
		if b.whiteToMove {
			return GeneratePawnMoves(b, i, N, NDIAGONALS, i&0xF0 == 0x10)
		} else {
			return GeneratePawnMoves(b, i, S, SDIAGONALS, i&0xF0 == 0x60)
		}
	case KNIGHT:
		return GenerateSingleMoves(b, i, KNIGHTMOVES)
	case BISHOP:
		return GenerateSlides(b, i, DIAGONALS)
	case ROOK:
		return GenerateSlides(b, i, LINES)
	case QUEEN:
		return GenerateSlides(b, i, DIAGONALSANDLINES)
	case KING:
		kingMoves := GenerateSingleMoves(b, i, DIAGONALSANDLINES)
		var castleKingsideAllowed, castleQueensideAllowed bool
		if b.whiteToMove {
			castleKingsideAllowed = b.castling&8 == 8
			castleQueensideAllowed = b.castling&4 == 4
		} else {
			castleKingsideAllowed = b.castling&2 == 2
			castleQueensideAllowed = b.castling&1 == 1
		}
		if castleKingsideAllowed && CastlingLegal(b, i, E) {
			kingMoves = append(kingMoves, Move{i, i + 2})
		}
		if castleQueensideAllowed && CastlingLegal(b, i, W) {
			kingMoves = append(kingMoves, Move{i, i - 2})
		}
		return kingMoves
	}

	// TODO: Raise error
	return make([]Move, 0)
}

func GenerateSingleMoves(b Board, i int, offsets []int) []Move {
	result := make([]Move, 0)
	ownColour := GetColour(b.squares[i])
	for _, offset := range offsets {
		if toSquare := i + offset; LegalSquareIndex(toSquare) {
			if b.squares[toSquare] == 0 || GetColour(b.squares[toSquare]) != ownColour {
				result = append(result, Move{i, toSquare})
			}
		}
	}
	return result
}

func GenerateSlides(b Board, i int, offsets []int) []Move {
	result := make([]Move, 0)
	ownColour := GetColour(b.squares[i])
	for _, offset := range offsets {
		for toSquare := i + offset; toSquare >= 0 && LegalSquareIndex(toSquare); toSquare += offset {
			if b.squares[toSquare] == 0 {
				result = append(result, Move{i, toSquare})
			} else if GetColour(b.squares[toSquare]) != ownColour {
				result = append(result, Move{i, toSquare})
				break
			} else {
				break
			}
		}
	}
	return result
}

func GeneratePawnMoves(b Board, i int, forward int, captureOffsets []int, onHomeRank bool) []Move {
	ownColour := GetColour(b.squares[i])
	var result []Move
	if onHomeRank {
		result = GeneratePawnSlides(b, i, forward)
	} else {
		result = make([]Move, 0)
		if toSquare := i + forward; toSquare&0x88 == 0 && b.squares[toSquare] == 0 {
			result = append(result, Move{i, toSquare})
		}
	}

	for _, offset := range captureOffsets {
		toSquare := i + offset
		if (toSquare)&0x88 > 0 {
			continue
		}
		if b.ep == toSquare {
			result = append(result, Move{i, toSquare})
			continue
		}
		if b.squares[toSquare] == 0 {
			continue
		}
		if GetColour(b.squares[toSquare]) == ownColour {
			continue
		}
		result = append(result, Move{i, toSquare})
	}
	return result
}

// GeneratePawnSlides returns the valid non-capturing moves for
// a pawn on the home rank. Behaviour is undefined for pawns on
// any other rank.
func GeneratePawnSlides(b Board, i int, offset int) []Move {
	result := make([]Move, 0)
	toSquare := i + offset
	if b.squares[toSquare] == 0 {
		result = append(result, Move{i, toSquare})
		toSquare := i + offset
		if b.squares[toSquare] == 0 {
			result = append(result, Move{i, toSquare})
		}
	}

	return result
}

func LegalSquareIndex(i int) bool {
	return i >= 0 && i&0x88 == 0
}

// CastlingLegal checks that the correct squares are empty
// and not attacked.
func CastlingLegal(b Board, i int, direction int) bool {
	var rookPosition int
	switch direction {
	case W:
		if (b.whiteToMove && b.castling&4 == 0) || (!b.whiteToMove && b.castling&1 == 0) {
			return false
		}
		rookPosition = i - 4
	case E:
		if (b.whiteToMove && b.castling&8 == 0) || (!b.whiteToMove && b.castling&2 == 0) {
			return false
		}
		rookPosition = i + 3
	default:
		// TODO: Throw error
		return false
	}
	for sq := i + direction; sq != rookPosition; sq += direction {
		if b.squares[sq] != EMPTY {
			return false
		}
	}
	// TODO: Check attacks on intermediate squares
	return true
}

func NotationToSquareIndex(notation string) int {
	return int((notation[0] - "a"[0]) + (notation[1]-"1"[0])*16)
}
