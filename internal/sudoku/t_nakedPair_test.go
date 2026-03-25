package sudoku

import "testing"

func TestApplyNakedPairs(t *testing.T) {
	t.Run("naked_pair_no_elimination_when_only_two_empties", func(t *testing.T) {
		// Row 0: [1, _, 6, 2, 8, _, 5, 3, 9] — only (0,1) and (0,5) empty; they form {3,7}
		// but there is no third empty in the row, so nothing to eliminate.
		b := New()
		b.SetCell(0, 0, 1)
		b.SetCell(0, 2, 6)
		b.SetCell(0, 3, 2)
		b.SetCell(0, 4, 8)
		b.SetCell(0, 6, 5)
		b.SetCell(0, 7, 3)
		b.SetCell(0, 8, 9)

		b.AddCandidate(0, 1, 3)
		b.AddCandidate(0, 1, 7)
		b.AddCandidate(0, 5, 3)
		b.AddCandidate(0, 5, 7)

		if applied, n := b.ApplyNakedPairs(); applied || n != 0 {
			t.Fatal("ApplyNakedPairs should return false when no candidate can be removed")
		}
	})

	t.Run("naked_pair_removes_from_other_empties_in_row", func(t *testing.T) {
		// Row 0: [1, _, _, 2, 8, _, 5, 3, 9] — pair {3,7} at (0,1) and (0,5); (0,2) also wrongly has 3,7.
		b := New()
		b.SetCell(0, 0, 1)
		b.SetCell(0, 3, 2)
		b.SetCell(0, 4, 8)
		b.SetCell(0, 6, 5)
		b.SetCell(0, 7, 3)
		b.SetCell(0, 8, 9)

		b.AddCandidate(0, 1, 3)
		b.AddCandidate(0, 1, 7)
		b.AddCandidate(0, 2, 3)
		b.AddCandidate(0, 2, 6)
		b.AddCandidate(0, 2, 7)
		b.AddCandidate(0, 5, 3)
		b.AddCandidate(0, 5, 7)

		applied, n := b.ApplyNakedPairs()
		if !applied || n != 1 {
			t.Fatalf("ApplyNakedPairs should remove 3 and 7 from (0,2); applied=%v applications=%d", applied, n)
		}

		if !b.HasCandidate(0, 1, 3) || !b.HasCandidate(0, 1, 7) {
			t.Error("pair cell (0,1) should still have candidates 3 and 7")
		}
		if !b.HasCandidate(0, 5, 3) || !b.HasCandidate(0, 5, 7) {
			t.Error("pair cell (0,5) should still have candidates 3 and 7")
		}
		if b.HasCandidate(0, 2, 3) || b.HasCandidate(0, 2, 7) {
			t.Error("cell (0,2) should lose candidates 3 and 7")
		}
		if !b.HasCandidate(0, 2, 6) {
			t.Error("cell (0,2) should still have candidate 6")
		}
	})

	t.Run("two_naked_pairs_same_row_both_apply", func(t *testing.T) {
		// Row 0: fixed digits separate two naked pairs; middle cell wrongly has 1,2,3.
		b := New()
		b.SetCell(0, 2, 5)
		b.SetCell(0, 3, 6)
		b.SetCell(0, 5, 7)
		b.SetCell(0, 6, 8)
		// Pair {1,2} at cols 0,1
		b.AddCandidate(0, 0, 1)
		b.AddCandidate(0, 0, 2)
		b.AddCandidate(0, 1, 1)
		b.AddCandidate(0, 1, 2)
		// Pair {3,4} at cols 7,8
		b.AddCandidate(0, 7, 3)
		b.AddCandidate(0, 7, 4)
		b.AddCandidate(0, 8, 3)
		b.AddCandidate(0, 8, 4)
		// Middle empty (0,4): should lose 1,2 from first pair and 3 from second
		b.AddCandidate(0, 4, 1)
		b.AddCandidate(0, 4, 2)
		b.AddCandidate(0, 4, 3)

		applied, n := b.ApplyNakedPairs()
		if !applied || n < 2 {
			t.Fatalf("expected both row naked pairs to apply (applications >= 2), got applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(0, 4, 1) || b.HasCandidate(0, 4, 2) || b.HasCandidate(0, 4, 3) {
			t.Error("(0,4) should lose 1, 2, and 3")
		}
	})
}
