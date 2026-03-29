package main

import (
	"os"

	"github.com/PhilippSchweizer/sudoku-engine/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		os.Exit(1)
	}
}
