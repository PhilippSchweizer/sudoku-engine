package sudoku

import "math/bits"

func (b Board) nakedSingle() (found bool, row, col, val int) {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) != 0 {
				continue
			}
			mask := b.GetCandidates(r, c)
			if bits.OnesCount16(mask) != 1 {
				continue
			}
			val := bits.TrailingZeros16(mask) + 1
			return true, r, c, val
		}
	}
	return false, -1, -1, 0
}

func (b *Board) ApplyNakedSingles() (applied bool) {
	found, row, col, val := b.nakedSingle()
	if !found {
		return false
	}

	b.SetCellAndUpdateCandidates(row, col, val)
	return true
}
