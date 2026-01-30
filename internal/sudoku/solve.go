package sudoku

import (
	"math/bits"
)

func Solve(b Board) (Board, bool) {
	if BoardSolved(b) {
		return b, true
	}

	if !BoardValid(b) {
		return b, false
	}

	var emptyCell [2]int = findNextEmpty(b)

	for c := 1; c <= 9; c++ {
		newBoard := b
		newBoard.SetCell(emptyCell[0], emptyCell[1], c)
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
		newBoard := b
		newBoard.SetCell(emptyCell[0], emptyCell[1], c)
		subCount := CountSolutions(newBoard)
		count += subCount
		if count >= 2 {
			return count
		}
	}
	return count
}

func findNextEmpty(b Board) [2]int {
	for row := range len(b.Cells) {
		for col := range len(b.Cells) {
			if b.Cell(row, col) == 0 {
				return [2]int{row, col}
			}
		}
	}
	return [2]int{-1, -1} // no empty cell found
}

func findNakedSingle(b Board) (found bool, row, col, val int) {
	for row := range 9 {
		for col := range 9 {
			if b.Cell(row, col) != 0 {
				continue
			}
			mask := b.GetCandidates(row, col)
			if bits.OnesCount16(mask) != 1 {
				continue
			}
			val := bits.TrailingZeros16(mask) + 1
			return true, row, col, val
		}
	}
	return false, -1, -1, 0
}

func findHiddenSingle(b Board) (found bool, row, col, val int) {
	for r := range 9 {
		if found, rr, cc, v := findHiddenSingleInUnit(b, func(i int) (int, int) { return r, i }); found {
			return true, rr, cc, v
		}
	}
	for c := range 9 {
		if found, rr, cc, v := findHiddenSingleInUnit(b, func(i int) (int, int) { return i, c }); found {
			return true, rr, cc, v
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if found, rr, cc, v := findHiddenSingleInUnit(b, func(i int) (int, int) { return br + i/3, bc + i%3 }); found {
			return true, rr, cc, v
		}
	}
	return false, -1, -1, 0
}

func findNakedPair(b Board) (found bool, r1, c1, r2, c2 int, d1, d2 int) {
	for r := range 9 {
		if found, rr1, cc1, rr2, cc2, dd1, dd2 := findNakedPairInUnit(b, func(i int) (int, int) { return r, i }); found {
			return true, rr1, cc1, rr2, cc2, dd1, dd2
		}
	}
	for c := range 9 {
		if found, rr1, cc1, rr2, cc2, dd1, dd2 := findNakedPairInUnit(b, func(i int) (int, int) { return i, c }); found {
			return true, rr1, cc1, rr2, cc2, dd1, dd2
		}
	}
	for box := range 9 {
		br := (box / 3) * 3
		bc := (box % 3) * 3
		if found, rr1, cc1, rr2, cc2, dd1, dd2 := findNakedPairInUnit(b, func(i int) (int, int) { return br + i/3, bc + i%3 }); found {
			return true, rr1, cc1, rr2, cc2, dd1, dd2
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

func findHiddenSingleInUnit(b Board, getCell func(i int) (row, col int)) (found bool, row, col, val int) {
	for d := 1; d <= 9; d++ {
		count := 0
		var singleR, singleC int
		for i := range 9 {
			r, c := getCell(i)
			if b.Cell(r, c) != 0 {
				continue
			}
			if b.HasCandidate(r, c, d) {
				count++
				singleR, singleC = r, c
			}
		}
		if count == 1 {
			return true, singleR, singleC, d
		}
	}
	return false, -1, -1, 0
}

func findNakedPairInUnit(b Board, getCell func(i int) (row, col int)) (found bool, r1, c1, r2, c2, d1, d2 int) {
	// Find two empty cells in the unit that have exactly the same two candidates.
	for i := 0; i < 9; i++ {
		ri, ci := getCell(i)
		if b.Cell(ri, ci) != 0 {
			continue
		}
		mi := b.GetCandidates(ri, ci)
		if bits.OnesCount16(mi) != 2 {
			continue
		}
		for j := i + 1; j < 9; j++ {
			rj, cj := getCell(j)
			if b.Cell(rj, cj) != 0 {
				continue
			}
			mj := b.GetCandidates(rj, cj)
			if mj != mi || bits.OnesCount16(mj) != 2 {
				continue
			}
			// Same two candidates. Extract the two digits from the mask.
			d1 = bits.TrailingZeros16(mi) + 1
			d2 = bits.TrailingZeros16(mi&^(uint16(1)<<(d1-1))) + 1
			return true, ri, ci, rj, cj, d1, d2
		}
	}
	return false, -1, -1, -1, -1, 0, 0
}

func (b *Board) ApplyNakedSingle() (applied bool) {
	found, r, c, v := findNakedSingle(*b)
	if !found {
		return false
	}
	if found {
		b.SetCellWithCandidateUpdate(r, c, v)
		return true
	}
	return false
}

func (b *Board) ApplyHiddenSingle() (applied bool) {
	found, r, c, v := findHiddenSingle(*b)
	if !found {
		return false
	}
	b.SetCellWithCandidateUpdate(r, c, v)
	return true
}

func (b *Board) ApplyNakedPair() (applied bool) {
	found, r1, c1, r2, c2, d1, d2 := findNakedPair(*b)

	if !found {
		return false
	}

	var sameRow bool = false
	var sameCol bool = false
	var sameBox bool = false

	if r1 == r2 {
		sameRow = true
	} else if c1 == c2 {
		sameCol = true
	} else {
		sameBox = true
	}

	if sameRow {
		for r := range 9 {
			if r == c1 || r == c2 {
				continue
			}
			b.RemoveCandidate(r1, r, d1)
			b.RemoveCandidate(r1, r, d2)
		}
	}
	if sameCol {
		for c := range 9 {
			if c == r1 || c == r2 {
				continue
			}
			b.RemoveCandidate(c, c1, d1)
			b.RemoveCandidate(c, c1, d2)
		}
	}
	if sameBox {
		br := (r1 / 3) * 3
		bc := (c1 / 3) * 3
		for i := range 9 {
			if ((br+i/3) == r1 && (bc+i%3) == c1) || ((br+i/3) == r2 && (bc+i%3) == c2) {
				continue
			}
			b.RemoveCandidate(br+i/3, bc+i%3, d1)
			b.RemoveCandidate(br+i/3, bc+i%3, d2)
		}
	}
	return true
}
