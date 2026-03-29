package sudoku

import "testing"

func TestApplyPointings_rowLocked(t *testing.T) {
	// Box 0: digit 7 only on row 0 inside the box — (0,0) and (0,1). Same row on (0,5) outside box
	// should lose 7.
	b := New()
	b.UpdateCandidates()
	for _, rc := range [][2]int{{0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}} {
		b.stripCandidate(rc[0], rc[1], 7)
	}

	applied, n := b.ApplyPointings()
	if !applied || n < 1 {
		t.Fatalf("ApplyPointings should eliminate 7 from row 0 outside box 0; applied=%v n=%d", applied, n)
	}
	for c := 3; c < 9; c++ {
		if b.HasCandidate(0, c, 7) {
			t.Errorf("cell (0,%d) should not have candidate 7 after pointing", c)
		}
	}
	if !b.HasCandidate(0, 0, 7) || !b.HasCandidate(0, 1, 7) {
		t.Error("pair cells in box should keep 7")
	}
}

func TestApplyPointings_columnLocked(t *testing.T) {
	// Box 0: digit 4 only in column 0 — (0,0),(1,0),(2,0). Strip 4 from (3,0)..(8,0).
	b := New()
	b.UpdateCandidates()
	for _, rc := range [][2]int{{0, 1}, {0, 2}, {1, 1}, {1, 2}, {2, 1}, {2, 2}} {
		b.stripCandidate(rc[0], rc[1], 4)
	}

	applied, n := b.ApplyPointings()
	if !applied || n < 1 {
		t.Fatalf("ApplyPointings should eliminate 4 from column 0 outside box 0; applied=%v n=%d", applied, n)
	}
	for r := 3; r < 9; r++ {
		if b.HasCandidate(r, 0, 4) {
			t.Errorf("cell (%d,0) should not have candidate 4 after pointing", r)
		}
	}
}
