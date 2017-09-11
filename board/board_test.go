package board

import (
	"testing"
)

func TestBoardFromFEN(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if b.squares[0] != WHITE|ROOK {
		t.Error("Piece 0 should be white rook, not %s", GetPieceType(b.squares[0]))
	}
	if !b.whiteToMove {
		t.Error("Should be white to move")
	}
	if b.castling != 0x0F {
		t.Error("All castling should be available")
	}
	if b.ep != -1 {
		t.Errorf("e.p. square should be -1, not %d", b.ep)
	}

	// Test e.p.
	b2 := FromFEN("8/8/8/8/8/8/8/8 w KQkq c3 0 1")
	if b2.ep != 0x22 {
		t.Errorf("e.p. should be 0x22, not %X", b2.ep)
	}
}

func TestGetPieceType(t *testing.T) {
	p := WHITE | KNIGHT
	if pt := GetPieceType(p); pt != KNIGHT {
		t.Error("Piece should be a knight")
	}
	if pt := GetColour(p); pt != WHITE {
		t.Error("Piece should be white")
	}
}

func TestPieceFromNotation(t *testing.T) {
	if p := PieceFromNotation('n'); p != BLACK|KNIGHT {
		t.Error("n should be a black knight not %d", p)
	}
	if p := PieceFromNotation('R'); p != 8+ROOK {
		t.Error("R should be a white rook")
	}
}

func TestGenerateMoves(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	m := GenerateMoves(b)
	if len(m) != 20 {
		t.Errorf("Initial position should return 20 moves, not %d", len(m))
	}

	b.whiteToMove = false
	m = GenerateMoves(b)
	if len(m) != 20 {
		t.Errorf("Initial position with black to move should return 20 moves, not %d", len(m))
	}

	b = FromFEN("7r/8/8/8/8/P7/8/4K2R w K - 0 1")
	m = GeneratePieceMoves(b, 4)
	if len(m) != 6 {
		t.Errorf("Castling should be allowed", m)
	}
	m = GeneratePieceMoves(b, 32)
	if len(m) != 1 {
		t.Errorf("Pawn should have 1 move, not %d", len(m))
	}

	m = GeneratePieceMoves(b, 7)
	if len(m) != 9 {
		t.Errorf("Rook should have 9 moves, not %d", len(m))
	}

}

func TestNotationToSquareIndex(t *testing.T) {
	correct := map[string]int{
		"a1": 0x00,
		"a3": 0x20,
		"d2": 0x13,
		"h8": 0x77,
	}

	for notation, index := range correct {
		i := NotationToSquareIndex(notation)
		if i != index {
			t.Errorf("%s should be index %d, not %d", notation, index, i)
		}
	}

}
