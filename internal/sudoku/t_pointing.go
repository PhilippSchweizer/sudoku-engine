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
/* func (b Board) pointing() (found, alongRow bool, box, line, val int) {
	for boxIdx := range 9 {
		for v := 1; v <= 9; v++ {
			if ok, ar, ln := b.pointingInBox(boxIdx, v); ok {
				return true, ar, boxIdx, ln, v
			}
		}
	}

	return false, false, -1, -1, 0
} */

// ApplyPointings performs pointing eliminations: when every candidate for a digit in a box lies on
// one row or column inside that box, that digit is removed from the rest of that line outside the box.
// Repeats until a full pass over all boxes and digits makes no change.
func (b *Board) ApplyPointings() (applied bool, applications int) {
	overall := false
	for {
		changed := false
		passApps := 0

		for box := range 9 {
			br := (box / 3) * 3
			bc := (box % 3) * 3
			for v := 1; v <= 9; v++ {
				ok, alongRow, line := b.pointingInBox(box, v)
				if !ok {
					continue
				}
				boxDigitChanged := false
				if alongRow {
					r := line
					for c := range 9 {
						if c >= bc && c < bc+3 {
							continue
						}
						if b.removeCandidateIfPresent(r, c, v) {
							boxDigitChanged = true
							changed = true
						}
					}
				} else {
					c := line
					for r := range 9 {
						if r >= br && r < br+3 {
							continue
						}
						if b.removeCandidateIfPresent(r, c, v) {
							boxDigitChanged = true
							changed = true
						}
					}
				}
				if boxDigitChanged {
					passApps++
				}
			}
		}

		applications += passApps
		if !changed {
			break
		}
		overall = true
	}
	return overall, applications
}
