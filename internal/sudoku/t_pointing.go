package sudoku

// pointingInBox reports whether all candidates for digit in box lie on one row or column inside the box.
// alongRow true means line is a row index; false means column index.
func (b Board) pointingInBox(box, digit int) (found, alongRow bool, line int) {
	br := (box / 3) * 3
	bc := (box % 3) * 3

	var hits int
	var r0, c0 int
	rowLocked, colLocked := true, true

	for i := range 9 {
		r, c := br+i/3, bc+i%3
		if b.Cell(r, c) != 0 {
			continue
		}
		if !b.HasCandidate(r, c, digit) {
			continue
		}
		if hits == 0 {
			r0, c0 = r, c
		} else {
			if r != r0 {
				rowLocked = false
			}
			if c != c0 {
				colLocked = false
			}
		}
		hits++
	}

	if hits == 0 {
		return false, false, -1
	}
	if rowLocked {
		return true, true, r0
	}
	if colLocked {
		return true, false, c0
	}
	return false, false, -1
}

// pointing finds the first box where digit is locked to one row or column within the box.
func (b Board) pointing() (found, alongRow bool, box, line, val int) {
	for boxIdx := range 9 {
		for v := 1; v <= 9; v++ {
			if ok, ar, ln := b.pointingInBox(boxIdx, v); ok {
				return true, ar, boxIdx, ln, v
			}
		}
	}

	return false, false, -1, -1, 0
}
