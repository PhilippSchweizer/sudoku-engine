package sudoku

import "math/rand"

func Generate() Board {
	b := New()
	generated, _ := generate(b)
	return generated
}

func generate(b Board) (Board, bool) {
	if b.IsSolved() {
		return b, true
	}

	row, col := b.nextEmpty()
	if row < 0 {
		return b, true
	}

	path := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(path), func(i, j int) {
		path[i], path[j] = path[j], path[i]
	})

	for _, v := range path {
		newBoard := b
		newBoard.SetCell(row, col, v)
		if !newBoard.UnitsValidAt(row, col) {
			continue
		}
		full, ok := generate(newBoard)
		if ok {
			return full, true
		}
	}

	return b, false

}

func (b Board) AllCellPositions() [][2]int {
	allCells := [][2]int{}
	for r := range 9 {
		for c := range 9 {
			allCells = append(allCells, [2]int{r, c})
		}
	}
	return allCells
}

func GeneratePuzzle() (puzzle, solution Board) {
	solution = Generate()
	puzzle = solution
	allCells := puzzle.AllCellPositions()
	rand.Shuffle(len(allCells), func(i, j int) {
		allCells[i], allCells[j] = allCells[j], allCells[i]

	})

	for i := range allCells {
		row, col := allCells[i][0], allCells[i][1]
		currentValue := puzzle.Cell(row, col)
		puzzle.SetCell(row, col, 0)
		if CountSolutions(puzzle) != 1 {
			puzzle.SetCell(row, col, currentValue)
		}
	}

	return puzzle, solution
}
