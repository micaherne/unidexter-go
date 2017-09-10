package board

import (
	"testing"
)

func TestBoardFromFEN(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if b.squares[0].GetType() != ROOK+(WHITE*8) {
		t.Error("Piece 0 should be white rook")
	}
}

func TestGetType(t *testing.T) {
	p := Piece((WHITE * 8) | KNIGHT)
	if pt := p.GetType(); pt != KNIGHT {
		t.Error("Piece should be a knight")
	}
	if pt := p.GetColour(); pt != WHITE {
		t.Error("Piece should be white")
	}
}

func TestPieceFromNotation(t *testing.T) {
	if p := PieceFromNotation('n'); p != BLACK|KNIGHT {
		t.Error("n should be a black knight not %d", p)
	}
}
