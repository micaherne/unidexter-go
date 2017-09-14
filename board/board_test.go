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

	b = FromFEN("7r/8/8/8/8/P7/8/R3K2R w KQ - 0 1")
	m = GeneratePieceMoves(b, 4)
	if len(m) != 7 {
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

	m = GeneratePieceMoves(b, 0x11)
	if len(m) != 0 {
		t.Errorf("Invalid piece should return 0 moves")
	}

}

func TestMakeMove(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	move := Move{0x14, 0x34}

	MakeMove(&b, move)

	if b.whiteToMove {
		t.Errorf("Should be black to move")
	}

	if b.squares[0x34] != WHITE|PAWN {
		t.Errorf("e4 should be white pawn")
	}
	if b.squares[0x14] != EMPTY {
		t.Errorf("e2 should be empty")
	}
}

func TestIsCheck(t *testing.T) {
	checkPositions := []string{
		"3rk3/8/8/8/8/8/8/3K4 w - -",
		"4k3/8/8/8/8/4n3/8/3K4 w - -",
		"4k3/8/8/7b/8/8/8/3K4 w - -",
		"4k3/8/8/8/q7/8/8/3K4 w - -",
		"4k3/8/8/8/8/8/4p3/3K4 w - -",
	}
	for i, pos := range checkPositions {
		b := FromFEN(pos)
		if !IsCheck(b, WHITE) {
			t.Errorf("White king should be in check in position %d", i)
		}
	}
	nonCheckPositions := []string{
		"4k3/8/8/8/8/4p3/8/3K4 w - -",
		"4k3/8/8/8/8/8/8/3KR3 w - -",
	}
	for i, pos := range nonCheckPositions {
		b := FromFEN(pos)
		if IsCheck(b, WHITE) {
			t.Errorf("White king should not be in check in position %d", i)
		}
	}
	blackCheckPositions := []string{
		"4k3/8/8/8/8/8/8/3KR3 w - -",
		"2Q1k3/4n3/8/8/8/8/8/3KR3 w - -",
	}
	for i, pos := range blackCheckPositions {
		b := FromFEN(pos)
		if !IsCheck(b, BLACK) {
			t.Errorf("Black king should be in check in position %d", i)
		}
	}
	blackNonCheckPositions := []string{
		"4k3/4n3/8/8/8/8/8/3KR3 w - -",
		"4k3/3nn3/8/8/8/8/8/3KR3 w - -",
	}
	for i, pos := range blackNonCheckPositions {
		b := FromFEN(pos)
		if IsCheck(b, BLACK) {
			t.Errorf("Black king should not be in check in position %d", i)
		}
	}
}

func TestRayDirection(t *testing.T) {
	r1 := RayDirection(0x00, 0x77)
	if r1 != NE {
		t.Errorf("a1 to h8 should be NE not %d", r1)
	}
	r2 := RayDirection(0x00, 0x07)
	if r2 != E {
		t.Errorf("a1 to h1 should be E not %d", r2)
	}
	r3 := RayDirection(0x07, 0x00)
	if r3 != W {
		t.Errorf("a8 to a1 should be W not %d", r3)
	}
	r4 := RayDirection(0x00, 0x12)
	if r4 != 0 {
		t.Errorf("a1 to b3 should return zero not %d", r4)
	}
	// Check that all returned values are correct.
	for from := 0; from < 128; from++ {
		if !LegalSquareIndex(from) {
			continue
		}
	ToLoop:
		for to := 0; to < 128; to++ {
			if !LegalSquareIndex(to) {
				continue
			}
			if from == to {
				continue
			}
			r := RayDirection(from, to)
			if r == 0 {
				continue
			}
			for i := 1; i < 8; i++ {
				if from+(i*r) == to {
					continue ToLoop
				}
			}
			t.Errorf("Can't get from %X to %X in direction %d", from, to, r)
			break
		}
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
