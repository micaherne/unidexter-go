package board

import (
	"math/rand"
)

type zobristKeybase struct {
	PiecePosition [15][128]uint64
	WhiteToMove   uint64
	Castling      [4]uint64 // KQkq
	EpFile        [8]uint64
}

var ZobristKeys zobristKeybase

func InitZobristKeys() {
	r := rand.New(rand.NewSource(27092014))
	for piece := 1; piece < 7; piece++ {
		for square := 0; square < 128; square++ {
			if !LegalSquareIndex(square) {
				continue
			}
			ZobristKeys.PiecePosition[piece][square] = r.Uint64()
			ZobristKeys.PiecePosition[piece|BLACK][square] = r.Uint64()
		}
	}
	ZobristKeys.WhiteToMove = r.Uint64()
	for i := 0; i < 4; i++ {
		ZobristKeys.Castling[i] = r.Uint64()
	}
	for i := 0; i < 8; i++ {
		ZobristKeys.EpFile[i] = r.Uint64()
	}
}

func (b *Board) CalculateZobristHash() {
	b.zobristKey = 0
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank<<4 | file
			b.zobristKey ^= ZobristKeys.PiecePosition[b.squares[square]][square]
		}
	}
	if b.whiteToMove {
		b.zobristKey ^= ZobristKeys.WhiteToMove
	}
	for i := 0; i < 4; i++ {
		if b.castling&(1<<uint(i)) != 0 {
			b.zobristKey ^= ZobristKeys.Castling[i]
		}
	}
	if b.ep > 0 {
		b.zobristKey ^= ZobristKeys.EpFile[b.ep&7]
	}
}
