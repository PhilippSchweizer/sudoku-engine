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
	t.Run("naked_pair_two_empty_cells", func(t *testing.T) {
		// Row 0: [1, 0, 6, 2, 8, 0, 5, 3, 9] -> 3 & 7 in 0,1 and 0,5
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
		// b.UpdateCandidates()

		// Apply naked pair: should remove 3 and 7 from other cells in row 0
		applied := b.ApplyNakedPair()
		if !applied {
			t.Fatal("ApplyNakedPair should find and apply the pair")
		}

		if !b.HasCandidate(0, 1, 3) || !b.HasCandidate(0, 1, 7) {
			t.Error("pair cell (0,1) should still have candidates 3 and 7")
		}

		if !b.HasCandidate(0, 5, 3) || !b.HasCandidate(0, 5, 7) {
			t.Error("pair cell (0,5) should still have candidates 3 and 7")
		}
	})

	t.Run("naked_pair_three_empty_cells", func(t *testing.T) {
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
		b.AddCandidate(0, 0, 3)
		b.AddCandidate(0, 0, 7)
		b.AddCandidate(0, 0, 4)

		applied := b.ApplyNakedPair()

		if !applied {
			t.Fatal("ApplyNakedPair should find and apply the pair.")
		}

		if b.HasCandidate(0, 0, 3) || b.HasCandidate(0, 0, 7) {
			t.Error("cell (0,0) should lose candidates 3 and 7 after naked pair elimination")
		}
		if !b.HasCandidate(0, 0, 4) {
			t.Error("cell (0,0) should still have candidate 4")
		}
	})
}

func TestApplyHiddenPair(t *testing.T) {
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

	applied := b.ApplyHiddenPair()
	if !applied {
		t.Fatal("ApplyHiddenPair should find and apply the hidden pair")
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
}
