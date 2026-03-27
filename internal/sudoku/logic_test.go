package sudoku

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

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

// writeLogicSolveBoardDump writes stats plus Start and Progress as ASCII grids with pencil marks.
// Start is shown with candidates recomputed from givens (Start may not store pencil marks before solving).
func writeLogicSolveBoardDump(path string, got LogicSolveResult) error {
	var b strings.Builder
	b.WriteString("LogicSolveResult\n")
	b.WriteString(strings.Repeat("=", 48))
	b.WriteByte('\n')
	b.WriteString("valid: ")
	b.WriteString(strconv.FormatBool(got.valid))
	b.WriteString(", solvedByLogic: ")
	b.WriteString(strconv.FormatBool(got.solvedByLogic))
	b.WriteString(", rounds: ")
	b.WriteString(strconv.Itoa(got.rounds))
	b.WriteString(", maxTechnique: ")
	b.WriteString(strconv.Itoa(int(got.MaxTechnique)))
	b.WriteString("\n")
	b.WriteString("nakedSingle=")
	b.WriteString(strconv.Itoa(got.nakedSingleCount))
	b.WriteString(" hiddenSingle=")
	b.WriteString(strconv.Itoa(got.hiddenSingleCount))
	b.WriteString(" nakedPair=")
	b.WriteString(strconv.Itoa(got.nakedPairCount))
	b.WriteString(" hiddenPair=")
	b.WriteString(strconv.Itoa(got.hiddenPairCount))
	b.WriteString(" pointing=")
	b.WriteString(strconv.Itoa(got.pointingCount))
	b.WriteString(" nakedTriple=")
	b.WriteString(strconv.Itoa(got.nakedTripleCount))
	b.WriteString("\n\n")

	start := got.Start
	start.UpdateCandidates()
	b.WriteString("Start (givens + pencil marks)\n")
	b.WriteString(strings.Repeat("-", 48))
	b.WriteByte('\n')
	b.WriteString(start.FormatWithPencilMarks())
	b.WriteByte('\n')

	b.WriteString("Progress (after SolveByLogic)\n")
	b.WriteString(strings.Repeat("-", 48))
	b.WriteByte('\n')
	b.WriteString(got.Progress.FormatWithPencilMarks())
	b.WriteByte('\n')

	return os.WriteFile(path, []byte(b.String()), 0o644)
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
	if got.nakedSingleCount != 0 || got.hiddenSingleCount != 0 || got.nakedPairCount != 0 || got.hiddenPairCount != 0 || got.pointingCount != 0 || got.nakedTripleCount != 0 {
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
	if got.nakedSingleCount != 0 || got.hiddenSingleCount != 0 || got.nakedPairCount != 0 || got.hiddenPairCount != 0 || got.pointingCount != 0 || got.nakedTripleCount != 0 {
		t.Error("no technique should apply on empty board")
	}
}

// TestSolveByLogic_generatedPuzzle exercises the logic solver on a randomly generated puzzle.
// Run with -v to see LogicSolveResult fields in the test log.
//
// If the logged numbers never change between invocations, you are usually seeing cached output:
// Go re-runs passing tests only when inputs change; otherwise it replays the prior log. Use
//   go test ./internal/sudoku/... -run TestSolveByLogic_generatedPuzzle -v -count=1
// to disable caching and generate a fresh puzzle each time.
func TestSolveByLogic_generatedPuzzle(t *testing.T) {
	t.Parallel()
	puzzle, solution := GeneratePuzzle()
	got := SolveByLogic(puzzle)

	t.Logf("LogicSolveResult: valid=%v solvedByLogic=%v rounds=%d maxTechnique=%v\n"+
		"  counts: nakedSingle=%d hiddenSingle=%d nakedPair=%d hiddenPair=%d pointing=%d nakedTriple=%d\n"+
		"  progress matches generator solution: %v",
		got.valid, got.solvedByLogic, got.rounds, got.MaxTechnique,
		got.nakedSingleCount, got.hiddenSingleCount, got.nakedPairCount, got.hiddenPairCount, got.pointingCount, got.nakedTripleCount,
		got.solvedByLogic && cellsEqual(got.Progress, solution),
	)

	// Written under package testdata/ (cwd is this package when `go test` runs) so you can open it after the test.
	dumpPath := filepath.Join("testdata", "logic_solve_boards.txt")
	if err := os.MkdirAll(filepath.Dir(dumpPath), 0o755); err != nil {
		t.Fatalf("mkdir testdata: %v", err)
	}
	if err := writeLogicSolveBoardDump(dumpPath, got); err != nil {
		t.Fatalf("write board dump: %v", err)
	}
	abs, _ := filepath.Abs(dumpPath)
	if abs != "" {
		dumpPath = abs
	}
	t.Logf("board dump (Start + Progress with pencil marks): %s", dumpPath)

	if !got.valid {
		t.Fatal("generated puzzle should be valid")
	}
	if !cellsEqual(got.Start, puzzle) {
		t.Error("Start should preserve input puzzle")
	}
}
