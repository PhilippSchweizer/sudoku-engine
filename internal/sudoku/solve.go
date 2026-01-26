package sudoku

func Solve(b Board) (Board, bool) {
	if BoardSolved(b) {
		return b, true
	}

	if !BoardValid(b) {
		return b, false
	}

	var emptyCell [2]int = findNextEmpty(b)

	for c := 1; c <= 9; c++ {
		newBoard := Set(b, emptyCell[0], emptyCell[1], c)
		solvedBoard, solved := Solve(newBoard)
		if solved {
			return solvedBoard, true
		}
	}

	return b, false
}

func findNextEmpty(b Board) [2]int {
	for row := range len(b) {
		for col := range len(b) {
			if b[row][col] == 0 {
				return [2]int{row, col}
			}
		}
	}
	return [2]int{-1, -1} // no empty cell found
}
