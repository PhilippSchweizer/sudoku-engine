package sudoku

func (b Board) IsValid() bool {
	// loop through rows
	for r := range 9 {
		var unit [9]int
		for i := range 9 {
			unit[i] = b.Cell(r, i)
		}
		if hasDuplicateInUnit(unit) {
			return false
		}
	}

	// loop through columns
	for c := range 9 {
		var unit [9]int
		for i := range 9 {
			unit[i] = b.Cell(i, c)
		}
		if hasDuplicateInUnit(unit) {
			return false
		}
	}

	// loop through boxes
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		var unit [9]int
		for i := range 9 {
			unit[i] = b.Cell(br+i/3, bc+i%3)
		}
		if hasDuplicateInUnit(unit) {
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
