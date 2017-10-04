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

}

func TestCalculateZobristHash(t *testing.T) {

	InitZobristKeys()

	correct := map[string][]Move{
		InitialPositionFEN:           []Move{{0x14, 0x34, 0}},
		"4k2r/8/8/8/8/8/8/4K3 b k -": []Move{{0x74, 0x76, 0}},
	}

	for fen, moves := range correct {
		b := FromFEN(fen)
		b.CalculateZobristHash()
		if b.zobristKey == 0 {
			t.Errorf("key should not be zero %s %b", b, b.zobristKey)
		}
		initialPosKey := b.zobristKey

		for _, move := range moves {
			MakeMove(b, move)

			b.CalculateZobristHash()

			move2Key := b.zobristKey

			if initialPosKey == move2Key {
				t.Errorf("keys should not be equal")
			}

			UndoMove(b)

			b.CalculateZobristHash()
			if b.zobristKey != initialPosKey {
				t.Errorf("keys should be equal after undo, %s %s", fen, move)
			}
		}

	}

}

func BenchmarkInitZobristKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		InitZobristKeys()
	}
}
