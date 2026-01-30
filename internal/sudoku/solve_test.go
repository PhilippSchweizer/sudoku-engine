package sudoku

import "testing"

func TestSolve(t *testing.T) {
	t.Run("solve_empty", func(t *testing.T) {
		result, solved := Solve(New())
		if !result.IsSolved() || !result.IsValid() {
			t.Error("Board solved by Solve() should be complete and valid.")
		}

		if !solved {
			t.Error("should be able to solve empty board")
		}
	})

	t.Run("already_solved", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{5, 3, 4, 6, 7, 8, 9, 1, 2},
				{6, 7, 2, 1, 9, 5, 3, 4, 8},
				{1, 9, 8, 3, 4, 2, 5, 6, 7},
				{8, 5, 9, 7, 6, 1, 4, 2, 3},
				{4, 2, 6, 8, 5, 3, 7, 9, 1},
				{7, 1, 3, 9, 2, 4, 8, 5, 6},
				{9, 6, 1, 5, 3, 7, 2, 8, 4},
				{2, 8, 7, 4, 1, 9, 6, 3, 5},
				{3, 4, 5, 2, 8, 6, 1, 7, 9},
			},
		}
		_, solved := Solve(b)
		if !solved {
			t.Error("should recognize alread solved board as solved")
		}
	})

	t.Run("invalid_board", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{5, 3, 4, 6, 7, 8, 9, 1, 2},
				{6, 7, 2, 1, 9, 5, 3, 4, 8},
				{1, 9, 8, 3, 4, 2, 5, 6, 7},
				{8, 5, 9, 7, 6, 1, 4, 2, 3},
				{4, 2, 6, 8, 5, 3, 7, 9, 1},
				{7, 1, 3, 9, 2, 4, 8, 5, 6},
				{9, 6, 1, 5, 3, 7, 2, 8, 4},
				{2, 8, 7, 4, 1, 9, 6, 3, 5},
				{3, 4, 5, 2, 8, 1, 1, 7, 9},
			},
		}
		_, solved := Solve(b)
		if solved {
			t.Error("Invalid board should be identified as such by Solve().")
		}
	})

	t.Run("difficult", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{0, 1, 6, 0, 8, 0, 5, 3, 0},
				{0, 0, 4, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 9, 0, 0, 0, 8, 0},
				{4, 0, 0, 0, 1, 0, 0, 0, 0},
				{9, 0, 0, 0, 0, 0, 3, 0, 0},
				{0, 3, 1, 7, 0, 0, 0, 0, 6},
				{0, 0, 0, 0, 0, 2, 0, 0, 7},
				{3, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 5, 8, 0, 9, 0, 6, 0, 0},
			},
		}
		_, solved := Solve(b)
		if !solved {
			t.Error("should be able to solve difficult sudoku puzzle")
		}
	})
}

func TestCountSolutions(t *testing.T) {
	t.Run("count_empty_board_solutions", func(t *testing.T) {
		if CountSolutions(New()) <= 1 {
			t.Error("should find more than one solution for empty board")
		}
	})

	t.Run("complete_valid_board", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{5, 3, 4, 6, 7, 8, 9, 1, 2},
				{6, 7, 2, 1, 9, 5, 3, 4, 8},
				{1, 9, 8, 3, 4, 2, 5, 6, 7},
				{8, 5, 9, 7, 6, 1, 4, 2, 3},
				{4, 2, 6, 8, 5, 3, 7, 9, 1},
				{7, 1, 3, 9, 2, 4, 8, 5, 6},
				{9, 6, 1, 5, 3, 7, 2, 8, 4},
				{2, 8, 7, 4, 1, 9, 6, 3, 5},
				{3, 4, 5, 2, 8, 6, 1, 7, 9},
			},
		}
		if CountSolutions(b) != 1 {
			t.Error("should return 1 solution for complete & valid board")
		}
	})

	t.Run("multiple_solutions", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{1, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 1, 0, 0, 0, 0, 0},
				{0, 0, 2, 0, 0, 0, 1, 0, 0},
				{0, 1, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1, 2, 0, 0, 0},
				{0, 2, 0, 0, 0, 0, 0, 1, 0},
				{0, 0, 1, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 2, 0, 1, 0, 0, 0},
				{2, 0, 0, 0, 0, 0, 0, 0, 1},
			},
		}
		if CountSolutions(b) < 2 {
			t.Error("should return more than one solution")
		}
	})
}

func TestApplyNakedPair(t *testing.T) {
	// Set up a board with a naked pair in row 0: cells (0,1) and (0,5) only have candidates {3, 7}.
	// Fill rest of row 0 so the row is valid; leave (0,1) and (0,5) empty with only {3,7}.
	b := New()
	b.SetCell(0, 0, 1)
	b.SetCell(0, 2, 6)
	b.SetCell(0, 3, 2)
	b.SetCell(0, 4, 8)
	b.SetCell(0, 5, 0) // empty, will set candidates below
	b.SetCell(0, 6, 5)
	b.SetCell(0, 7, 3)
	b.SetCell(0, 8, 9)
	// Row 0: [1, 0, 6, 2, 8, 0, 5, 3, 9] -> only (0,1) and (0,5) empty
	// Set (0,1) and (0,5) to only candidates {3, 7} (no UpdateCandidates so no other pair)
	b.AddCandidate(0, 1, 3)
	b.AddCandidate(0, 1, 7)
	b.AddCandidate(0, 5, 3)
	b.AddCandidate(0, 5, 7)
	// Apply naked pair: should remove 3 and 7 from other cells in row 0 (none empty except pair)
	applied := b.ApplyNakedPair()
	if !applied {
		t.Fatal("ApplyNakedPair should find and apply the pair")
	}
	// Pair cells still have 3 and 7
	if !b.HasCandidate(0, 1, 3) || !b.HasCandidate(0, 1, 7) {
		t.Error("pair cell (0,1) should still have candidates 3 and 7")
	}
	if !b.HasCandidate(0, 5, 3) || !b.HasCandidate(0, 5, 7) {
		t.Error("pair cell (0,5) should still have candidates 3 and 7")
	}

	// Case 2: row with one more empty cell that has 3,7 — should lose 3,7 after elimination
	b2 := New()
	b2.SetCell(0, 0, 1)
	b2.SetCell(0, 2, 6)
	b2.SetCell(0, 3, 2)
	b2.SetCell(0, 4, 8)
	b2.SetCell(0, 6, 5)
	b2.SetCell(0, 7, 3)
	b2.SetCell(0, 8, 9)
	// (0,1), (0,5) = pair with {3,7}; (0,0) empty with 3,7 so we can give it 3,7 and expect removal
	b2.AddCandidate(0, 1, 3)
	b2.AddCandidate(0, 1, 7)
	b2.AddCandidate(0, 5, 3)
	b2.AddCandidate(0, 5, 7)
	b2.AddCandidate(0, 0, 3)
	b2.AddCandidate(0, 0, 7)
	b2.AddCandidate(0, 0, 4)
	applied2 := b2.ApplyNakedPair()
	if !applied2 {
		t.Fatal("ApplyNakedPair should find and apply the pair (case 2)")
	}
	if b2.HasCandidate(0, 0, 3) || b2.HasCandidate(0, 0, 7) {
		t.Error("cell (0,0) should lose candidates 3 and 7 after naked pair elimination")
	}
	if !b2.HasCandidate(0, 0, 4) {
		t.Error("cell (0,0) should still have candidate 4")
	}
}
