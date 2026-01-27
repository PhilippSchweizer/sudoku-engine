package sudoku

type Board struct {
	Cells      [9][9]int
	Candidates [9][9]uint16
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
