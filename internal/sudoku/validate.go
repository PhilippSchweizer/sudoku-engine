package sudoku

func (b Board) IsValid() bool {
	for r := range 9 {
		if hasDuplicateInUnit(b.RowAsUnit(r)) {
			return false
		}
	}
	for c := range 9 {
		if hasDuplicateInUnit(b.ColAsUnit(c)) {
			return false
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if hasDuplicateInUnit(b.BoxAsUnit(br, bc)) {
			return false
		}
	}
	return true
}

func (b Board) IsSolved() bool {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) == 0 {
				return false
			}
		}
	}
	return b.IsValid()
}

func (b Board) UnitsValidAt(row, col int) bool {
	return !hasDuplicateInUnit(b.RowAsUnit(row)) &&
		!hasDuplicateInUnit(b.ColAsUnit(col)) &&
		!hasDuplicateInUnit(b.BoxAsUnit(row, col))
}

func (b Board) RowAsUnit(row int) [9]int {
	unit := [9]int{}
	for i := range 9 {
		unit[i] = b.Cell(row, i)
	}
	return unit
}

func (b Board) ColAsUnit(col int) [9]int {
	unit := [9]int{}
	for i := range 9 {
		unit[i] = b.Cell(i, col)
	}
	return unit
}

func (b Board) BoxAsUnit(row, col int) [9]int {
	unit := [9]int{}
	br := (row / 3) * 3
	bc := (col / 3) * 3
	for i := range 9 {
		unit[i] = b.Cell(br+i/3, bc+i%3)
	}
	return unit
}

func hasDuplicateInUnit(unit [9]int) bool {
	seen := [9]bool{}
	for _, v := range unit {
		if v == 0 {
			continue
		}
		if seen[v-1] {
			return true
		} else {
			seen[v-1] = true
		}
	}
	return false
}
