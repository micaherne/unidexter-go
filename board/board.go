package board

import (
	"bytes"
	"fmt"
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

var PROMOTIONPIECES = []int{QUEEN, ROOK, BISHOP, KNIGHT}

type Board struct {
	squares     [128]int
	whiteToMove bool
	castling    int // KQkq
	ep          int
	blackKing   int
	whiteKing   int
	fullMove    int
	halfMove    int
	moveHistory []MoveUndo
}

type Move struct {
	from      int
	to        int
	promotion int
}

type MoveUndo struct {
	from        int
	to          int
	isPromotion bool // We don't care what the promoted piece is.
	captured    int
	ep          int
	halfMove    int
	castling    int
}

func (m Move) String() string {
	return SquareIndexToNotation(m.from) + SquareIndexToNotation(m.to)
}

const InitialPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func GetPieceType(p int) int {
	return p & 7
}

func GetColour(p int) int {
	return p & WHITE
}

func GetOpponentColour(p int) int {
	return p&WHITE ^ WHITE
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

func PieceToNotation(piece int) string {
	notation := " "
	colour := GetColour(piece)
	pieceType := GetPieceType(piece)
	switch pieceType {
	case PAWN:
		notation = "p"
	case KNIGHT:
		notation = "n"
	case BISHOP:
		notation = "b"
	case ROOK:
		notation = "r"
	case QUEEN:
		notation = "q"
	case KING:
		notation = "k"
	}
	if colour == WHITE {
		notation = strings.ToUpper(notation)
	}
	return notation
}

func FromFEN(fen string) *Board {
	b := Board{}
	fenParts := strings.Split(fen, " ")
	boardParts := strings.Split(fenParts[0], "/")
	for rank, line := range boardParts {
		file := 0
		for _, char := range line {
			if n, err := strconv.Atoi(string(char)); err == nil {
				file += n
			} else {
				piece := PieceFromNotation(char)
				square := (7-rank)*16 + file
				b.squares[square] = piece
				if piece == WHITE|KING {
					b.whiteKing = square
				} else if piece == BLACK|KING {
					b.blackKing = square
				}
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

	if len(fenParts) > 4 {
		if halfMove, err := strconv.Atoi(fenParts[4]); err == nil {
			b.halfMove = halfMove
		}
	} else {
		b.halfMove = 0
	}

	if len(fenParts) > 5 {
		if fullMove, err := strconv.Atoi(fenParts[5]); err == nil {
			b.fullMove = fullMove
		}
	} else {
		b.fullMove = 1
	}

	return &b
}

func ToFEN(b *Board) string {
	var result bytes.Buffer
	for rankStart := 0x70; rankStart >= 0; rankStart -= 16 {
		empties := 0
		for i := 0; i < 8; i++ {
			square := rankStart + i
			if b.squares[square] == EMPTY {
				empties++
				continue
			}
			if empties > 0 {
				result.WriteString(strconv.Itoa(empties))
				empties = 0
			}
			result.WriteString(PieceToNotation(b.squares[square]))
		}
		if empties > 0 {
			result.WriteString(strconv.Itoa(empties))
		}
		if rankStart > 0 {
			result.WriteString("/")
		}
	}

	result.WriteString(" ")
	if b.whiteToMove {
		result.WriteString("w")
	} else {
		result.WriteString("b")
	}
	result.WriteString(" ")
	for i, symbol := range []rune{'K', 'Q', 'k', 'q'} {
		if b.castling&(1<<uint(3-i)) > 0 {
			result.WriteRune(symbol)
		}
	}
	result.WriteString(" ")
	if b.ep > 0 {
		result.WriteString(SquareIndexToNotation(b.ep))
	} else {
		result.WriteString("-")
	}
	result.WriteString(" ")
	result.WriteString(fmt.Sprintf("%d %d", b.halfMove, b.fullMove))

	return result.String()
}

// GenerateMoves generates pseudo-legal moves for the position given
func GenerateMoves(b *Board) []Move {
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

func GeneratePieceMoves(b *Board, i int) []Move {
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
			kingMoves = append(kingMoves, Move{i, i + 2, EMPTY})
		}
		if castleQueensideAllowed && CastlingLegal(b, i, W) {
			kingMoves = append(kingMoves, Move{i, i - 2, EMPTY})
		}
		return kingMoves
	}

	// TODO: Raise error
	return make([]Move, 0)
}

func GenerateSingleMoves(b *Board, i int, offsets []int) []Move {
	result := make([]Move, 0)
	ownColour := GetColour(b.squares[i])
	for _, offset := range offsets {
		if toSquare := i + offset; LegalSquareIndex(toSquare) {
			if b.squares[toSquare] == 0 || GetColour(b.squares[toSquare]) != ownColour {
				result = append(result, Move{i, toSquare, EMPTY})
			}
		}
	}
	return result
}

func GenerateSlides(b *Board, i int, offsets []int) []Move {
	result := make([]Move, 0)
	ownColour := GetColour(b.squares[i])
	for _, offset := range offsets {
		for toSquare := i + offset; toSquare >= 0 && LegalSquareIndex(toSquare); toSquare += offset {
			if b.squares[toSquare] == 0 {
				result = append(result, Move{i, toSquare, EMPTY})
			} else if GetColour(b.squares[toSquare]) != ownColour {
				result = append(result, Move{i, toSquare, EMPTY})
				break
			} else {
				break
			}
		}
	}
	return result
}

func GeneratePawnMoves(b *Board, i int, forward int, captureOffsets []int, onHomeRank bool) []Move {
	ownColour := GetColour(b.squares[i])
	var result []Move
	if onHomeRank {
		result = GeneratePawnSlides(b, i, forward)
	} else {
		result = make([]Move, 0)
		if toSquare := i + forward; toSquare&0x88 == 0 && b.squares[toSquare] == 0 {
			result = append(result, Move{i, toSquare, EMPTY})
		}
	}

	for _, offset := range captureOffsets {
		toSquare := i + offset
		if (toSquare)&0x88 > 0 {
			continue
		}
		if b.ep == toSquare {
			result = append(result, Move{i, toSquare, EMPTY})
			continue
		}
		if b.squares[toSquare] == 0 {
			continue
		}
		if GetColour(b.squares[toSquare]) == ownColour {
			continue
		}
		result = append(result, Move{i, toSquare, EMPTY})
	}

	// Convert moves to promotion moves where required.
	// Loads start of array with queens, then rooks etc.
	if ownColour == WHITE && i&0xF0 == 0x60 || ownColour == BLACK && i&0xF0 == 0x10 {
		var promotions []Move
		promotions = make([]Move, len(result)*4)
		for pieceIndex, piece := range PROMOTIONPIECES {
			for moveIndex, move := range result {
				promotions[pieceIndex*len(result)+moveIndex] = Move{move.from, move.to, ownColour | piece}
			}
		}
		return promotions
	}

	return result
}

// GeneratePawnSlides returns the valid non-capturing moves for
// a pawn on the home rank. Behaviour is undefined for pawns on
// any other rank.
func GeneratePawnSlides(b *Board, i int, offset int) []Move {
	result := make([]Move, 0)
	toSquare := i + offset
	if b.squares[toSquare] == 0 {
		result = append(result, Move{i, toSquare, EMPTY})
		toSquare += offset
		if b.squares[toSquare] == 0 {
			result = append(result, Move{i, toSquare, EMPTY})
		}
	}

	return result
}

// LegalMove decides whether a pseudo-legal move is actually legal
// i.e. does it result in the moving side being in check?
// TODO: This is a very slow implementation - do static calculation
func LegalMove(b *Board, move Move) bool {
	colourMoving := ColourToMove(b)
	MakeMove(b, move)
	result := !IsCheck(b, colourMoving)
	UndoMove(b)
	return result
}

func ColourToMove(b *Board) int {
	if b.whiteToMove {
		return WHITE
	}
	return BLACK
}

func MakeMove(b *Board, move Move) {
	undo := MoveUndo{
		from:        move.from,
		to:          move.to,
		isPromotion: move.promotion != EMPTY,
		captured:    b.squares[move.to],
		ep:          b.ep,
		halfMove:    b.halfMove,
		castling:    b.castling,
	}

	resetEp := true

	movedPiece := GetPieceType(b.squares[move.from])

	if movedPiece == PAWN {
		if move.to == b.ep && (move.to&0x0F != move.from&0x0F) {
			capturedSquare := move.from&0xF0 | move.to&0x0F
			undo.captured = b.squares[capturedSquare]
			b.squares[capturedSquare] = EMPTY
		} else if offset := move.to - move.from; offset == 32 || offset == -32 {
			b.ep = move.from + offset/2
			resetEp = false
		}
	} else if movedPiece == KING {
		if GetColour(b.squares[move.from]) == WHITE {
			b.whiteKing = move.to
			b.castling &= 0x03
		} else {
			b.blackKing = move.to
			b.castling &= 0x0C
		}
		offset := move.to - move.from

		if offset == 2 { // Kingside castling
			// Rook
			b.squares[move.from+1] = b.squares[move.to+1]
			b.squares[move.to+1] = EMPTY
		} else if offset == -2 { // Queenside castling
			// Rook
			b.squares[move.from-1] = b.squares[move.to-2]
			b.squares[move.to-2] = EMPTY
		}
	} else if movedPiece == ROOK {
		if b.whiteToMove {
			if move.from == 0x00 && b.castling&0x04 == 0x04 {
				b.castling &= ^0x04
			}
			if move.from == 0x07 && b.castling&0x08 == 0x08 {
				b.castling &= ^0x08
			}
		} else {
			if move.from == 0x70 && b.castling&0x01 == 0x01 {
				b.castling &= ^0x01
			}
			if move.from == 0x77 && b.castling&0x02 == 0x02 {
				b.castling &= ^0x02
			}
		}
	}

	// Update castling if rook captured
	if b.squares[move.to]&ROOK == ROOK {
		if GetColour(b.squares[move.from]) == WHITE {
			if move.to == 0x70 {
				b.castling &= ^0x01
			} else if move.to == 0x77 {
				b.castling &= ^0x02
			}
		} else {
			if move.to == 0x00 {
				b.castling &= ^0x04
			} else if move.to == 0x07 {
				b.castling &= ^0x08
			}
		}
	}

	if resetEp {
		b.ep = -1
	}

	b.moveHistory = append(b.moveHistory, undo)

	// Move the actual piece
	if move.promotion == EMPTY { // Checking undo.promotion would be marginally quicker.
		b.squares[move.to] = b.squares[move.from]
	} else {
		b.squares[move.to] = move.promotion
	}

	b.squares[move.from] = EMPTY
	b.whiteToMove = !b.whiteToMove
}

func UndoMove(b *Board) {
	a := b.moveHistory
	var lastMove MoveUndo
	lastMove, b.moveHistory = a[len(a)-1], a[:len(a)-1]

	b.squares[lastMove.from] = b.squares[lastMove.to]
	b.squares[lastMove.to] = lastMove.captured

	b.whiteToMove = !b.whiteToMove
	b.ep = lastMove.ep
	b.halfMove = lastMove.halfMove
	b.castling = lastMove.castling

	// Promotion
	if lastMove.isPromotion {
		b.squares[lastMove.from] = PAWN | GetColour(b.squares[lastMove.from])
	}

	movedPiece := GetPieceType(b.squares[lastMove.from])

	if movedPiece == PAWN && lastMove.to == lastMove.ep && (lastMove.to&0x0F != lastMove.from&0x0F) {
		capturedSquare := lastMove.from&0xF0 | lastMove.to&0x0F
		b.squares[capturedSquare] = lastMove.captured
		b.squares[lastMove.to] = EMPTY
	} else if movedPiece == KING {
		if GetColour(b.squares[lastMove.from]) == WHITE {
			b.whiteKing = lastMove.from
		} else {
			b.blackKing = lastMove.from
		}
		offset := lastMove.to - lastMove.from

		if offset == 2 { // Kingside castling
			// Rook
			b.squares[lastMove.to+1] = b.squares[lastMove.from+1]
			b.squares[lastMove.from+1] = EMPTY
		} else if offset == -2 { // Queenside castling
			// Rook
			b.squares[lastMove.to-2] = b.squares[lastMove.from-1]
			b.squares[lastMove.from-1] = EMPTY
		}
	}
}

func MakeMoveFromNotation(b *Board, move string) {
	// TODO: Support promotions and castling
	m := Move{NotationToSquareIndex(move[:2]), NotationToSquareIndex(move[2:4]), EMPTY}
	MakeMove(b, m)
}

func LegalSquareIndex(i int) bool {
	return i >= 0 && i&0x88 == 0
}

// CastlingLegal checks that the correct squares are empty
// and not attacked.
func CastlingLegal(b *Board, i int, direction int) bool {
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

	if IsCheck(b, GetColour(b.squares[i])) {
		return false
	}

	// Check attacks on intermediate squares
	colour := GetOpponentColour(b.squares[i])
	if IsAttacked(b, i+direction, colour) {
		return false
	}
	if IsAttacked(b, i+direction+direction, colour) {
		return false
	}
	return true
}

// IsAttacked determines whether the given square is attacked by
// the given colour.
// Note that for e.p. only the ep square returns true, not the attacked pawn's square
func IsAttacked(b *Board, square int, colour int) bool {

	var pawnAttacks []int
	if colour == WHITE {
		pawnAttacks = NDIAGONALS
	} else { // Assume calling code is correct to avoid extra if
		pawnAttacks = SDIAGONALS
	}

	// Knights
	for _, knightMove := range KNIGHTMOVES {
		testSquare := square - knightMove
		if LegalSquareIndex(testSquare) && b.squares[testSquare] == colour|KNIGHT {
			return true
		}
	}

	// Rays
	for _, dir := range DIAGONALS {
		for i := 1; i < 8; i++ {
			testSquare := square - (i * dir)
			if !LegalSquareIndex(testSquare) {
				break
			}
			if b.squares[testSquare] != EMPTY {
				if b.squares[testSquare] == colour|BISHOP {
					return true
				}
				if b.squares[testSquare] == colour|QUEEN {
					return true
				}
				break
			}
		}
	}
	for _, dir := range LINES {
		for i := 1; i < 8; i++ {
			testSquare := square - (i * dir)
			if !LegalSquareIndex(testSquare) {
				break
			}
			if b.squares[testSquare] != EMPTY {
				if b.squares[testSquare] == colour|ROOK {
					return true
				}
				if b.squares[testSquare] == colour|QUEEN {
					return true
				}
				break
			}
		}
	}

	// Pawns
	for _, pawnMove := range pawnAttacks {
		testSquare := square - pawnMove
		if LegalSquareIndex(testSquare) && b.squares[testSquare] == colour|PAWN {
			return true
		}
	}

	// King
	var offset int
	if colour == WHITE {
		offset = square - b.whiteKing
	} else {
		offset = square - b.blackKing
	}

	for _, dir := range DIAGONALSANDLINES {
		if offset == dir {
			return true
		}
	}

	return false
}

// IsCheck tests whether the given colour is in check.
func IsCheck(b *Board, colour int) bool {
	var kingPosition int
	var opponentColour int
	if colour == WHITE {
		kingPosition = b.whiteKing
		opponentColour = BLACK
	} else { // Assume calling code is correct to avoid extra if
		kingPosition = b.blackKing
		opponentColour = WHITE
	}
	return IsAttacked(b, kingPosition, opponentColour)
}

// RayDirection gets the direction from one square to another
// or zero if they are not on the same line or diagonal.
// Assumes from and to are distinct, valid squares.
func RayDirection(from int, to int) int {
	result := 0
	if from&0xF0 == to&0xF0 {
		if from&0x0F > to&0x0F {
			return W
		}
		return E
	} else if from&0x0F == to&0x0F {
		if from&0xF0 > to&0xF0 {
			return S
		}
		return N
	} else if xdiff, ydiff := from&0x0F-to&0x0F, from&0xF0>>4-to&0xF0>>4; xdiff^2 == ydiff^2 {
		if xdiff > 0 {
			result += S
		} else {
			result += N
		}
		if ydiff > 0 {
			result += W
		} else {
			result += E
		}
	}
	return result
}

func NotationToSquareIndex(notation string) int {
	return int((notation[0] - "a"[0]) + (notation[1]-"1"[0])*16)
}

func SquareIndexToNotation(square int) string {
	return string(byte(square&0x0F)+"a"[0]) + string(byte(square&0xF0>>4)+"1"[0])
}

func (b *Board) String() string {
	result := bytes.Buffer{}
	separator := strings.Repeat("+-", 8) + "+\n"
	for rankStart := 0x70; rankStart >= 0; rankStart -= 16 {

		result.Write([]byte(separator + "|"))

		for squareIndex := rankStart; squareIndex < rankStart+8; squareIndex++ {
			// TODO: Change string concatenation to direct buffer write
			notation := PieceToNotation(b.squares[squareIndex])
			notation += "|"
			result.Write([]byte(notation))
		}

		result.Write([]byte("\n"))
	}

	result.Write([]byte(separator))

	otherStuff := fmt.Sprintf("White to move?: %t\nCastling (KQkq): %b\nEn passent: %X\nHalf move: %d", b.whiteToMove, b.castling, b.ep, b.halfMove)
	result.Write([]byte(otherStuff))
	return result.String()
}
