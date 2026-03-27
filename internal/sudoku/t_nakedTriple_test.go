package sudoku

import "testing"

func TestApplyNakedTriples(t *testing.T) {
	t.Run("naked_triple_no_elimination_when_only_three_empties", func(t *testing.T) {
		// Row 0: six givens in cols 0..5; cols 6,7,8 are the triple — nowhere else to eliminate.
		b := New()
		for c, v := range []int{4, 5, 6, 7, 8, 9} {
			b.SetCell(0, c, v)
		}

		b.AddCandidate(0, 6, 1)
		b.AddCandidate(0, 6, 2)
		b.AddCandidate(0, 7, 2)
		b.AddCandidate(0, 7, 3)
		b.AddCandidate(0, 8, 1)
		b.AddCandidate(0, 8, 3)

		if applied, n := b.ApplyNakedTriples(); applied || n != 0 {
			t.Fatalf("expected no elimination; applied=%v n=%d", applied, n)
		}
	})

	t.Run("naked_triple_removes_from_other_empty_in_row", func(t *testing.T) {
		b := New()
		for c, v := range []int{4, 5, 6, 7, 8} {
			b.SetCell(0, c, v)
		}
		b.AddCandidate(0, 5, 1)
		b.AddCandidate(0, 5, 2)
		b.AddCandidate(0, 6, 2)
		b.AddCandidate(0, 6, 3)
		b.AddCandidate(0, 7, 1)
		b.AddCandidate(0, 7, 3)
		b.AddCandidate(0, 8, 1)
		b.AddCandidate(0, 8, 9)

		applied, n := b.ApplyNakedTriples()
		if !applied || n < 1 {
			t.Fatalf("expected at least one application; applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(0, 8, 1) || b.HasCandidate(0, 8, 2) || b.HasCandidate(0, 8, 3) {
			t.Error("(0,8) should lose candidates 1, 2, and 3")
		}
		if !b.HasCandidate(0, 8, 9) {
			t.Error("(0,8) should still have 9")
		}
	})

	t.Run("naked_triple_removes_from_other_empty_in_column", func(t *testing.T) {
		b := New()
		for r, v := range []int{4, 5, 6, 7, 8} {
			b.SetCell(r, 3, v)
		}
		b.AddCandidate(5, 3, 1)
		b.AddCandidate(5, 3, 2)
		b.AddCandidate(6, 3, 2)
		b.AddCandidate(6, 3, 3)
		b.AddCandidate(7, 3, 1)
		b.AddCandidate(7, 3, 3)
		b.AddCandidate(8, 3, 1)
		b.AddCandidate(8, 3, 9)

		applied, n := b.ApplyNakedTriples()
		if !applied || n < 1 {
			t.Fatalf("expected at least one application; applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(8, 3, 1) || b.HasCandidate(8, 3, 2) || b.HasCandidate(8, 3, 3) {
			t.Error("(8,3) should lose candidates 1, 2, and 3")
		}
		if !b.HasCandidate(8, 3, 9) {
			t.Error("(8,3) should still have 9")
		}
	})

	t.Run("naked_triple_removes_from_other_empty_in_box", func(t *testing.T) {
		b := New()
		b.SetCell(0, 0, 4)
		b.SetCell(0, 1, 5)
		b.SetCell(0, 2, 6)
		b.SetCell(2, 1, 7)
		b.SetCell(2, 2, 8)

		b.AddCandidate(1, 0, 1)
		b.AddCandidate(1, 0, 2)
		b.AddCandidate(1, 1, 2)
		b.AddCandidate(1, 1, 3)
		b.AddCandidate(1, 2, 1)
		b.AddCandidate(1, 2, 3)
		b.AddCandidate(2, 0, 1)
		b.AddCandidate(2, 0, 9)

		applied, n := b.ApplyNakedTriples()
		if !applied || n < 1 {
			t.Fatalf("expected at least one application; applied=%v n=%d", applied, n)
		}
		if b.HasCandidate(2, 0, 1) || b.HasCandidate(2, 0, 2) || b.HasCandidate(2, 0, 3) {
			t.Error("(2,0) should lose candidates 1, 2, and 3")
		}
		if !b.HasCandidate(2, 0, 9) {
			t.Error("(2,0) should still have 9")
		}
	})
}
