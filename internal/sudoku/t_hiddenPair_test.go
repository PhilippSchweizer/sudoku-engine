package sudoku

import "testing"

func TestApplyHiddenPairs(t *testing.T) {
	t.Run("single_hidden_pair_in_row", func(t *testing.T) {
		// Hidden pair in row 0: digits 1 and 2 appear only in (0,1) and (0,5),
		// and those cells have an extra candidate 3 that should be removed.
		b := New()
		b.AddCandidate(0, 1, 1)
		b.AddCandidate(0, 1, 2)
		b.AddCandidate(0, 1, 3)
		b.AddCandidate(0, 5, 1)
		b.AddCandidate(0, 5, 2)
		b.AddCandidate(0, 5, 3)
		// Make digit 3 appear elsewhere in the row so it's not part of the hidden pair.
		b.AddCandidate(0, 0, 3)
		b.AddCandidate(0, 0, 4)

		applied, n := b.ApplyHiddenPairs()
		if !applied || n < 1 {
			t.Fatalf("ApplyHiddenPairs should find and apply the hidden pair; applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(0, 1, 3) || b.HasCandidate(0, 5, 3) {
			t.Error("hidden pair should remove candidate 3 from the pair cells")
		}
		if !b.HasCandidate(0, 1, 1) || !b.HasCandidate(0, 1, 2) {
			t.Error("pair cell (0,1) should keep candidates 1 and 2")
		}
		if !b.HasCandidate(0, 5, 1) || !b.HasCandidate(0, 5, 2) {
			t.Error("pair cell (0,5) should keep candidates 1 and 2")
		}
		if !b.HasCandidate(0, 0, 3) || !b.HasCandidate(0, 0, 4) {
			t.Error("other cells in row should keep their candidates")
		}
	})

	t.Run("two_pairs_same_row", func(t *testing.T) {
		// Row 0: two disjoint hidden pairs. Digits 1,2 only in (0,0)-(0,1) with chaff 8; digit 8 also in (0,2)
		// so (1,2) is the only hidden pair on those cells. Digits 3,4 only in (0,7)-(0,8) with chaff 9; 9 also in (0,2).
		b := New()
		b.AddCandidate(0, 0, 1)
		b.AddCandidate(0, 0, 2)
		b.AddCandidate(0, 0, 8)
		b.AddCandidate(0, 1, 1)
		b.AddCandidate(0, 1, 2)
		b.AddCandidate(0, 1, 8)
		b.AddCandidate(0, 2, 8)
		b.AddCandidate(0, 2, 9)
		b.SetCell(0, 3, 5)
		b.SetCell(0, 4, 6)
		b.SetCell(0, 5, 7)
		b.SetCell(0, 6, 8)
		b.AddCandidate(0, 7, 3)
		b.AddCandidate(0, 7, 4)
		b.AddCandidate(0, 7, 9)
		b.AddCandidate(0, 8, 3)
		b.AddCandidate(0, 8, 4)
		b.AddCandidate(0, 8, 9)

		applied, n := b.ApplyHiddenPairs()
		if !applied || n < 2 {
			t.Fatalf("expected both hidden pairs in row 0 in one pass (applications >= 2), applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(0, 0, 8) || b.HasCandidate(0, 1, 8) {
			t.Error("pair (1,2) cells should lose chaff 8")
		}
		if b.HasCandidate(0, 7, 9) || b.HasCandidate(0, 8, 9) {
			t.Error("pair (3,4) cells should lose chaff 9")
		}
		if !b.HasCandidate(0, 0, 1) || !b.HasCandidate(0, 0, 2) {
			t.Error("(0,0) should keep 1 and 2")
		}
		if !b.HasCandidate(0, 7, 3) || !b.HasCandidate(0, 7, 4) {
			t.Error("(0,7) should keep 3 and 4")
		}
		if !b.HasCandidate(0, 2, 8) || !b.HasCandidate(0, 2, 9) {
			t.Error("(0,2) should still have 8 and 9")
		}
	})
}
