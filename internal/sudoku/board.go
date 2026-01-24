package sudoku

type Board [9][9]int

func New() Board {
	return Board{}
}

func Get(board Board, row, col int) int {
	return board[row][col]
}

func Set(board Board, row, col, val int) Board {
	if val < 0 || val > 9 {
		return board
	}
	b := board
	b[row][col] = val
	return b
}
