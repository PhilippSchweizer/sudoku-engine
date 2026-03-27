package sudoku

import "math/bits"

func (b *Board) applyNakedTripleInUnit(
	getCell func(i int) (row, col int),
	eliminate func(b *Board, r1, c1, r2, c2, r3, c3, v1, v2, v3 int) (tripleChanged bool),
) (tripleApplications int, changed bool) {
	for i := range 9 {
		row, col := getCell(i)
		if b.Cell(row, col) != 0 {
			continue
		}
		mask_i := b.GetCandidates(row, col)
		if bits.OnesCount16(mask_i) != 2 && bits.OnesCount16(mask_i) != 3 {
			continue
		}

		for j := i + 1; j < 9; j++ {
			row2, col2 := getCell(j)
			if b.Cell(row2, col2) != 0 {
				continue
			}
			mask_j := b.GetCandidates(row2, col2)
			if bits.OnesCount16(mask_j) != 2 && bits.OnesCount16(mask_j) != 3 {
				continue
			}

			for k := j + 1; k < 9; k++ {
				row3, col3 := getCell(k)
				if b.Cell(row3, col3) != 0 {
					continue
				}
				mask_k := b.GetCandidates(row3, col3)
				if bits.OnesCount16(mask_k) != 2 && bits.OnesCount16(mask_k) != 3 {
					continue
				}
				U := mask_i | mask_j | mask_k
				if bits.OnesCount16(U) != 3 {
					continue
				}
				v1 := bits.TrailingZeros16(U) + 1
				m2 := U & ^(uint16(1) << (v1 - 1))
				v2 := bits.TrailingZeros16(m2) + 1
				m3 := m2 & ^(uint16(1) << (v2 - 1))
				v3 := bits.TrailingZeros16(m3) + 1

				if eliminate(b, row, col, row2, col2, row3, col3, v1, v2, v3) {
					tripleApplications++
					changed = true
				}
			}
		}
	}
	return tripleApplications, changed
}

func (b *Board) ApplyNakedTriples() (applied bool, applications int) {
	overall := false
	for {
		changed := false
		passApps := 0

		for r := range 9 {
			pa, ch := b.applyNakedTripleInUnit(
				func(i int) (int, int) { return r, i },
				func(b *Board, r1, c1, r2, c2, r3, c3, v1, v2, v3 int) bool {
					tripleChanged := false
					for c := range 9 {
						if c == c1 || c == c2 || c == c3 {
							continue
						}
						rm1 := b.removeCandidateIfPresent(r1, c, v1)
						rm2 := b.removeCandidateIfPresent(r1, c, v2)
						rm3 := b.removeCandidateIfPresent(r1, c, v3)
						if rm1 || rm2 || rm3 {
							tripleChanged = true
						}
					}
					return tripleChanged
				},
			)
			passApps += pa
			if ch {
				changed = true
			}
		}

		for c := range 9 {
			pa, ch := b.applyNakedTripleInUnit(
				func(i int) (int, int) { return i, c },
				func(b *Board, r1, c1, r2, c2, r3, c3, v1, v2, v3 int) bool {
					tripleChanged := false
					for r := range 9 {
						if r == r1 || r == r2 || r == r3 {
							continue
						}
						rm1 := b.removeCandidateIfPresent(r, c1, v1)
						rm2 := b.removeCandidateIfPresent(r, c1, v2)
						rm3 := b.removeCandidateIfPresent(r, c1, v3)
						if rm1 || rm2 || rm3 {
							tripleChanged = true
						}
					}
					return tripleChanged
				},
			)
			passApps += pa
			if ch {
				changed = true
			}
		}

		for box := range 9 {
			br := (box / 3) * 3
			bc := (box % 3) * 3
			pa, ch := b.applyNakedTripleInUnit(
				func(i int) (int, int) { return br + i/3, bc + i%3 },
				func(b *Board, r1, c1, r2, c2, r3, c3, v1, v2, v3 int) bool {
					tripleChanged := false
					for i := range 9 {
						r, co := br+i/3, bc+i%3
						if (r == r1 && co == c1) || (r == r2 && co == c2) || (r == r3 && co == c3) {
							continue
						}
						rm1 := b.removeCandidateIfPresent(r, co, v1)
						rm2 := b.removeCandidateIfPresent(r, co, v2)
						rm3 := b.removeCandidateIfPresent(r, co, v3)
						if rm1 || rm2 || rm3 {
							tripleChanged = true
						}
					}
					return tripleChanged
				},
			)
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
