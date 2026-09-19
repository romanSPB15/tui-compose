package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

const doubleClickWindow = 400 * time.Millisecond

type Picker struct {
	title string

	all      []string
	filtered []string

	query    []rune
	inputCur int

	cursor  int
	offset  int
	hovered int

	focused bool
	wnd     tui.Window

	lastClick    time.Time
	lastClickIdx int

	titleStyle, promptStyle, borderStyle, queryStyle, cursorStyle tui.Style
	itemStyle, activeStyle, hoverStyle, helpStyle, matchStyle     tui.Style

	result string
}

func NewPicker(items []string, title string) *Picker {
	sorted := make([]string, len(items))
	copy(sorted, items)
	sort.Strings(sorted)

	p := &Picker{
		title:        title,
		all:          sorted,
		hovered:      -1,
		lastClickIdx: -1,
		titleStyle:   tui.FrBrightWhite | tui.Bold,
		promptStyle:  tui.FrGreen | tui.Bold,
		borderStyle:  tui.FrBrightBlack,
		queryStyle:   tui.FrWhite,
		cursorStyle:  tui.BgWhite | tui.FrBlack,
		itemStyle:    tui.FrWhite,
		activeStyle:  tui.BgBlue | tui.FrWhite | tui.Bold,
		hoverStyle:   tui.BgBrightBlack | tui.FrWhite,
		helpStyle:    tui.FrBrightBlack,
		matchStyle:   tui.FrYellow | tui.Bold,
	}
	p.applyFilter()
	return p
}

func (p *Picker) Result() string { return p.result }

func (p *Picker) Width() int  { return p.wnd.Width() }
func (p *Picker) Height() int { return p.wnd.Height() }

func (p *Picker) effectiveQuery() string {
	return strings.ToLower(strings.TrimSpace(string(p.query)))
}

func (p *Picker) topOffset() int {
	if p.title != "" {
		return 1
	}
	return 0
}

func (p *Picker) filterY() int { return p.topOffset() }
func (p *Picker) listTop() int { return p.topOffset() + 2 }

func (p *Picker) applyFilter() {
	q := p.effectiveQuery()
	out := make([]string, 0, len(p.all))
	if q == "" {
		out = append(out, p.all...)
	} else {
		for _, s := range p.all {
			if strings.Contains(strings.ToLower(s), q) {
				out = append(out, s)
			}
		}
	}
	p.filtered = out
	p.cursor = 0
	p.offset = 0
	p.hovered = -1
}

func (p *Picker) listH() int {
	h := p.wnd.Height() - 4 - p.topOffset()
	if h < 1 {
		h = 1
	}
	return h
}

func (p *Picker) clampOffset() {
	lh := p.listH()
	if p.cursor < p.offset {
		p.offset = p.cursor
	}
	if p.cursor >= p.offset+lh {
		p.offset = p.cursor - lh + 1
	}
	maxOff := len(p.filtered) - lh
	if maxOff < 0 {
		maxOff = 0
	}
	if p.offset > maxOff {
		p.offset = maxOff
	}
	if p.offset < 0 {
		p.offset = 0
	}
}

func (p *Picker) rowToIndex(y int) int {
	top := p.listTop()
	lh := p.listH()
	if y < top || y >= top+lh {
		return -1
	}
	idx := p.offset + (y - top)
	if idx < 0 || idx >= len(p.filtered) {
		return -1
	}
	return idx
}

func (p *Picker) matchRange(text string) (int, int) {
	q := p.effectiveQuery()
	if q == "" {
		return -1, -1
	}
	lower := strings.ToLower(text)
	idx := strings.Index(lower, q)
	if idx < 0 {
		return -1, -1
	}
	start := utf8.RuneCountInString(text[:idx])
	end := start + utf8.RuneCountInString(q)
	return start, end
}

func (p *Picker) Render(buf [][]cell.Cell) {
	h := len(buf)
	if h == 0 {
		return
	}
	w := len(buf[0])
	if w == 0 {
		return
	}

	empty := cell.Style{}
	ts := tui.ConvertToCellStyle(p.titleStyle)
	bs := tui.ConvertToCellStyle(p.borderStyle)
	ps := tui.ConvertToCellStyle(p.promptStyle)
	qs := tui.ConvertToCellStyle(p.queryStyle)
	cs := tui.ConvertToCellStyle(p.cursorStyle)
	its := tui.ConvertToCellStyle(p.itemStyle)
	as := tui.ConvertToCellStyle(p.activeStyle)
	hvs := tui.ConvertToCellStyle(p.hoverStyle)
	hs := tui.ConvertToCellStyle(p.helpStyle)
	ms := tui.ConvertToCellStyle(p.matchStyle)

	for y := 0; y < h; y++ {
		row := buf[y]
		for x := 0; x < w && x < len(row); x++ {
			row[x] = cell.Cell{Char: ' ', Style: empty}
		}
	}

	filterY := p.filterY()

	if p.title != "" {
		titleRunes := []rune(p.title)
		if len(titleRunes) > w {
			titleRunes = titleRunes[:w]
		}
		start := (w - len(titleRunes)) / 2
		if start < 0 {
			start = 0
		}
		for i, r := range titleRunes {
			x := start + i
			if x < w {
				buf[0][x] = cell.Cell{Char: r, Style: ts}
			}
		}
	}

	prompt := "> "
	px := 0
	for _, r := range prompt {
		if px >= w {
			break
		}
		buf[filterY][px] = cell.Cell{Char: r, Style: ps}
		px++
	}
	for i, r := range p.query {
		x := px + i
		if x >= w {
			break
		}
		buf[filterY][x] = cell.Cell{Char: r, Style: qs}
	}
	if p.focused {
		cx := px + p.inputCur
		if cx < w {
			var ch rune = ' '
			if p.inputCur < len(p.query) {
				ch = p.query[p.inputCur]
			}
			buf[filterY][cx] = cell.Cell{Char: ch, Style: cs}
		}
	}

	div1Y := filterY + 1
	if div1Y < h {
		for x := 0; x < w && x < len(buf[div1Y]); x++ {
			buf[div1Y][x] = cell.Cell{Char: '-', Style: bs}
		}
	}
	if h >= 3 {
		div2Y := h - 2
		for x := 0; x < w && x < len(buf[div2Y]); x++ {
			buf[div2Y][x] = cell.Cell{Char: '-', Style: bs}
		}
	}

	top := p.listTop()
	lh := p.listH()
	if lh > 0 {
		for i := 0; i < lh && p.offset+i < len(p.filtered); i++ {
			idx := p.offset + i
			row := top + i
			if row >= h {
				break
			}
			base := its
			prefix := "[ ]  "
			switch {
			case idx == p.cursor:
				base = as
				prefix = "[x] "
			case idx == p.hovered:
				base = hvs
			}

			for x := 0; x < w && x < len(buf[row]); x++ {
				buf[row][x] = cell.Cell{Char: ' ', Style: base}
			}
			for j, r := range prefix {
				if j < w {
					buf[row][j] = cell.Cell{Char: r, Style: base}
				}
			}

			text := p.filtered[idx]
			runes := []rune(text)
			mStart, mEnd := p.matchRange(text)

			for j, r := range runes {
				x := 4 + j
				if x >= w {
					break
				}
				st := base
				if mStart >= 0 && j >= mStart && j < mEnd {
					st = base.Merge(ms)
				}
				buf[row][x] = cell.Cell{Char: r, Style: st}
			}
		}
	}

	if len(p.filtered) == 0 && h >= 5 {
		msg := "Ничего не найдено"
		runes := []rune(msg)
		row := top + lh/2
		if row >= h {
			row = h - 1
		}
		start := (w - len(runes)) / 2
		for i, r := range runes {
			x := start + i
			if x >= 0 && x < w && x < len(buf[row]) {
				buf[row][x] = cell.Cell{Char: r, Style: hs}
			}
		}
	}

	row := h - 1
	help := " ENTER, двойной клик - выбрать | стрелки, колесо мыши - навигация | ESC - отмена"
	for i, r := range []rune(help) {
		if i < w && i < len(buf[row]) {
			buf[row][i] = cell.Cell{Char: r, Style: hs}
		}
	}
}

func (p *Picker) Send(ev tui.Event) {
	switch e := ev.(type) {
	case *tui.WindowEvent:
		p.wnd = e.Window
	case *tui.CheckFocusableEvent:
		e.Result = true
	case *tui.FocusEvent:
		p.focused = e.Focused
	case *input.KeyboardEvent:
		p.handleKey(e)
	case *input.MouseEvent:
		p.handleMouse(e)
	case *tui.MouseHoverEvent:
		if !e.Entered {
			if p.hovered != -1 {
				p.hovered = -1
				p.wnd.Redraw()
			}
		}
	}
}

func (p *Picker) handleKey(e *input.KeyboardEvent) {
	switch e.Key {
	case input.KeyEsc:
		p.result = ""
		p.wnd.Quit()
		return
	case input.KeyEnter:
		if p.cursor >= 0 && p.cursor < len(p.filtered) {
			p.result = p.filtered[p.cursor]
		}
		p.wnd.Quit()
		return
	case input.KeyArrowUp:
		if p.cursor > 0 {
			p.cursor--
		}
		p.clampOffset()
	case input.KeyArrowDown:
		if p.cursor < len(p.filtered)-1 {
			p.cursor++
		}
		p.clampOffset()
	case input.KeyHome:
		p.cursor = 0
		p.clampOffset()
	case input.KeyEnd:
		if len(p.filtered) > 0 {
			p.cursor = len(p.filtered) - 1
		}
		p.clampOffset()
	case input.KeyPgUp:
		p.cursor -= p.listH()
		if p.cursor < 0 {
			p.cursor = 0
		}
		p.clampOffset()
	case input.KeyPgDown:
		p.cursor += p.listH()
		if p.cursor >= len(p.filtered) {
			p.cursor = len(p.filtered) - 1
			if p.cursor < 0 {
				p.cursor = 0
			}
		}
		p.clampOffset()
	case input.KeyBackspace:
		if p.inputCur > 0 {
			p.query = append(p.query[:p.inputCur-1], p.query[p.inputCur:]...)
			p.inputCur--
			p.applyFilter()
			p.clampOffset()
		}
	case input.KeyDelete:
		if p.inputCur < len(p.query) {
			p.query = append(p.query[:p.inputCur], p.query[p.inputCur+1:]...)
			p.applyFilter()
			p.clampOffset()
		}
	case input.KeyArrowLeft:
		if p.inputCur > 0 {
			p.inputCur--
		}
	case input.KeyArrowRight:
		if p.inputCur < len(p.query) {
			p.inputCur++
		}
	case input.KeyCtrlU:
		p.query = p.query[:0]
		p.inputCur = 0
		p.applyFilter()
		p.clampOffset()
	case input.KeyCtrlW:
		for p.inputCur > 0 && p.query[p.inputCur-1] == ' ' {
			p.query = append(p.query[:p.inputCur-1], p.query[p.inputCur:]...)
			p.inputCur--
		}
		for p.inputCur > 0 && p.query[p.inputCur-1] != ' ' {
			p.query = append(p.query[:p.inputCur-1], p.query[p.inputCur:]...)
			p.inputCur--
		}
		p.applyFilter()
		p.clampOffset()
	default:
		if e.Rune != 0 {
			p.query = append(p.query[:p.inputCur],
				append([]rune{e.Rune}, p.query[p.inputCur:]...)...)
			p.inputCur++
			p.applyFilter()
			p.clampOffset()
		}
	}
	p.wnd.Redraw()
}

func (p *Picker) handleMouse(e *input.MouseEvent) {
	switch e.Action {
	case input.MouseWheelUp:
		p.cursor--
		if p.cursor < 0 {
			p.cursor = 0
		}
		p.clampOffset()
		p.wnd.Redraw()

	case input.MouseWheelDown:
		p.cursor++
		if p.cursor >= len(p.filtered) {
			p.cursor = len(p.filtered) - 1
			if p.cursor < 0 {
				p.cursor = 0
			}
		}
		p.clampOffset()
		p.wnd.Redraw()

	case input.MouseMove:
		idx := p.rowToIndex(e.Pos.Y)
		if idx != p.hovered {
			p.hovered = idx
			p.wnd.Redraw()
		}

	case input.MousePress:
		idx := p.rowToIndex(e.Pos.Y)
		if idx < 0 {
			return
		}
		now := time.Now()
		if p.lastClickIdx == idx && now.Sub(p.lastClick) < doubleClickWindow {
			p.result = p.filtered[idx]
			p.wnd.Quit()
			return
		}
		p.cursor = idx
		p.lastClick = now
		p.lastClickIdx = idx
		p.clampOffset()
		p.wnd.Redraw()
	}
}

func openTTY() (*os.File, *os.File, error) {
	if runtime.GOOS == "windows" {
		in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
		if err != nil {
			return nil, nil, err
		}
		out, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
		if err != nil {
			in.Close()
			return nil, nil, err
		}
		return in, out, nil
	}
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}
	return tty, tty, nil
}

func main() {
	code, _ := run()
	os.Exit(code)
}

func run() (int, string) {
	var title string
	flag.StringVar(&title, "t", "", "заголовок")
	flag.StringVar(&title, "title", "", "заголовок")
	flag.Parse()

	var items []string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		for sc.Scan() {
			line := sc.Text()
			if line != "" {
				items = append(items, line)
			}
		}
	} else {
		items = flag.Args()
	}

	if len(items) == 0 {
		fmt.Fprintln(os.Stderr, "tpick: нет элементов")
		return 1, ""
	}

	origStdout := os.Stdout

	ttyIn, ttyOut, err := openTTY()
	if err != nil {
		ttyIn, ttyOut = os.Stdin, os.Stdout
	}

	wnd := tui.NewWindow(tui.WithIO(ttyIn, ttyOut, ttyOut))
	wnd.SetTitle("tpick")

	picker := NewPicker(items, title)
	wnd.SetContent(picker)
	wnd.Focus().SetIndex(0)
	wnd.Run()

	if r := picker.Result(); r != "" {
		fmt.Fprintln(origStdout, r)
		return 0, ""
	}
	return 1, ""
}
