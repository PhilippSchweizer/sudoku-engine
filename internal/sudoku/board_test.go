package sudoku

import (
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	for row := range 9 {
		for col := range 9 {
			if b.Cell(row, col) != 0 {
				t.Errorf("Cell (%d.%d) should be 0. Got %d", row, col, b.Cell(row, col))
			}
		}
	}
}

func TestGetSet(t *testing.T) {
	b := New()
	b.SetCell(7, 7, 7)
	if b.Cell(7, 7) != 7 {
		t.Error("Cell should have correct value after setting.")
	}
}

func TestSetCellPanicsOnInvalidValue(t *testing.T) {
	panicIfNoPanic := func(t *testing.T, val int) {
		t.Helper()
		if r := recover(); r == nil {
			t.Errorf("SetCell(%d) should panic", val)
		}
	}

	t.Run("value_too_high", func(t *testing.T) {
		defer panicIfNoPanic(t, 10)
		b := New()
		b.SetCell(0, 0, 10)
	})

	t.Run("value_negative", func(t *testing.T) {
		defer panicIfNoPanic(t, -1)
		b := New()
		b.SetCell(0, 0, -1)
	})
}

func TestAddRemoveCandidate(t *testing.T) {
	t.Run("add_candidate", func(t *testing.T) {
		b := New()
		b.AddCandidate(0, 1, 5)
		if !b.HasCandidate(0, 1, 5) {
			t.Error("Cell 0.1 should have 5 as candidate.")
		}
	})

	t.Run("add_candidate_twice", func(t *testing.T) {
		b := New()
		b.AddCandidate(0, 1, 5)
		b.AddCandidate(0, 1, 5)
		if !b.HasCandidate(0, 1, 5) {
			t.Error("Cell 0.1 should have 5 as candidate.")
		}
	})

	t.Run("remove_candidate", func(t *testing.T) {
		b := New()
		b.UpdateCandidates()
		b.RemoveCandidate(1, 2, 6)
		if b.HasCandidate(1, 2, 6) {
			t.Error("Cell 1.2 should not have 6 as candidate after removal.")
		}
	})

	t.Run("remove_nonexistent_candidate", func(t *testing.T) {
		b := New()
		b.RemoveCandidate(0, 0, 5)
		if b.HasCandidate(0, 0, 5) {
			t.Error("Should not have candidate after removal.")
		}
	})

	t.Run("get_candidates", func(t *testing.T) {
		b := New()
		b.SetCell(8, 8, 1)
		b.SetCell(6, 6, 6)
		b.SetCell(8, 4, 4)
		b.UpdateCandidates()
		if b.GetCandidates(8, 7) != 0b111010110 {
			t.Error("Should return only candidates that don't contradict set cells.")
		}
	})
}

func TestUpdateCandidates(t *testing.T) {
	t.Run("empty_board", func(t *testing.T) {
		b := New()
		b.UpdateCandidates()
		if b.GetCandidates(2, 2) != 0b111111111 {
			t.Error("Every empty cell should have every candidate after updating candidates.")
		}
	})

	t.Run("filled_cell", func(t *testing.T) {
		b := New()
		b.SetCell(5, 5, 8)
		b.UpdateCandidates()
		if b.GetCandidates(5, 5) != 0 {
			t.Errorf("Non empty cell should not have any candidates.")
		}
	})

	t.Run("row_conflict", func(t *testing.T) {
		b := New()
		b.SetCell(1, 1, 1)
		b.UpdateCandidates()
		if b.HasCandidate(1, 2, 1) {
			t.Errorf("Cell 2 of row 1 should not have 1 as candidate.")
		}
	})

	t.Run("col_conflict", func(t *testing.T) {
		b := New()
		b.SetCell(2, 2, 2)
		b.UpdateCandidates()
		if b.HasCandidate(3, 2, 2) {
			t.Error("Cell 3 of col 2 should not have 2 as candidate.")
		}
	})

	t.Run("box_conflict", func(t *testing.T) {
		b := New()
		b.SetCell(4, 4, 4)
		b.UpdateCandidates()
		if b.HasCandidate(3, 3, 4) {
			t.Error("Cell 3.3 should not have 4 as candidate.")
		}
	})
}

func TestSetCellWithCandidateUpdate(t *testing.T) {
	b := New()
	b.SetCellWithCandidateUpdate(4, 4, 4)
	if b.HasCandidate(3, 4, 4) || b.HasCandidate(4, 5, 4) || b.HasCandidate(3, 3, 4) {
		t.Error("Cell candidates in same row, column and box should have been updated.")
	}
}
