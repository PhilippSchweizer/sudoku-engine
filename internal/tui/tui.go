package tui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/PhilippSchweizer/sudoku-engine/internal/sudoku"
)

type insertMode int

const sudokuTitle = `
             |\         ,,             
              \\        ||             
 _-_, \\ \\  / \\  /'\\ ||/\ \\ \\ 
||_.  || || || || || || ||_< || || 
 ~ || || || || || || || || | || || 
,-_-  \\/\\  \\/  \\,/  \\,\ \\/\\ 
`

const (
	modeValue insertMode = iota
	modePencil
)

const (
	gapBoardSidebar = 1
	gridH           = '-' // horizontal rule fill
	gridV           = '|' // vertical separator
	gridX           = '*' // rule intersections (corners / column joins)
	// Sidebar at or below this width uses the compact "SUDOKU" title. The ASCII
	// banner is wider than most sidebars but still draws fine; comparing maxW > w
	// was almost always true and hid the art.
	titleFallbackMaxSidebarW = 25
)

type model struct {
	puzzle    sudoku.Board
	current   sudoku.Board
	cursorRow int
	cursorCol int
	mode      insertMode
	termW     int
	termH     int
	sidebarW  int
	slotCharW int
	vp        viewport.Model
}

// Tokyo Night–inspired palette (https://github.com/enkia/tokyo-night-vscode-theme)
var (
	tnCursorBG      = lipgloss.Color("#bb9af7")
	tnCursorFG      = lipgloss.Color("#1a1b26")
	tnBorder        = lipgloss.Color("#3b4261")
	tnGridAccent    = lipgloss.Color("#9aa5ce") // board rim + 3×3 box borders
	tnPencil        = lipgloss.Color("#565f89")
	tnGiven         = lipgloss.Color("#7aa2f7")
	tnUser          = lipgloss.Color("#9ece6a")
	tnTitle         = lipgloss.Color("#bb9af7")
	tnLegend        = lipgloss.Color("#565f89")
	tnLegendGrad    = lipgloss.Color("#7dcfff") // cyan accent for border blend
	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(tnTitle)
	legendItemStyle = lipgloss.NewStyle().Foreground(tnLegend)
	vpSepStyle      = lipgloss.NewStyle().
			BorderRight(true).
			BorderForeground(tnBorder).
			MarginRight(gapBoardSidebar).
			PaddingRight(1)
	gapBetweenSlots = " "
)

func (m *model) gridStrokeMuted() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tnBorder)
}

func (m *model) gridStrokeAccent() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tnGridAccent)
}

// gridVertSeg is one text row of the vertical bar (styled).
func (m *model) gridVertSeg(boxEdge bool) string {
	if boxEdge {
		return m.gridStrokeAccent().Render(string(gridV))
	}
	return m.gridStrokeMuted().Render(string(gridV))
}

// gridVertColumn is three rows tall so it lines up with a stacked cell column.
func (m *model) gridVertColumn(boxEdge bool) string {
	s := m.gridVertSeg(boxEdge)
	return lipgloss.JoinVertical(lipgloss.Left, s, s, s)
}

func (m *model) cellContentWidth() int {
	return m.slotCharW*3 + len(gapBetweenSlots)*2
}

// Run starts the sudoku TUI (alternate screen, full-window layout).
func Run() error {
	puzzle, _ := sudoku.GeneratePuzzle()
	current := puzzle

	m := newModel(puzzle, current)
	p := tea.NewProgram(m)
	_, err := p.Run()
	return err
}

func newModel(puzzle, current sudoku.Board) *model {
	vp := viewport.New(viewport.WithWidth(60), viewport.WithHeight(24))
	vp.MouseWheelEnabled = true
	vp.MouseWheelDelta = 3
	return &model{
		puzzle:    puzzle,
		current:   current,
		cursorRow: 0,
		cursorCol: 0,
		mode:      modeValue,
		termW:     80,
		termH:     24,
		sidebarW:  38,
		slotCharW: 3,
		vp:        vp,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) boardMinWidth() int {
	return 2 + 9*m.cellContentWidth() + 8
}

func (m *model) adjustSlotForViewport(vpW int) {
	frame := m.vp.Style.GetHorizontalFrameSize()
	avail := vpW - frame
	if avail < 12 {
		avail = 12
	}
	for _, sw := range []int{3, 2, 1} {
		m.slotCharW = sw
		if m.boardMinWidth() <= avail {
			return
		}
	}
	m.slotCharW = 1
}

func (m *model) relayoutViewport() {
	if m.termW <= 0 || m.termH <= 0 {
		return
	}
	sb := 38
	if m.termW < 100 {
		sb = m.termW / 3
	}
	if sb < 24 {
		sb = 24
	}
	if sb > m.termW/2 {
		sb = m.termW / 2
	}
	m.sidebarW = sb
	vpW := m.termW - sb - gapBoardSidebar
	if vpW < 12 {
		vpW = 12
	}
	m.adjustSlotForViewport(vpW)
	m.vp.SetWidth(vpW)
	m.vp.SetHeight(m.termH)
	m.vp.Style = vpSepStyle.Width(vpW + m.vp.Style.GetHorizontalFrameSize())
}

func (m *model) syncViewportToCursor() {
	frameV := m.vp.Style.GetVerticalFrameSize()
	innerH := m.vp.Height() - frameV
	if innerH <= 0 {
		return
	}
	top := m.lineIndexOfSudokuRowTop(m.cursorRow)
	// cursor uses 3 content lines [top, top+3)
	if top < m.vp.YOffset() {
		m.vp.SetYOffset(top)
		return
	}
	if top+3 > m.vp.YOffset()+innerH {
		m.vp.SetYOffset(top + 3 - innerH)
	}
}

// lineIndexOfSudokuRowTop returns 0-based line in renderBoardFlat for subrow 0 of sudoku row r.
func (m *model) lineIndexOfSudokuRowTop(r int) int {
	y := 1 // line 0 is top border
	for i := 0; i < r; i++ {
		y += 3 // three pencil subrows
		if i < 8 {
			y++ // horizontal rule (thin or band) after each sudoku row
		}
	}
	return y
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termW = msg.Width
		m.termH = msg.Height
		m.relayoutViewport()
		m.vp.SetContent(m.renderBoardFlat())
		m.syncViewportToCursor()
		return m, nil

	case tea.MouseMsg:
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd

	case tea.KeyPressMsg:
		k := msg.String()
		switch k {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.mode == modeValue {
				m.mode = modePencil
			} else {
				m.mode = modeValue
			}
			return m, nil
		case "f":
			m.current.FillCandidatesFromLegal()
		case "pgdown":
			m.vp.PageDown()
			return m, nil
		case "pgup":
			m.vp.PageUp()
			return m, nil
		case "up":
			m.cursorRow = (m.cursorRow + 8) % 9
		case "down":
			m.cursorRow = (m.cursorRow + 1) % 9
		case "left":
			m.cursorCol = (m.cursorCol + 8) % 9
		case "right":
			m.cursorCol = (m.cursorCol + 1) % 9
		case "k":
			m.cursorRow = (m.cursorRow + 8) % 9
		case "j":
			m.cursorRow = (m.cursorRow + 1) % 9
		case "h":
			m.cursorCol = (m.cursorCol + 8) % 9
		case "l":
			m.cursorCol = (m.cursorCol + 1) % 9
		case "w":
			m.cursorCol = (m.cursorCol + 3) % 9
		case "W":
			m.cursorCol = (m.cursorCol + 6) % 9
		case "e":
			m.cursorRow = (m.cursorRow + 3) % 9
		case "E":
			m.cursorRow = (m.cursorRow + 6) % 9
		case "x", "backspace", "delete":
			if m.puzzle.Cell(m.cursorRow, m.cursorCol) == 0 {
				m.current.ClearUserCell(m.cursorRow, m.cursorCol)
			}
		default:
			if len(k) == 1 {
				ch := k[0]
				if ch == '0' {
					if m.mode == modeValue && m.puzzle.Cell(m.cursorRow, m.cursorCol) == 0 {
						m.current.ClearUserCell(m.cursorRow, m.cursorCol)
					}
					m.vp.SetContent(m.renderBoardFlat())
					m.syncViewportToCursor()
					return m, nil
				}
				if ch >= '1' && ch <= '9' {
					m.handleDigit(int(ch - '0'))
				}
			}
		}
		m.vp.SetContent(m.renderBoardFlat())
		m.syncViewportToCursor()
		return m, nil
	}
	return m, nil
}

func (m *model) handleDigit(v int) {
	r, c := m.cursorRow, m.cursorCol
	if m.puzzle.Cell(r, c) != 0 {
		return
	}
	if m.current.Cell(r, c) != 0 {
		if m.mode == modePencil {
			return
		}
		m.current.ClearUserCell(r, c)
		m.current.SetCellAndUpdateCandidates(r, c, v)
		return
	}
	switch m.mode {
	case modeValue:
		m.current.SetCellAndUpdateCandidates(r, c, v)
	case modePencil:
		if m.current.HasCandidate(r, c, v) {
			m.current.RemoveCandidate(r, c, v)
		} else {
			m.current.AddCandidate(r, c, v)
		}
	}
}

func (m *model) widenPencilSubline(r, c, sub int) string {
	raw := m.current.PencilMarkLine(r, c, sub)
	runes := []rune(raw)
	if len(runes) != 3 {
		runes = []rune{' ', ' ', ' '}
	}
	parts := make([]string, 3)
	for i := range 3 {
		ch := string(runes[i])
		if ch == "" {
			ch = " "
		}
		parts[i] = lipgloss.NewStyle().Width(m.slotCharW).Align(lipgloss.Center).Render(ch)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts[0], gapBetweenSlots, parts[1], gapBetweenSlots, parts[2])
}

func (m *model) cellForeground(r, c, sub int) color.Color {
	if r == m.cursorRow && c == m.cursorCol {
		return tnCursorFG
	}
	if m.current.Cell(r, c) != 0 {
		if m.puzzle.Cell(r, c) != 0 {
			return tnGiven
		}
		return tnUser
	}
	return tnPencil
}

func (m *model) styledCellInteriorLine(r, c, sub int) string {
	line := m.widenPencilSubline(r, c, sub)
	fg := m.cellForeground(r, c, sub)
	st := lipgloss.NewStyle().
		Foreground(fg).
		Width(m.cellContentWidth()).
		Align(lipgloss.Center)
	if r == m.cursorRow && c == m.cursorCol {
		st = st.Background(tnCursorBG)
	}
	if m.current.Cell(r, c) == 0 && !(r == m.cursorRow && c == m.cursorCol) {
		st = st.Faint(true)
	}
	return st.Render(line)
}

// styledCellColumn stacks the three pencil sub-rows per cell. JoinVertical(Center)
// keeps line padding consistent when ANSI SGR changes width; we avoid Bold on the
// value row so the digit lines up with blank rows (bold alone often looks vertically shifted).
func (m *model) styledCellColumn(r, c int) string {
	lines := make([]string, 3)
	for sub := range 3 {
		lines[sub] = m.styledCellInteriorLine(r, c, sub)
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines[0], lines[1], lines[2])
}

func (m *model) gridRowWithCells(r int) string {
	parts := make([]string, 0, 20)
	parts = append(parts, m.gridVertColumn(true))
	for c := 0; c < 9; c++ {
		parts = append(parts, m.styledCellColumn(r, c))
		if c < 8 {
			parts = append(parts, m.gridVertColumn(c%3 == 2))
		}
	}
	parts = append(parts, m.gridVertColumn(true))
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (m *model) renderBoardFlat() string {
	rows := make([]string, 0, 32)
	rows = append(rows, m.gridTopBorder())
	for r := 0; r < 9; r++ {
		rows = append(rows, m.gridRowWithCells(r))
		if r%3 == 2 && r < 8 {
			rows = append(rows, m.gridBandSeparator())
		} else if r < 8 {
			rows = append(rows, m.gridInnerSeparator())
		}
	}
	rows = append(rows, m.gridBottomBorder())
	return strings.Join(rows, "\n")
}

// gridRuleRowAllAccent is a full-width horizontal rule using only the accent stroke
// (outer top/bottom borders and band between 3×3 blocks).
func (m *model) gridRuleRowAllAccent() string {
	cw := m.cellContentWidth()
	var b strings.Builder
	b.WriteByte(gridX)
	for c := 0; c < 9; c++ {
		b.WriteString(strings.Repeat(string(gridH), cw))
		if c < 8 {
			b.WriteByte(gridX)
		}
	}
	b.WriteByte(gridX)
	return m.gridStrokeAccent().Render(b.String())
}

// gridRuleRowInnerSeparators uses the muted stroke; accent intersections on 3×3
// column boundaries and at the board corners.
func (m *model) gridRuleRowInnerSeparators() string {
	cw := m.cellContentWidth()
	parts := make([]string, 0, 22)
	muted := m.gridStrokeMuted()
	accent := m.gridStrokeAccent()
	parts = append(parts, accent.Render(string(gridX)))
	for c := 0; c < 9; c++ {
		parts = append(parts, muted.Render(strings.Repeat(string(gridH), cw)))
		if c < 8 {
			st := muted
			if c%3 == 2 {
				st = accent
			}
			parts = append(parts, st.Render(string(gridX)))
		}
	}
	parts = append(parts, accent.Render(string(gridX)))
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (m *model) gridTopBorder() string {
	return m.gridRuleRowAllAccent()
}

func (m *model) gridBottomBorder() string {
	return m.gridRuleRowAllAccent()
}

func (m *model) gridInnerSeparator() string {
	return m.gridRuleRowInnerSeparators()
}

func (m *model) gridBandSeparator() string {
	return m.gridRuleRowAllAccent()
}

// expandTitleTabs replaces tabs so rune counts match terminal cells. Lipgloss
// renders tabs as 4 spaces by default; a single '\t' in the const would otherwise
// throw off padding/centering vs len(line).
func expandTitleTabs(s string) string {
	return strings.ReplaceAll(s, "\t", strings.Repeat(" ", 4))
}

const sudokuTitleFallbackWord = "SUDOKU"

// renderSudokuTitleFallback is used when the ASCII banner is wider than the sidebar.
func renderSudokuTitleFallback(w int) string {
	word := sudokuTitleFallbackWord
	cols := lipgloss.Blend1D(len(word), tnTitle, tnGiven, tnLegendGrad)
	segs := make([]string, 0, len(word))
	for i, r := range word {
		bg := cols[i]
		if bg == nil {
			bg = tnTitle
		}
		segs = append(segs, lipgloss.NewStyle().
			Bold(true).
			Background(bg).
			Foreground(tnCursorFG).
			Padding(0, 1).
			Render(string(r)))
	}
	line := lipgloss.JoinHorizontal(lipgloss.Top, segs...)
	return lipgloss.PlaceHorizontal(w, lipgloss.Center, line)
}

// renderSudokuTitle styles the multi-line ASCII banner to fit sidebar width w.
func renderSudokuTitle(w int) string {
	if w <= titleFallbackMaxSidebarW {
		return renderSudokuTitleFallback(w)
	}
	// Do not use TrimSpace on the whole block: it removes leading spaces on line 1
	// (they follow the opening newline from the raw string const).
	art := strings.TrimPrefix(strings.TrimSuffix(sudokuTitle, "\n"), "\n")
	raw := strings.Split(art, "\n")
	lines := make([]string, 0, len(raw))
	for _, line := range raw {
		line = expandTitleTabs(line)
		// Do not TrimRight: trailing spaces are part of the column grid.
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return ""
	}
	// Use display cell width (not rune count) so padding matches the terminal.
	maxW := 0
	for _, line := range lines {
		if mw := lipgloss.Width(line); mw > maxW {
			maxW = mw
		}
	}
	for i := range lines {
		for lipgloss.Width(lines[i]) < maxW {
			lines[i] += " "
		}
	}
	// Vertical sweep: title purple → given blue → cyan (ties into legend border blend).
	palette := lipgloss.Blend1D(len(lines), tnTitle, tnGiven, tnLegendGrad)
	parts := make([]string, len(lines))
	for i, line := range lines {
		fg := palette[i]
		if fg == nil {
			fg = tnTitle
		}
		parts[i] = lipgloss.NewStyle().
			Bold(true).
			Foreground(fg).
			Render(line)
	}
	block := lipgloss.JoinVertical(lipgloss.Left, parts...)
	// Center the whole banner once; per-line Center + Width(w) misaligns rows when
	// lipgloss string width differs from len([]rune).
	return lipgloss.PlaceHorizontal(w, lipgloss.Center, block)
}

func (m *model) renderSidebar() string {
	modeStr := "VALUE"
	if m.mode == modePencil {
		modeStr = "PENCIL"
	}
	w := m.sidebarW
	if w > m.termW-4 && m.termW > 4 {
		w = m.termW / 3
	}
	if w < 24 {
		w = 24
	}
	title := renderSudokuTitle(w)
	modeLine := titleStyle.Width(w).Align(lipgloss.Center).Render(fmt.Sprintf("[%s mode]", modeStr))
	legendBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForegroundBlend(tnTitle, tnLegendGrad).
		Padding(1, 2).
		Width(w)
	innerW := w - legendBox.GetHorizontalFrameSize()
	if innerW < 8 {
		innerW = 8
	}
	legendKeyStyle := lipgloss.NewStyle().
		Foreground(tnGridAccent).
		Align(lipgloss.Left).
		Padding(0, 2, 0, 0)
	legendTbl := table.New().
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderHeader(false).
		BorderColumn(false).
		BorderRow(false).
		Width(innerW).
		Wrap(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return legendKeyStyle
			}
			return legendItemStyle
		}).
		Rows(
			[]string{"Move", "←↓↑→ or hjkl"},
			[]string{"Jump box right", "w"},
			[]string{"Jump box left", "W"},
			[]string{"Jump box up", "E"},
			[]string{"Jump box down", "e"},
			[]string{"toggle input mode", "[tab]"},
			[]string{"fill pencil marks", "f"},
			[]string{"scroll grid", "pgup/pgdn"},
			[]string{"", "mouse wheel"},
			[]string{"enter digit / pencil mark", "1-9"},
			[]string{"clear cell", "x / 0 / backspace"},
			[]string{"quit", "q / Ctrl+C"},
		)
	leg := legendBox.Render(legendTbl.String())
	return lipgloss.JoinVertical(lipgloss.Left, title, "", modeLine, "", leg)
}

func (m *model) View() tea.View {
	if m.termW <= 0 || m.termH <= 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}
	m.relayoutViewport()
	content := m.renderBoardFlat()
	m.vp.SetContent(content)
	m.syncViewportToCursor()

	side := m.renderSidebar()
	row := lipgloss.JoinHorizontal(lipgloss.Top, m.vp.View(), side)
	out := lipgloss.Place(m.termW, m.termH, lipgloss.Left, lipgloss.Top, row)
	v := tea.NewView(out)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
