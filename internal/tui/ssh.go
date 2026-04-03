package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"

	"github.com/PhilippSchweizer/sudoku-engine/internal/sudoku"
)

// NewRemoteSessionModel returns a Bubble Tea model for one SSH session (use with charm.land/wish/v2/bubbletea).
// A fresh puzzle is generated for each connection.
func NewRemoteSessionModel(_ ssh.Session) (tea.Model, []tea.ProgramOption) {
	puzzle, _ := sudoku.GeneratePuzzle()
	current := puzzle
	return newModel(puzzle, current), nil
}
