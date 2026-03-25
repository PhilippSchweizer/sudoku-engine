package sudoku

// applyHiddenPairsInUnit finds every hidden pair in one unit (getCell maps slot 0..8 to coordinates)
// and strips all other candidates from the two pair cells.
func (b *Board) applyHiddenPairsInUnit(getCell func(i int) (row, col int)) (pairApplications int, changed bool) {
	type coord struct{ r, c int }
	var positions [9][]coord
	for p := range 9 {
		row, col := getCell(p)
		if b.Cell(row, col) != 0 {
			continue
		}
		for v := 1; v <= 9; v++ {
			if b.HasCandidate(row, col, v) {
				positions[v-1] = append(positions[v-1], coord{row, col})
			}
		}
	}
	for v1 := range 9 {
		if len(positions[v1]) != 2 {
			continue
		}
		for v2 := v1 + 1; v2 < 9; v2++ {
			if len(positions[v2]) != 2 {
				continue
			}
			if positions[v1][0] != positions[v2][0] || positions[v1][1] != positions[v2][1] {
				continue
			}
			val1, val2 := v1+1, v2+1
			c1 := positions[v1][0]
			c2 := positions[v1][1]
			pairChanged := false
			for v := 1; v <= 9; v++ {
				if v == val1 || v == val2 {
					continue
				}
				rm1 := b.removeCandidateIfPresent(c1.r, c1.c, v)
				rm2 := b.removeCandidateIfPresent(c2.r, c2.c, v)
				if rm1 || rm2 {
					pairChanged = true
					changed = true
				}
			}
			if pairChanged {
				pairApplications++
			}
		}
	}
	return pairApplications, changed
}

func (b *Board) ApplyHiddenPairs() (applied bool, applications int) {
	overall := false
	for {
		changed := false
		passApps := 0

		for r := range 9 {
			pa, ch := b.applyHiddenPairsInUnit(func(i int) (int, int) { return r, i })
			passApps += pa
			if ch {
				changed = true
			}
		}
		for c := range 9 {
			pa, ch := b.applyHiddenPairsInUnit(func(i int) (int, int) { return i, c })
			passApps += pa
			if ch {
				changed = true
			}
		}
		for box := range 9 {
			br := (box / 3) * 3
			bc := (box % 3) * 3
			pa, ch := b.applyHiddenPairsInUnit(func(i int) (int, int) { return br + i/3, bc + i%3 })
			passApps += pa
			if ch {
				changed = true
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
