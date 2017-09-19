package board

import (
	"fmt"
	"testing"
)

func TestFromFEN(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if b.squares[0] != WHITE|ROOK {
		t.Errorf("Piece 0 should be white rook, not %d", GetPieceType(b.squares[0]))
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
	if b.halfMove != 0 {
		t.Errorf("halfMove should be 0, not %d", b.halfMove)
	}
	if b.fullMove != 1 {
		t.Errorf("fullMove should be 1, not %d", b.fullMove)
	}

	// Test e.p.
	b2 := FromFEN("8/8/8/8/8/8/8/8 w KQkq c3 0 1")
	if b2.ep != 0x22 {
		t.Errorf("e.p. should be 0x22, not %X", b2.ep)
	}

	b3 := FromFEN("r1bqk1nr/ppp2ppp/2n5/1BbpP3/8/5N2/PPPP1PPP/RNBQK2R w KQkq d6")
	if b3.ep != 0x53 {
		t.Errorf("e.p. should be 0x53, not %X", b3.ep)
	}
}

func TestToFEN(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", // kiwipete
		"r3k2r/p1ppqpb1/bn2Pnp1/4N3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1",  // kiwipete + d5d6
	}
	for i, fen := range fens {
		b := FromFEN(fen)
		if fenOut := ToFEN(b); fenOut != fen {
			t.Errorf("Test %d. %s\ndoes not equal\n%s", i, fen, fenOut)
		}
	}
}

func TestGetPieceType(t *testing.T) {
	p := WHITE | KNIGHT
	if pt := GetPieceType(p); pt != KNIGHT {
		t.Error("Piece should be a knight")
	}
	if c := GetColour(p); c != WHITE {
		t.Error("Piece should be white")
	}
	if oc := GetOpponentColour(p); oc != BLACK {
		t.Error("Opponent should be black")
	}
	p = BLACK | PAWN
	if pt := GetPieceType(p); pt != PAWN {
		t.Error("Piece should be a pawnt")
	}
	if c := GetColour(p); c != BLACK {
		t.Error("Piece should be black")
	}
	if oc := GetOpponentColour(p); oc != WHITE {
		t.Error("Opponent should be white")
	}
}

func TestPieceFromNotation(t *testing.T) {
	if p := PieceFromNotation('n'); p != BLACK|KNIGHT {
		t.Errorf("n should be a black knight not %d", p)
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
		t.Errorf("Castling should be allowed: %X", m)
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

	b = FromFEN("8/8/8/8/8/8/6k1/4K2R w K - 0 1")
	m = GenerateMoves(b)
	// Note - 2 of these moves are illegal
	if len(m) != 14 {
		t.Errorf("Should be 14, not %d", len(m))
	}
	// "Kiwipete" perft position
	/* b = FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -")
	m = GeneratePieceMoves(b, 0x10)
	fmt.Println(m)*/

}

func TestMakeMove(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	move := Move{0x14, 0x34, EMPTY}

	MakeMove(b, move)

	if b.whiteToMove {
		t.Errorf("Should be black to move")
	}

	if b.squares[0x34] != WHITE|PAWN {
		t.Errorf("e4 should be white pawn")
	}
	if b.squares[0x14] != EMPTY {
		t.Errorf("e2 should be empty")
	}

	// Castling king side
	b = FromFEN("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/2N2N2/PPPP1PPP/R1BQK2R w KQkq -")
	move = Move{0x04, 0x06, EMPTY}

	MakeMove(b, move)

	if b.squares[0x05] != WHITE|ROOK {
		t.Errorf("f1 should be white rook after castling")
	}
	if b.squares[0x06] != WHITE|KING {
		t.Errorf("g1 should be white king after castling")
	}
	if b.squares[0x07] != EMPTY {
		t.Errorf("h1 should be empty after castling")
	}

	b = FromFEN("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/2N2N2/PPPP1PPP/R1BQK2R b KQkq -")

	move = Move{0x74, 0x76, EMPTY}
	MakeMove(b, move)

	if b.squares[0x75] != BLACK|ROOK {
		t.Errorf("f8 should be black rook after castling")
	}
	if b.squares[0x76] != BLACK|KING {
		t.Errorf("g8 should be black king after castling")
	}
	if b.squares[0x77] != EMPTY {
		t.Errorf("h8 should be empty after castling")
	}

	// Castling queens side
	b = FromFEN("r3kbnr/pppb1ppp/2np4/4p3/4P2q/2NPBQ2/PPP2PPP/R3KBNR w KQkq -")

	move = Move{0x04, 0x02, EMPTY}

	MakeMove(b, move)

	if b.squares[0x03] != WHITE|ROOK {
		t.Errorf("d1 should be white rook after castling")
	}
	if b.squares[0x02] != WHITE|KING {
		t.Errorf("c1 should be white king after castling")
	}
	if b.squares[0x00] != EMPTY {
		t.Errorf("a1 should be empty after castling")
	}

	// e.p.
	b = FromFEN("r1bqk1nr/ppp2ppp/2n5/1BbpP3/8/5N2/PPPP1PPP/RNBQK2R w KQkq d6 0 1")
	move = Move{0x44, 0x53, EMPTY}
	MakeMove(b, move)
	if b.squares[0x53] != WHITE|PAWN {
		t.Errorf("d6 should be a white pawn")
	}
	if b.squares[0x43] != EMPTY {
		t.Errorf("d5 should be empty")
	}
	if b.moveHistory[0].captured != BLACK|PAWN {
		t.Errorf("captured should be black pawn")
	}

	b = FromFEN("r2k3r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/1R2K2R w K -")
	MakeMove(b, Move{0x62, 0x42, EMPTY})
	if b.ep != 0x52 {
		t.Errorf("e.p. square should be c6 after c5c7, not %X", b.ep)
	}

	// castling
	// white
	correct := map[Move]int{
		{0x00, 0x10, EMPTY}: 0x0B,
		{0x07, 0x37, EMPTY}: 0x07,
	}
	for move, castling := range correct {
		b = FromFEN("r3k2r/8/8/8/8/8/8/R3K2R w KQkq -")
		MakeMove(b, move)
		if b.castling != castling {
			t.Errorf("Castling after %s should be %b, not %b", move, castling, b.castling)
		}
	}
	// black
	correct = map[Move]int{
		{0x70, 0x60, EMPTY}: 0x0E,
		{0x77, 0x37, EMPTY}: 0x0D,
	}
	for move, castling := range correct {
		b = FromFEN("r3k2r/8/8/8/8/8/8/R3K2R b KQkq -")
		MakeMove(b, move)
		if b.castling != castling {
			t.Errorf("Castling after %s should be %b, not %b", move, castling, b.castling)
		}
	}

	b = FromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/1R2K2R b Kkq -")
	if b.castling != 11 {
		t.Errorf("Kkq should be 11")
	}
	MakeMove(b, Move{0x74, 0x73, EMPTY})
	if b.castling != 8 {
		t.Errorf("Castling should be 8")
	}

}

func TestMakeMoveFromNotation(t *testing.T) {
	b := FromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	MakeMoveFromNotation(b, "e2e4")

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

func TestUndoMove(t *testing.T) {
	tests := map[string][]Move{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1": []Move{
			Move{0x14, 0x34, EMPTY},
			Move{0x64, 0x44, EMPTY},
		},
		"r1bqk1nr/pppp1ppp/2n5/1Bb1p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq -": []Move{
			Move{0x04, 0x06, EMPTY},
		},
		// e.p.
		"r1bqk1nr/ppp2ppp/2n5/1BbpP3/8/5N2/PPPP1PPP/RNBQK2R w KQkq d6": []Move{
			Move{0x44, 0x53, EMPTY},
		},
	}

	for fen, moves := range tests {
		a := FromFEN(fen)
		b := FromFEN(fen)

		for _, move := range moves {
			MakeMove(b, move)
		}

		for i := 0; i < len(moves); i++ {
			UndoMove(b)
		}

		// Can't use reflect.DeepEqual as moveHistory capacities are different
		if fmt.Sprint(a) != fmt.Sprint(b) {
			fmt.Println(a, b)
			t.Error("Undo failed")
		}
	}

}

func TestIsAttacked(t *testing.T) {
	b := FromFEN("7r/8/8/8/8/P7/8/R3K2R w KQ - 0 1")
	if IsAttacked(b, 0x03, BLACK) {
		t.Error("a4 is not attacked by black")
	}
	if !IsAttacked(b, 0x07, BLACK) {
		t.Error("h1 is attacked by black")
	}
	b = FromFEN("8/8/8/8/8/8/6k1/4K2R w K -")
	if !IsAttacked(b, 0x06, BLACK) {
		t.Error("g1 is attacked by black")
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

func TestSquareIndexToNotation(t *testing.T) {
	correct := map[int]string{
		0x00: "a1",
		0x01: "b1",
	}
	for square, name := range correct {
		if n := SquareIndexToNotation(square); n != name {
			t.Errorf("Name of %X should be %s, not %s", square, name, n)
		}
	}
}
