package sudoku

import "testing"

func TestBoardValid(t *testing.T) {
	t.Run("empty_board", func(t *testing.T) {
		if !BoardValid(New()) {
			t.Error("empty board should be valid")
		}
	})
	t.Run("duplicate_in_row", func(t *testing.T) {
		b := New()
		r := 0
		b.SetCell(r, 0, 1)
		b.SetCell(r, 3, 1)

		if BoardValid(b) != false {
			t.Errorf("Duplicate in row %d not detected", (r + 1))
		}
	})

	t.Run("duplicate_in_column", func(t *testing.T) {
		b := New()
		c := 8
		b.SetCell(3, c, 5)
		b.SetCell(5, c, 5)

		if BoardValid(b) != false {
			t.Errorf("Duplicate in column %d not detected", (c + 1))
		}
	})

	t.Run("duplicate_in_box", func(t *testing.T) {
		b := New()
		b.SetCell(0, 0, 1)
		b.SetCell(2, 2, 1)

		if BoardValid(b) != false {
			t.Error("Duplicate in box not detected")
		}
	})
}

func TestBoardSolved(t *testing.T) {
	t.Run("solved_board", func(t *testing.T) {
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
		if !BoardSolved(b) {
			t.Error("Valid full board should be evaluated as solved.")
		}
	})

	t.Run("incomplete_has_zero", func(t *testing.T) {
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
				{3, 4, 5, 2, 8, 6, 1, 7, 0},
			},
		}
		if BoardSolved(b) {
			t.Error("Board with empty cell should not be evaluated as solved.")
		}
	})

	t.Run("full_but_invalid_duplicate", func(t *testing.T) {
		b := Board{
			Cells: [9][9]int{
				{5, 5, 4, 6, 7, 8, 9, 1, 2},
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
		if BoardSolved(b) {
			t.Error("Full board with duplicate in row should not be evaluated as solved.")
		}
	})

	t.Run("empty_board", func(t *testing.T) {
		if BoardSolved(New()) {
			t.Error("Empty board should not be evaluated as solved.")
		}
	})
}
