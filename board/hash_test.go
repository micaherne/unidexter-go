package board

import (
	"testing"
)

func TestInitZobristKeys(t *testing.T) {
	InitZobristKeys()
	if ZobristKeys.PiecePosition[3][0x44] == 0 {
		t.Errorf("key should not be zero")
	}

	if ZobristKeys.PiecePosition[0][0x44] != 0 {
		t.Errorf("key should be zero for empty piece")
	}

	b := FromFEN(InitialPositionFEN)
	b.CalculateZobristHash()
	if b.zobristKey == 0 {
		t.Errorf("key should not be zero")
	}
	initialPosKey := b.zobristKey

	MakeMove(b, Move{0x14, 0x34, 0})

	b.CalculateZobristHash()

	move2Key := b.zobristKey

	if initialPosKey == move2Key {
		t.Errorf("keys should not be equal")
	}

	UndoMove(b)

	b.CalculateZobristHash()
	if b.zobristKey != initialPosKey {
		t.Errorf("keys should be equal after undo")
	}

}
