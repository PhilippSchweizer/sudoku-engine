package sudoku

func BoardValid(b Board) bool {
	for row := range 9 {
		var unit [9]int
		for i := range 9 {
			unit[i] = b[row][i]
		}
		if hasDuplicateInUnit(unit) {
			return false
		}
	}

	for column := range 9 {
		var unit [9]int
		for i := range 9 {
			unit[i] = b[i][column]
		}
		if hasDuplicateInUnit(unit) {
			return false
		}
	}

	// TODO: boxes
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		var unit [9]int
		for i := range 9 {
			unit[i] = b[br+i/3][bc+i%3]
		}
		if hasDuplicateInUnit(unit) {
			return false
		}

	}
	return true
}

func hasDuplicateInUnit(unit [9]int) bool {
	seen := [9]bool{}
	for _, v := range unit {
		if v == 0 {
			continue
		}
		if seen[v-1] == true {
			return true
		} else {
			seen[v-1] = true
		}
	}
	return false
	/* for i := 0; i < len(unit)-1; i++ {
		for j := i + 1; j < len(unit); j++ {
			if unit[i] == unit[j] {
				return true
			}
		}
	} */
}
