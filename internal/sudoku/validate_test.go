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
		b = Set(b, r, 0, 1)
		b = Set(b, r, 3, 1)

		if BoardValid(b) != false {
			t.Errorf("Duplicate in row %d not detected", (r + 1))
		}
	})

	t.Run("duplicate_in_column", func(t *testing.T) {
		b := New()
		c := 8
		b = Set(b, 3, c, 5)
		b = Set(b, 5, c, 5)

		if BoardValid(b) != false {
			t.Errorf("Duplicate in column %d not detected", (c + 1))
		}
	})

	t.Run("duplicate_in_box", func(t *testing.T) {
		b := New()
		b = Set(b, 0, 0, 1)
		b = Set(b, 2, 2, 1)

		if BoardValid(b) != false {
			t.Error("Duplicate in box not detected")
		}
	})
}
