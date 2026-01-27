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

func CountSolutions(b Board) int {
	if BoardSolved(b) {
		return 1
	}

	if !BoardValid(b) {
		return 0
	}

	emptyCell := findNextEmpty(b)
	if emptyCell[0] == -1 {
		return 1
	}
	count := 0

	for c := 1; c <= 9; c++ {
		newBoard := Set(b, emptyCell[0], emptyCell[1], c)
		subCount := CountSolutions(newBoard)
		count += subCount
		if count >= 2 {
			return count
		}
	}
	return count
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
