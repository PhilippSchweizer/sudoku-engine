package sudoku

import "testing"

func TestGenerate(t *testing.T) {
	// t.Run("")
	b := Generate()
	if !b.IsSolved() {
		t.Error("Generate() should always return a full board.")
	}
	if !b.IsValid() {
		t.Error("Generate() should always return a valid board.")
	}
}

func TestGeneratePuzzle(t *testing.T) {
	t.Run("two_boards_valid", func(t *testing.T) {
		puzzle, solution := GeneratePuzzle()
		if !puzzle.IsValid() && !solution.IsValid() {
			t.Error("Returned boards should be valid.")
		}
	})

	t.Run("solution_is_solved", func(t *testing.T) {
		_, solution := GeneratePuzzle()
		if !solution.IsSolved() {
			t.Error("Solution board should be complete.")
		}
	})

	t.Run("puzzle_is_subgrid_of_solution", func(t *testing.T) {
		puzzle, solution := GeneratePuzzle()
		for r := range 9 {
			for c := range 9 {
				if puzzle.Cell(r, c) != solution.Cell(r, c) && puzzle.Cell(r, c) != 0 {
					t.Error("Puzzle board should be a subgrid of solution.")
				}
			}
		}
	})
}
