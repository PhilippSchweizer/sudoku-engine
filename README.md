# Golang Sudoku Engine (WIP)

> A Go-based Sudoku backend focused on learning. This is a work in progress.

## Why This Exists
- Learning Go, algorithms, and efficiency tradeoffs
- Building a portfolio piece (I coded it myself, i swear)
- Having fun

## Status
- **Current:** Backtracking solver, validation, and candidate tracking
- **WIP:** Puzzle generation with human-solvable constraints
- **Planned:** API surface for frontend/desktop/mobile/shell clients

## Features (WIP)
- Backtracking solver with solution counting
- Board validation (rows, columns, boxes)
- Candidate tracking + basic human-style techniques (naked/hidden singles, pairs)
- Early groundwork for difficulty evaluation

## Usage
### Run the demo
```bash
go run .
```

### Example (API usage)
```go
package main

import (
	"fmt"

	"github.com/PhilippSchweizer/sudoku-engine/internal/sudoku"
)

func main() {
	board := sudoku.New()
	solved, ok := sudoku.Solve(board)
	if ok {
		fmt.Printf("%+v\n", solved)
	}
}
```

## Tests
```bash
go test ./...
```

## Roadmap
- Implement puzzle generation with unique solutions
- Add difficulty grading based on human techniques
- Expose a clean API for UIs and other clients

## Notes
- This is an early-stage project; APIs and behavior may change

## License
TBD
