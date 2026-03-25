package sudoku

// removeCandidateIfPresent removes val from (row,col) if the cell is empty and has that candidate.
// Reports whether the candidate set changed.
func (b *Board) removeCandidateIfPresent(row, col, val int) bool {
	if b.Cell(row, col) != 0 {
		return false
	}
	if !b.HasCandidate(row, col, val) {
		return false
	}
	b.RemoveCandidate(row, col, val)
	return true
}
