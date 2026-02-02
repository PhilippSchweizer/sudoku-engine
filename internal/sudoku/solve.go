package sudoku

import (
	"math/bits"
)

func Solve(b Board) (Board, bool) {
	if !b.IsValid() {
		return b, false
	}
	return solve(b)
}

func solve(b Board) (Board, bool) {
	if b.IsSolved() {
		return b, true
	}

	row, col := b.nextEmpty()

	for v := 1; v <= 9; v++ {
		newBoard := b
		newBoard.SetCell(row, col, v)
		if !newBoard.UnitsValidAt(row, col) {
			continue
		}
		solvedBoard, solved := solve(newBoard)
		if solved {
			return solvedBoard, true
		}
	}

	return b, false
}

func CountSolutions(b Board) int {
	if !b.IsValid() {
		return 0
	}
	return countSolutions(b)
}

func countSolutions(b Board) int {
	if b.IsSolved() {
		return 1
	}

	row, col := b.nextEmpty()
	count := 0

	for v := 1; v <= 9; v++ {
		newBoard := b
		newBoard.SetCell(row, col, v)
		if !newBoard.UnitsValidAt(row, col) {
			continue
		}
		subCount := countSolutions(newBoard)
		count += subCount
		if count >= 2 {
			return count
		}
	}
	return count
}

func (b Board) nextEmpty() (row, col int) {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) == 0 {
				return r, c
			}
		}
	}
	return -1, -1 // no empty cell found
}

func (b Board) nakedSingle() (found bool, row, col, val int) {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) != 0 {
				continue
			}
			mask := b.GetCandidates(r, c)
			if bits.OnesCount16(mask) != 1 {
				continue
			}
			val := bits.TrailingZeros16(mask) + 1
			return true, r, c, val
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

func (b Board) nakedPair() (found bool, row1, col1, row2, col2 int, val1, val2 int) {
	for r := range 9 {
		if found, r1, c1, r2, c2, v1, v2 := b.nakedPairInUnit(func(i int) (int, int) { return r, i }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	for c := range 9 {
		if found, r1, c1, r2, c2, v1, v2 := b.nakedPairInUnit(func(i int) (int, int) { return i, c }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if found, r1, c1, r2, c2, v1, v2 := b.nakedPairInUnit(func(i int) (int, int) { return br + i/3, bc + i%3 }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

func (b Board) hiddenPair() (found bool, row1, col1, row2, col2 int, val1, val2 int) {
	for r := range 9 {
		if found, r1, c1, r2, c2, v1, v2 := b.hiddenPairInUnit(func(i int) (int, int) { return r, i }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	for c := range 9 {
		if found, r1, c1, r2, c2, v1, v2 := b.hiddenPairInUnit(func(i int) (int, int) { return i, c }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if found, r1, c1, r2, c2, v1, v2 := b.hiddenPairInUnit(func(i int) (int, int) { return br + i/3, bc + i%3 }); found {
			return true, r1, c1, r2, c2, v1, v2
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

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

func (b Board) nakedPairInUnit(getCell func(i int) (row, col int)) (found bool, row1, col1, row2, col2, val1, val2 int) {
	// Find two empty cells in the unit that have exactly the same two candidates.
	for i := range 9 {
		row, col := getCell(i)
		if b.Cell(row, col) != 0 {
			continue
		}
		mask := b.GetCandidates(row, col)
		if bits.OnesCount16(mask) != 2 {
			continue
		}
		for j := i + 1; j < 9; j++ {
			rowOther, colOther := getCell(j)
			if b.Cell(rowOther, colOther) != 0 {
				continue
			}
			otherMask := b.GetCandidates(rowOther, colOther)
			if otherMask != mask || bits.OnesCount16(otherMask) != 2 {
				continue
			}
			// Same two candidates. Extract the two digits from the mask.
			v1 := bits.TrailingZeros16(mask) + 1
			v2 := bits.TrailingZeros16(mask&^(uint16(1)<<(v1-1))) + 1
			return true, row, col, rowOther, colOther, v1, v2
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

func (b Board) hiddenPairInUnit(getCell func(i int) (row, col int)) (found bool, row1, col1, row2, col2, val1, val2 int) {
	type coord struct{ r, c int }
	positions := [9][]coord{}
	for p := range 9 {
		row, col := getCell(p)
		if b.Cell(row, col) != 0 {
			continue
		}
		for v := 1; v <= 9; v++ {
			if b.HasCandidate(row, col, v) {
				pos := coord{row, col}
				positions[v-1] = append(positions[v-1], pos)
			}
		}
	}

	for v1 := range 9 {
		if len(positions[v1]) != 2 {
			continue
		}
		for v2 := v1 + 1; v2 < 9; v2++ {
			if len(positions[v2]) != 2 {
				continue
			}
			if positions[v1][0] == positions[v2][0] && positions[v1][1] == positions[v2][1] {
				c1 := positions[v1][0]
				row1, col1 = c1.r, c1.c
				c2 := positions[v1][1]
				row2, col2 = c2.r, c2.c
				return true, row1, col1, row2, col2, v1 + 1, v2 + 1
			}
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

func (b *Board) ApplyNakedSingle() (applied bool) {
	found, row, col, val := b.nakedSingle()
	if !found {
		return false
	}
	if found {
		b.SetCellAndUpdateCandidates(row, col, val)
		return true
	}
	return false
}

func (b *Board) ApplyHiddenSingle() (applied bool) {
	found, row, col, val := b.hiddenSingle()
	if !found {
		return false
	}
	b.SetCellAndUpdateCandidates(row, col, val)
	return true
}

func (b *Board) ApplyNakedPair() (applied bool) {
	found, row1, col1, row2, col2, val1, val2 := b.nakedPair()

	if !found {
		return false
	}

	var sameRow bool = false
	var sameCol bool = false
	var sameBox bool = false

	if row1 == row2 {
		sameRow = true
	} else if col1 == col2 {
		sameCol = true
	} else {
		sameBox = true
	}

	if sameRow {
		for c := range 9 {
			if c == col1 || c == col2 {
				continue
			}
			b.RemoveCandidate(row1, c, val1)
			b.RemoveCandidate(row1, c, val2)
		}
	}
	if sameCol {
		for r := range 9 {
			if r == row1 || r == row2 {
				continue
			}
			b.RemoveCandidate(r, col1, val1)
			b.RemoveCandidate(r, col1, val2)
		}
	}
	if sameBox {
		br := (row1 / 3) * 3
		bc := (col1 / 3) * 3
		for i := range 9 {
			if ((br+i/3) == row1 && (bc+i%3) == col1) || ((br+i/3) == row2 && (bc+i%3) == col2) {
				continue
			}
			b.RemoveCandidate(br+i/3, bc+i%3, val1)
			b.RemoveCandidate(br+i/3, bc+i%3, val2)
		}
	}
	return true
}

func (b *Board) ApplyHiddenPair() (applied bool) {
	found, row1, col1, row2, col2, val1, val2 := b.hiddenPair()
	if !found {
		return false
	}
	for v := 1; v <= 9; v++ {
		if v == val1 || v == val2 {
			continue
		}
		b.RemoveCandidate(row1, col1, v)
		b.RemoveCandidate(row2, col2, v)

	}
	return true
}
