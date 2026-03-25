package sudoku

import "math/bits"

// applyNakedPairsInUnit finds every naked pair in one unit. getCell maps slot 0..8 to coordinates;
// eliminate removes the pair digits from other cells that share the constraint for that unit shape.
func (b *Board) applyNakedPairsInUnit(
	getCell func(i int) (row, col int),
	eliminate func(b *Board, r1, c1, r2, c2, v1, v2 int) (pairChanged bool),
) (pairApplications int, changed bool) {
	for i := range 9 {
		row, col := getCell(i)
		if b.Cell(row, col) != 0 {
			continue
		}
		mask := b.GetCandidates(row, col)
		if bits.OnesCount16(mask) != 2 {
			continue
		}
		v1 := bits.TrailingZeros16(mask) + 1
		v2 := bits.TrailingZeros16(mask&^(uint16(1)<<(v1-1))) + 1
		for j := i + 1; j < 9; j++ {
			row2, col2 := getCell(j)
			if b.Cell(row2, col2) != 0 {
				continue
			}
			otherMask := b.GetCandidates(row2, col2)
			if otherMask != mask || bits.OnesCount16(otherMask) != 2 {
				continue
			}
			if eliminate(b, row, col, row2, col2, v1, v2) {
				pairApplications++
				changed = true
			}
		}
	}
	return pairApplications, changed
}

func (b *Board) ApplyNakedPairs() (applied bool, applications int) {
	overall := false
	for {
		changed := false
		passApps := 0

		for r := range 9 {
			pa, ch := b.applyNakedPairsInUnit(
				func(i int) (int, int) { return r, i },
				func(b *Board, r1, c1, r2, c2, v1, v2 int) bool {
					pairChanged := false
					for c := range 9 {
						if c == c1 || c == c2 {
							continue
						}
						rm1 := b.removeCandidateIfPresent(r1, c, v1)
						rm2 := b.removeCandidateIfPresent(r1, c, v2)
						if rm1 || rm2 {
							pairChanged = true
						}
					}
					return pairChanged
				},
			)
			passApps += pa
			if ch {
				changed = true
			}
		}

		for c := range 9 {
			pa, ch := b.applyNakedPairsInUnit(
				func(i int) (int, int) { return i, c },
				func(b *Board, r1, c1, r2, c2, v1, v2 int) bool {
					pairChanged := false
					for rr := range 9 {
						if rr == r1 || rr == r2 {
							continue
						}
						rm1 := b.removeCandidateIfPresent(rr, c1, v1)
						rm2 := b.removeCandidateIfPresent(rr, c1, v2)
						if rm1 || rm2 {
							pairChanged = true
						}
					}
					return pairChanged
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
			pa, ch := b.applyNakedPairsInUnit(
				func(i int) (int, int) { return br + i/3, bc + i%3 },
				func(b *Board, r1, c1, r2, c2, v1, v2 int) bool {
					pairChanged := false
					for i := range 9 {
						r, co := br+i/3, bc+i%3
						if (r == r1 && co == c1) || (r == r2 && co == c2) {
							continue
						}
						rm1 := b.removeCandidateIfPresent(r, co, v1)
						rm2 := b.removeCandidateIfPresent(r, co, v2)
						if rm1 || rm2 {
							pairChanged = true
						}
					}
					return pairChanged
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
