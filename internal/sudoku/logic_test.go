package sudoku

import "testing"

func cellsEqual(a, b Board) bool {
	for r := range 9 {
		for c := range 9 {
			if a.Cell(r, c) != b.Cell(r, c) {
				return false
			}
		}
	}
	return true
}

// solvedTestGrid is a complete valid Sudoku (same as in solve_test.go).
func solvedTestGrid() Board {
	return Board{
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
}

func TestSolveByLogic_invalidBoard(t *testing.T) {
	t.Parallel()
	b := solvedTestGrid()
	b.SetCell(8, 7, 1) // duplicate with (8,6); breaks validity

	got := SolveByLogic(b)

	if got.valid {
		t.Fatal("expected valid=false for invalid board")
	}
	if !cellsEqual(got.Start, b) {
		t.Error("Start should preserve input cells")
	}
}

func TestSolveByLogic_alreadySolved(t *testing.T) {
	t.Parallel()
	b := solvedTestGrid()

	got := SolveByLogic(b)

	if !got.valid {
		t.Fatal("expected valid=true")
	}
	if got.solvedByLogic {
		t.Error("already-complete board should not count as solved by logic steps")
	}
	if !cellsEqual(got.Start, b) {
		t.Error("Start should match input")
	}
	if !got.Progress.IsSolved() || !got.Progress.IsValid() {
		t.Error("Progress should be complete and valid")
	}
	if !cellsEqual(got.Progress, b) {
		t.Error("Progress cells should match solved input")
	}
	if got.rounds != 0 {
		t.Errorf("rounds: want 0, got %d", got.rounds)
	}
	if got.nakedSingleCount != 0 || got.hiddenSingleCount != 0 || got.nakedPairCount != 0 || got.hiddenPairCount != 0 {
		t.Error("no techniques should run on already-solved board")
	}
	if got.MaxTechnique != TechniqueNakedSingle {
		t.Errorf("MaxTechnique: want TechniqueNakedSingle, got %v", got.MaxTechnique)
	}
}

func TestSolveByLogic_singleNakedSingleCompletes(t *testing.T) {
	t.Parallel()
	b := solvedTestGrid()
	b.SetCell(8, 8, 0) // only 9 is valid here

	got := SolveByLogic(b)

	if !got.valid {
		t.Fatal("expected valid=true")
	}
	if !got.solvedByLogic {
		t.Fatal("expected puzzle completed by naked single")
	}
	if !got.Progress.IsSolved() || !got.Progress.IsValid() {
		t.Fatal("Progress should be fully solved")
	}
	if got.nakedSingleCount < 1 {
		t.Errorf("expected at least one naked single, got %d", got.nakedSingleCount)
	}
	if got.MaxTechnique != TechniqueNakedSingle {
		t.Errorf("MaxTechnique: want TechniqueNakedSingle, got %v", got.MaxTechnique)
	}
	if !cellsEqual(got.Start, b) {
		t.Error("Start should still reflect original puzzle (one empty cell)")
	}
}

func TestSolveByLogic_stuckEmptyBoard(t *testing.T) {
	t.Parallel()
	// Empty grid: valid, but no naked/hidden single or pair applies (every empty still has 9 candidates).
	got := SolveByLogic(New())

	if !got.valid {
		t.Fatal("expected valid=true")
	}
	if got.solvedByLogic {
		t.Fatal("expected logic solver to make no progress on empty board")
	}
	if got.Progress.IsSolved() {
		t.Fatal("Progress should not be a complete grid")
	}
	if got.rounds != 0 {
		t.Errorf("rounds: want 0, got %d", got.rounds)
	}
	if got.nakedSingleCount != 0 || got.hiddenSingleCount != 0 || got.nakedPairCount != 0 || got.hiddenPairCount != 0 {
		t.Error("no technique should apply on empty board")
	}
}
