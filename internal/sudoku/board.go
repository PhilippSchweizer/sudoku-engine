package sudoku

type Board struct {
	Cells      [9][9]int
	Candidates [9][9]uint16
}

func bitFor(val int) uint16 {
	return uint16(1) << uint16(val-1)
}

func New() Board {
	return Board{}
}

func (b Board) Cell(row, col int) int {
	return b.Cells[row][col]
}

func (b *Board) SetCell(row, col, val int) {
	if val < 0 || val > 9 {
		panic("sudoku: SetCell value must be between 0 and 9")
	}
	b.Cells[row][col] = val
}

func (b *Board) SetCellAndUpdateCandidates(row, col, val int) {
	b.SetCell(row, col, val)
	b.updateCandidatesForPlacement(row, col, val)
}

func (b Board) HasCandidate(row, col, val int) bool {
	return b.Candidates[row][col]&bitFor(val) != 0
}

func (b Board) GetCandidates(row, col int) uint16 {
	return b.Candidates[row][col]
}

func (b *Board) AddCandidate(row, col, val int) {
	b.Candidates[row][col] |= bitFor(val)
}

func (b *Board) RemoveCandidate(row, col, val int) {
	b.Candidates[row][col] &^= bitFor(val)
}

func (b *Board) UpdateCandidates() {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) != 0 {
				b.Candidates[r][c] = 0
				continue
			}
			const allDigitsMask uint16 = (1 << 9) - 1
			usedMask := uint16(0)

			br := (r / 3) * 3
			bc := (c / 3) * 3
			for i := range 9 {
				if v := b.Cell(r, i); v != 0 {
					usedMask |= bitFor(v)

				}
				if v := b.Cell(i, c); v != 0 {
					usedMask |= bitFor(v)
				}
				if v := b.Cell(br+i/3, bc+i%3); v != 0 {
					usedMask |= bitFor(v)
				}
			}
			candidates := allDigitsMask &^ usedMask
			b.Candidates[r][c] = candidates
		}

	}

}

func (b *Board) updateCandidatesForPlacement(row, col, val int) {
	br := (row / 3) * 3
	bc := (col / 3) * 3

	for i := range 9 {
		b.RemoveCandidate(row, i, val)
		b.RemoveCandidate(i, col, val)
		b.RemoveCandidate(br+i/3, bc+i%3, val)
	}
}
