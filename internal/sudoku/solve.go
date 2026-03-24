package sudoku

func Solve(b Board) (Board, bool) {
	if !b.IsValid() {
		return b, false
	}
	return solve(b)
}

func solve(b Board) (Board, bool) {
	if b.IsSolved() {
		return b, true
	}

	row, col, mask := pickEmptyMRV(b)
	if row < 0 || mask == 0 {
		return b, false
	}

	for m := mask; m != 0; m &= m - 1 {
		v := bits.TrailingZeros16(m) + 1
		newBoard := b
		newBoard.SetCell(row, col, v)
		solvedBoard, solved := solve(newBoard)
		if solved {
			return solvedBoard, true
		}
	}

	return b, false
}

func CountSolutions(b Board) int {
	if !b.IsValid() {
		return 0
	}
	return countSolutions(b)
}

func countSolutions(b Board) int {
	if b.IsSolved() {
		return 1
	}

	row, col, mask := pickEmptyMRV(b)
	if row < 0 || mask == 0 {
		return 0
	}

	count := 0
	for m := mask; m != 0; m &= m - 1 {
		v := bits.TrailingZeros16(m) + 1
		newBoard := b
		newBoard.SetCell(row, col, v)
		subCount := countSolutions(newBoard)
		count += subCount
		if count >= 2 {
			return count
		}
	}
	return count
}

// validDigitMask returns a bitmask of digits 1–9 that can be placed at (row,col)
// without conflicting with row, column, or box. Empty if the cell is already filled.
func validDigitMask(b Board, row, col int) uint16 {
	if b.Cell(row, col) != 0 {
		return 0
	}
	const allDigits uint16 = (1 << 9) - 1
	var used uint16
	for i := range 9 {
		if v := b.Cell(row, i); v != 0 {
			used |= bitFor(v)
		}
		if v := b.Cell(i, col); v != 0 {
			used |= bitFor(v)
		}
	}
	br := (row / 3) * 3
	bc := (col / 3) * 3
	for i := range 9 {
		if v := b.Cell(br+i/3, bc+i%3); v != 0 {
			used |= bitFor(v)
		}
	}
	return allDigits &^ used
}

// pickEmptyMRV chooses an empty cell with the fewest candidate digits (MRV heuristic).
func pickEmptyMRV(b Board) (row, col int, mask uint16) {
	best := 10
	row, col = -1, -1
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) != 0 {
				continue
			}
			m := validDigitMask(b, r, c)
			n := bits.OnesCount16(m)
			if n < best {
				best = n
				row, col = r, c
				mask = m
				if n <= 1 {
					return row, col, mask
				}
			}
		}
	}
	return row, col, mask
}

func (b Board) nextEmpty() (row, col int) {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) == 0 {
				return r, c
			}
		}
	}
	return -1, -1 // no empty cell found
}
