package sudoku

import "testing"

func TestSolve(t *testing.T) {
	t.Run("solve_empty", func(t *testing.T) {
		_, solved := Solve(New())
		if !solved {
			t.Error("should be able to solve empty board")
		}
	})
}
