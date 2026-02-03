package main

import (
	"fmt"

	"github.com/PhilippSchweizer/sudoku-engine/internal/sudoku"
)

func main() {
	puzzle, solution := sudoku.GeneratePuzzle()
	fmt.Println(puzzle)
	fmt.Println(solution)
}
