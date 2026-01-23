package sudoku

import (
	"math/rand"
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if Get(b, row, col) != 0 {
				t.Errorf("Cell (%d.%d) should be 0. Got %d", row, col, Get(b, row, col))
			}
		}
	}
}

func TestGetSet(t *testing.T) {
	b := New()
	r := rand.Intn(9)
	c := rand.Intn(9)
	v := rand.Intn(9) + 1

	b = Set(b, r, c, v)

	if Get(b, r, c) != v {
		t.Errorf("Cell %d.%d should be %d. Got %d", r, c, v, Get(b, r, c))
	}
}
