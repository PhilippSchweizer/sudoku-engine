package sudoku

import "strings"

type Board struct {
	Cells      [9][9]int
	Candidates [9][9]uint16
}

func bitFor(val int) uint16 {
	return uint16(1) << uint16(val-1)
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

func (b *Board) SetCellAndUpdateCandidates(row, col, val int) {
	b.SetCell(row, col, val)
	b.updateCandidatesForPlacement(row, col, val)
}

func (b Board) HasCandidate(row, col, val int) bool {
	return b.Candidates[row][col]&bitFor(val) != 0
}

func (b Board) GetCandidates(row, col int) uint16 {
	return b.Candidates[row][col]
}

func (b *Board) AddCandidate(row, col, val int) {
	b.Candidates[row][col] |= bitFor(val)
}

func (b *Board) RemoveCandidate(row, col, val int) {
	b.Candidates[row][col] &^= bitFor(val)
}

func (b *Board) UpdateCandidates() {
	for r := range 9 {
		for c := range 9 {
			if b.Cell(r, c) != 0 {
				b.Candidates[r][c] = 0
				continue
			}
			const allDigitsMask uint16 = (1 << 9) - 1
			usedMask := uint16(0)

			br := (r / 3) * 3
			bc := (c / 3) * 3
			for i := range 9 {
				if v := b.Cell(r, i); v != 0 {
					usedMask |= bitFor(v)

				}
				if v := b.Cell(i, c); v != 0 {
					usedMask |= bitFor(v)
				}
				if v := b.Cell(br+i/3, bc+i%3); v != 0 {
					usedMask |= bitFor(v)
				}
			}
			candidates := allDigitsMask &^ usedMask
			b.Candidates[r][c] = candidates
		}

	}

}

func (b *Board) updateCandidatesForPlacement(row, col, val int) {
	br := (row / 3) * 3
	bc := (col / 3) * 3

	for i := range 9 {
		b.RemoveCandidate(row, i, val)
		b.RemoveCandidate(i, col, val)
		b.RemoveCandidate(br+i/3, bc+i%3, val)
	}
}

func (b Board) String() string {
	var sb strings.Builder
	line := "+-------+-------+-------+\n"
	sb.WriteString(line)
	for r := range 9 {
		sb.WriteString("| ")
		for c := range 9 {
			v := b.Cell(r, c)
			if v == 0 {
				sb.WriteString(". ")
			} else {
				sb.WriteByte(byte('0' + v))
				sb.WriteString(" ")
			}
			if c == 2 || c == 5 {
				sb.WriteString("| ")
			}
		}
		sb.WriteString("|\n")
		if r == 2 || r == 5 {
			sb.WriteString(line)
		}
	}
	sb.WriteString(line)
	return sb.String()
}

// pencilMarkLine returns one of three 3-character lines for the miniature 1–9 grid
// in an empty cell (subrow 0 → digits 1–3, 1 → 4–6, 2 → 7–9). Filled cells show the digit on the middle line only.
func (b Board) pencilMarkLine(row, col, subrow int) string {
	v := b.Cell(row, col)
	if v != 0 {
		if subrow == 1 {
			return string([]byte{' ', byte('0' + v), ' '})
		}
		return "   "
	}
	mask := b.Candidates[row][col]
	start := subrow*3 + 1
	var buf [3]byte
	for i := range 3 {
		d := start + i
		if mask&bitFor(d) != 0 {
			buf[i] = byte('0' + d)
		} else {
			buf[i] = '.'
		}
	}
	return string(buf[:])
}

// Content width of one FormatWithPencilMarks row (between outer '+' borders).
const pencilMarkGridInnerWidth = 41

func (b Board) pencilMarkHorizontalRule() string {
	return "+" + strings.Repeat("-", pencilMarkGridInnerWidth) + "+\n"
}

// FormatWithPencilMarks renders the grid with 3×3 pencil marks per empty cell
// (digits present, '.' eliminated). Filled cells show only the placed digit on the center line of the cell block.
func (b Board) FormatWithPencilMarks() string {
	var sb strings.Builder
	sb.WriteString(b.pencilMarkHorizontalRule())
	for r := range 9 {
		for sub := range 3 {
			sb.WriteString("| ")
			for c := range 9 {
				sb.WriteString(b.pencilMarkLine(r, c, sub))
				if c < 8 {
					if c%3 == 2 {
						sb.WriteString(" | ")
					} else {
						sb.WriteByte(' ')
					}
				}
			}
			sb.WriteString(" |\n")
		}
		if r < 8 && r%3 == 2 {
			sb.WriteString(b.pencilMarkHorizontalRule())
		}
	}
	sb.WriteString(b.pencilMarkHorizontalRule())
	return sb.String()
}
