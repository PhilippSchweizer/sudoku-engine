package main

import (
	"strings"

	"github.com/PhilippSchweizer/sudoku-engine/internal/sudoku"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	puzzle    sudoku.Board
	current   sudoku.Board
	solution  sudoku.Board
	cursorRow int
	cursorCol int
}

const (
	givenStart  = "\x1b[37m"
	userStart   = "\x1b[36m"
	cursorStart = "\x1b[46m\x1b[30m" // cyan background, black text
	colorEnd    = "\x1b[0m"          // reset all attributes
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		switch k {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			m.cursorRow = (m.cursorRow + 8) % 9
			//return m.cursorRow
		case "down", "j":
			m.cursorRow = (m.cursorRow + 1) % 9
		case "left", "h":
			m.cursorCol = (m.cursorCol + 8) % 9
		case "right", "l":
			m.cursorCol = (m.cursorCol + 1) % 9
		case "w":
			m.cursorCol = (m.cursorCol + 3) % 9
		case "W":
			m.cursorCol = (m.cursorCol + 6) % 9
		case "e":
			m.cursorRow = (m.cursorRow + 3) % 9
		case "E":
			m.cursorRow = (m.cursorRow + 6) % 9

		default:
			if len(k) == 1 && k[0] >= '0' && k[0] <= '9' {
				r, c := m.cursorRow, m.cursorCol
				if m.puzzle.Cell(r, c) == 0 {
					v := int(k[0] - '0')
					m.current.SetCell(r, c, v)

				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	// return m.puzzle.String()
	return renderBoardWithCursor(m.puzzle, m.current, m.cursorRow, m.cursorCol)
}

func main() {
	puzzle, solution := sudoku.GeneratePuzzle()
	current := puzzle
	m := model{
		puzzle:    puzzle,
		current:   current,
		solution:  solution,
		cursorRow: 0,
		cursorCol: 0,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
	}
}

func renderBoardWithCursor(puzzle, current sudoku.Board, row, col int) string {
	var sb strings.Builder
	line := "+-------------+-------------+-------------+\n"
	sb.WriteString(line)
	for r := range 9 {
		sb.WriteString("| ")
		for c := range 9 {
			v := current.Cell(r, c)
			isGiven := puzzle.Cell(r, c) != 0
			if row != r || col != c {
				if v == 0 {
					sb.WriteString(" .  ")
				} else {
					if isGiven {
						sb.WriteString(givenStart)
					} else {
						sb.WriteString(userStart)
					}
					sb.WriteString(" ")
					sb.WriteByte(byte('0' + v))
					sb.WriteString("  ")
					sb.WriteString(colorEnd)
				}
			} else {
				if v == 0 {
					sb.WriteString(cursorStart)
					// sb.WriteString("[.]")
					sb.WriteString(" . ")
					sb.WriteString(colorEnd)
					sb.WriteString(" ")
				} else {
					sb.WriteString(cursorStart)
					// sb.WriteString("[")
					sb.WriteString(" ")
					sb.WriteByte(byte('0' + v))
					// sb.WriteString("]")
					sb.WriteString(" ")
					sb.WriteString(colorEnd)
					sb.WriteString(" ")
				}
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
