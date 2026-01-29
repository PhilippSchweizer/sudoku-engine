package sudoku

import "testing"

func TestSolve(t *testing.T) {
	t.Run("solve_empty", func(t *testing.T) {
		result, solved := Solve(New())
		if !BoardSolved(result) || !BoardValid(result) {
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
