package sudoku

// hiddenSingleInUnit scans one unit (row, column, or box): getCell maps slot 0..8 to board coordinates.
func (b Board) hiddenSingleInUnit(getCell func(i int) (row, col int)) (found bool, row, col, val int) {
	for v := 1; v <= 9; v++ {
		count := 0
		var singleRow, singleCol int
		for i := range 9 {
			row, col := getCell(i)
			if b.Cell(row, col) != 0 {
				continue
			}
			if b.HasCandidate(row, col, v) {
				count++
				singleRow, singleCol = row, col
			}
		}
		if count == 1 {
			return true, singleRow, singleCol, v
		}
	}
	return false, -1, -1, 0
}

func (b Board) hiddenSingle() (found bool, row, col, val int) {
	for r := range 9 {
		if found, rowFound, colFound, v := b.hiddenSingleInUnit(func(i int) (int, int) { return r, i }); found {
			return true, rowFound, colFound, v
		}
	}
	for c := range 9 {
		if found, rowFound, colFound, v := b.hiddenSingleInUnit(func(i int) (int, int) { return i, c }); found {
			return true, rowFound, colFound, v
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if found, rowFound, colFound, v := b.hiddenSingleInUnit(func(i int) (int, int) { return br + i/3, bc + i%3 }); found {
			return true, rowFound, colFound, v
		}
	}
	return false, -1, -1, 0
}

func (b *Board) ApplyHiddenSingles() (applied bool) {
	found, row, col, val := b.hiddenSingle()
	if !found {
		return false
	}
	b.SetCellAndUpdateCandidates(row, col, val)
	return true
}
