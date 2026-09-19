package extra

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

const tabWidth = 4

// TextArea — многострочный текстовый редактор с переносом строк.
//
// Стрелки, Home/End, PgUp/PgDn, Enter, Backspace/Delete, Insert,
// вертикальный скролл (клавиши/колесо), мигающий курсор,
// опциональные номера строк и автоматический перенос длинных строк.
type TextArea struct {
	lines      []string
	cursorLine int
	cursorCol  int

	width  int
	height int

	focused   bool
	hovered   bool
	overwrite bool

	style, styleF, styleH cell.Style
	cursorStyle           cell.Style
	lineNumStyle          cell.Style

	showLineNumbers bool

	offsetTop int

	wnd tui.Window

	blinkEnabled  bool
	blinkInterval time.Duration
	cursorVisible bool
	blinkOnce     sync.Once

	OnChanged func(string)
}

// NewTextArea создаёт поле шириной width и высотой height строк.
func NewTextArea(width, height int) *TextArea {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return &TextArea{
		lines:         []string{""},
		width:         width,
		height:        height,
		cursorStyle:   cell.Style{Bg: "47", Fg: "34"},
		lineNumStyle:  cell.Style{Fg: "90"},
		cursorVisible: true,
		blinkEnabled:  true,
		blinkInterval: 500 * time.Millisecond,
	}
}

// WithStyle задаёт стиль текста без фокуса.
func (ta *TextArea) WithStyle(s tui.Style) *TextArea {
	ta.style = tui.ConvertToCellStyle(s)
	return ta
}

// WithFocusedStyle задаёт стиль текста в фокусе.
func (ta *TextArea) WithFocusedStyle(s tui.Style) *TextArea {
	ta.styleF = tui.ConvertToCellStyle(s)
	return ta
}

// WithHoverStyle задаёт стиль при наведении.
func (ta *TextArea) WithHoverStyle(s tui.Style) *TextArea {
	ta.styleH = tui.ConvertToCellStyle(s)
	return ta
}

// WithCursorStyle задаёт стиль курсора.
func (ta *TextArea) WithCursorStyle(s tui.Style) *TextArea {
	ta.cursorStyle = tui.ConvertToCellStyle(s)
	return ta
}

// WithLineNumberStyle задаёт стиль номеров строк.
func (ta *TextArea) WithLineNumberStyle(s tui.Style) *TextArea {
	ta.lineNumStyle = tui.ConvertToCellStyle(s)
	return ta
}

// WithLineNumbers включает или отключает номера строк.
func (ta *TextArea) WithLineNumbers(v bool) *TextArea {
	ta.showLineNumbers = v
	return ta
}

// WithBlinkInterval задаёт интервал мигания курсора.
func (ta *TextArea) WithBlinkInterval(d time.Duration) *TextArea {
	ta.blinkInterval = d
	return ta
}

// DisableBlink отключает мигание.
func (ta *TextArea) DisableBlink() *TextArea {
	ta.blinkEnabled = false
	ta.cursorVisible = true
	return ta
}

// EnableBlink включает мигание.
func (ta *TextArea) EnableBlink() *TextArea {
	ta.blinkEnabled = true
	return ta
}

// WithOverwrite включает или выключает режим затирания.
func (ta *TextArea) WithOverwrite(v bool) *TextArea {
	ta.overwrite = v
	return ta
}

// WithOnChanged задаёт обработчик изменения текста.
func (ta *TextArea) WithOnChanged(h func(string)) *TextArea {
	ta.OnChanged = h
	return ta
}

// WithText заменяет содержимое текстом, разбитым по \n.
func (ta *TextArea) WithText(text string) *TextArea {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	ta.lines = lines
	ta.cursorLine = 0
	ta.cursorCol = 0
	ta.offsetTop = 0
	rows := ta.buildVisualRows()
	ta.clampOffsets(rows)
	return ta
}

// Text возвращает всё содержимое через \n.
func (ta *TextArea) Text() string {
	return strings.Join(ta.lines, "\n")
}

// Cursor возвращает логическую позицию курсора.
func (ta *TextArea) Cursor() (line, col int) {
	return ta.cursorLine, ta.cursorCol
}

// Width возвращает полную ширину, включая номера строк.
func (ta *TextArea) Width() int {
	if ta.showLineNumbers {
		return ta.lineNumWidth() + 1 + ta.width
	}
	return ta.width
}

// Height возвращает высоту в строках.
func (ta *TextArea) Height() int {
	return ta.height
}

func (ta *TextArea) lineNumWidth() int {
	n := len(ta.lines)
	if n < 1 {
		n = 1
	}
	return len(strconv.Itoa(n))
}

// ---------- визуальные строки ----------

// lineData хранит expanded-руны строки и маппинги между
// логическими и визуальными колонками.
type lineData struct {
	expanded []rune
	vToL     []int // expanded-позиция → логическая колонка (длина len(expanded)+1)
	lToV     []int // логическая колонка → expanded-позиция (длина len(logical)+1)
}

type visualRow struct {
	line   int
	startE int
	endE   int
	data   *lineData
}

// makeLineData раскрывает табы и строит маппинги.
func makeLineData(line string) *lineData {
	logical := []rune(line)
	expanded := make([]rune, 0, len(logical))
	vToL := make([]int, 0, len(logical)+1)
	lToV := make([]int, 0, len(logical)+1)

	col := 0
	for i, r := range logical {
		lToV = append(lToV, col)
		if r == '\t' {
			n := tabWidth - col%tabWidth
			for j := 0; j < n; j++ {
				expanded = append(expanded, ' ')
				vToL = append(vToL, i)
			}
			col += n
		} else {
			expanded = append(expanded, r)
			vToL = append(vToL, i)
			col++
		}
	}
	vToL = append(vToL, len(logical))
	lToV = append(lToV, col)

	return &lineData{expanded: expanded, vToL: vToL, lToV: lToV}
}

// buildVisualRows строит визуальные строки с переносом по ta.width.
func (ta *TextArea) buildVisualRows() []visualRow {
	w := ta.width
	if w < 1 {
		w = 1
	}
	rows := make([]visualRow, 0, len(ta.lines))
	for i, line := range ta.lines {
		data := makeLineData(line)
		if len(data.expanded) == 0 {
			rows = append(rows, visualRow{line: i, startE: 0, endE: 0, data: data})
			continue
		}
		for start := 0; start < len(data.expanded); start += w {
			end := start + w
			if end > len(data.expanded) {
				end = len(data.expanded)
			}
			rows = append(rows, visualRow{line: i, startE: start, endE: end, data: data})
		}
	}
	return rows
}

// findCursorVisual находит визуальный ряд и колонку курсора.
func findCursorVisual(rows []visualRow, line, col int) (vr int, vcol int) {
	for i := len(rows) - 1; i >= 0; i-- {
		r := rows[i]
		if r.line != line {
			continue
		}
		if col >= len(r.data.lToV) {
			col = len(r.data.lToV) - 1
		}
		vc := r.data.lToV[col]
		if vc >= r.startE {
			return i, vc - r.startE
		}
	}
	return -1, 0
}

// setCursorFromVisual устанавливает курсор по визуальной позиции.
func (ta *TextArea) setCursorFromVisual(rows []visualRow, vrIdx, vcol int) {
	if vrIdx < 0 || vrIdx >= len(rows) {
		return
	}
	r := rows[vrIdx]
	vc := r.startE + vcol
	if vc < 0 {
		vc = 0
	}
	if vc >= len(r.data.vToL) {
		vc = len(r.data.vToL) - 1
	}
	ta.cursorLine = r.line
	ta.cursorCol = r.data.vToL[vc]
}

// clampOffsets держит курсор в видимой области.
func (ta *TextArea) clampOffsets(rows []visualRow) {
	vr, _ := findCursorVisual(rows, ta.cursorLine, ta.cursorCol)
	if vr < 0 {
		return
	}
	if vr < ta.offsetTop {
		ta.offsetTop = vr
	}
	if vr >= ta.offsetTop+ta.height {
		ta.offsetTop = vr - ta.height + 1
	}
	if ta.offsetTop < 0 {
		ta.offsetTop = 0
	}
	maxTop := len(rows) - ta.height
	if maxTop < 0 {
		maxTop = 0
	}
	if ta.offsetTop > maxTop {
		ta.offsetTop = maxTop
	}
}

// clampOffsetOnly подрезает offsetTop, не ориентируясь на курсор.
// Используется в Render, чтобы не возвращать прокрутку назад.
func (ta *TextArea) clampOffsetOnly(totalRows int) {
	maxTop := totalRows - ta.height
	if maxTop < 0 {
		maxTop = 0
	}
	if ta.offsetTop > maxTop {
		ta.offsetTop = maxTop
	}
	if ta.offsetTop < 0 {
		ta.offsetTop = 0
	}
}

// ---------- рендер ----------

// Render рисует поле в буфер.
func (ta *TextArea) Render(buf [][]cell.Cell) {
	if len(buf) == 0 {
		return
	}

	fieldStyle := ta.style
	switch {
	case ta.focused:
		fieldStyle = ta.styleF
	case ta.hovered && ta.styleH != (cell.Style{}):
		fieldStyle = ta.styleH
	}

	rows := ta.buildVisualRows()
	ta.clampOffsetOnly(len(rows))

	numWidth := ta.lineNumWidth()
	xStart := 0
	if ta.showLineNumbers {
		xStart = numWidth + 1
	}

	textWidth := ta.width
	totalW := ta.Width()
	if totalW > len(buf[0]) {
		totalW = len(buf[0])
	}

	for row := 0; row < ta.height && row < len(buf); row++ {
		for x := 0; x < totalW; x++ {
			buf[row][x] = cell.Cell{Char: ' ', Style: fieldStyle}
		}

		vr := ta.offsetTop + row
		if vr >= len(rows) {
			continue
		}
		r := rows[vr]

		// номера строк — только для первого сегмента логической строки
		if ta.showLineNumbers && r.startE == 0 {
			numStr := strconv.Itoa(r.line + 1)
			pad := numWidth - len(numStr)
			for i := 0; i < pad; i++ {
				buf[row][i] = cell.Cell{Char: ' ', Style: ta.lineNumStyle}
			}
			for i, ch := range numStr {
				buf[row][pad+i] = cell.Cell{Char: ch, Style: ta.lineNumStyle}
			}
			buf[row][numWidth] = cell.Cell{Char: ' ', Style: ta.lineNumStyle}
		}

		// сегмент текста
		seg := r.data.expanded[r.startE:r.endE]
		for i := 0; i < len(seg) && i < textWidth; i++ {
			x := xStart + i
			if x >= totalW {
				break
			}
			buf[row][x] = cell.Cell{Char: seg[i], Style: fieldStyle}
		}

		// курсор
		if ta.focused && ta.cursorVisible && r.line == ta.cursorLine {
			if cVr, cVc := findCursorVisual(rows, ta.cursorLine, ta.cursorCol); cVr == vr {
				cursorX := xStart + cVc
				if cursorX < xStart+textWidth && cursorX < totalW {
					var ch rune = ' '
					if cVc < len(seg) {
						ch = seg[cVc]
					}
					under := buf[row][cursorX]
					buf[row][cursorX] = cell.Cell{
						Char:  ch,
						Style: ta.cursorStyle.Merge(under.Style),
					}
				}
			}
		}
	}
}

// ---------- события ----------

// Send обрабатывает событие.
func (ta *TextArea) Send(ev tui.Event) {
	switch e := ev.(type) {
	case *tui.WindowEvent:
		ta.wnd = e.Window
	case *tui.CheckFocusableEvent:
		e.Result = true
	case *tui.FocusEvent:
		ta.focused = e.Focused
		if e.Focused {
			ta.cursorVisible = true
			ta.startBlink()
		}
	case *tui.MouseHoverEvent:
		ta.hovered = e.Entered
		if ta.wnd != nil {
			ta.wnd.Redraw()
		}
	case *input.KeyboardEvent:
		ta.handleKey(e)
	case *input.MouseEvent:
		switch e.Action {
		case input.MousePress:
			ta.handleClick(e)
		case input.MouseWheelUp:
			ta.scrollBy(-1)
		case input.MouseWheelDown:
			ta.scrollBy(1)
		}
	}
}

func (ta *TextArea) scrollBy(delta int) {
	rows := ta.buildVisualRows()
	ta.offsetTop += delta
	if ta.offsetTop < 0 {
		ta.offsetTop = 0
	}
	maxTop := len(rows) - ta.height
	if maxTop < 0 {
		maxTop = 0
	}
	if ta.offsetTop > maxTop {
		ta.offsetTop = maxTop
	}
	if ta.wnd != nil {
		ta.wnd.Redraw()
	}
}

func (ta *TextArea) handleKey(ev *input.KeyboardEvent) {
	line := []rune(ta.lines[ta.cursorLine])
	changed := false
	rows := ta.buildVisualRows()

	switch ev.Key {
	case input.KeyArrowUp:
		vr, vc := findCursorVisual(rows, ta.cursorLine, ta.cursorCol)
		if vr > 0 {
			ta.setCursorFromVisual(rows, vr-1, vc)
		}
	case input.KeyArrowDown:
		vr, vc := findCursorVisual(rows, ta.cursorLine, ta.cursorCol)
		if vr >= 0 && vr < len(rows)-1 {
			ta.setCursorFromVisual(rows, vr+1, vc)
		}
	case input.KeyArrowLeft:
		if ta.cursorCol > 0 {
			ta.cursorCol--
		} else if ta.cursorLine > 0 {
			ta.cursorLine--
			ta.cursorCol = len([]rune(ta.lines[ta.cursorLine]))
		}
	case input.KeyArrowRight:
		if ta.cursorCol < len(line) {
			ta.cursorCol++
		} else if ta.cursorLine < len(ta.lines)-1 {
			ta.cursorLine++
			ta.cursorCol = 0
		}
	case input.KeyPgUp:
		vr, vc := findCursorVisual(rows, ta.cursorLine, ta.cursorCol)
		vr -= ta.height
		if vr < 0 {
			vr = 0
		}
		ta.setCursorFromVisual(rows, vr, vc)
	case input.KeyPgDown:
		vr, vc := findCursorVisual(rows, ta.cursorLine, ta.cursorCol)
		vr += ta.height
		if vr >= len(rows) {
			vr = len(rows) - 1
		}
		ta.setCursorFromVisual(rows, vr, vc)
	case input.KeyHome:
		ta.cursorCol = 0
	case input.KeyEnd:
		ta.cursorCol = len(line)
	case input.KeyBackspace:
		if ta.cursorCol > 0 {
			line = append(line[:ta.cursorCol-1], line[ta.cursorCol:]...)
			ta.lines[ta.cursorLine] = string(line)
			ta.cursorCol--
			changed = true
		} else if ta.cursorLine > 0 {
			prevLen := len([]rune(ta.lines[ta.cursorLine-1]))
			ta.lines[ta.cursorLine-1] += string(line)
			ta.lines = append(ta.lines[:ta.cursorLine], ta.lines[ta.cursorLine+1:]...)
			ta.cursorLine--
			ta.cursorCol = prevLen
			changed = true
		}
	case input.KeyDelete:
		if ta.cursorCol < len(line) {
			line = append(line[:ta.cursorCol], line[ta.cursorCol+1:]...)
			ta.lines[ta.cursorLine] = string(line)
			changed = true
		} else if ta.cursorLine < len(ta.lines)-1 {
			ta.lines[ta.cursorLine] += ta.lines[ta.cursorLine+1]
			ta.lines = append(ta.lines[:ta.cursorLine+1], ta.lines[ta.cursorLine+2:]...)
			changed = true
		}
	case input.KeyEnter:
		before := string(line[:ta.cursorCol])
		after := string(line[ta.cursorCol:])
		ta.lines[ta.cursorLine] = before
		ta.lines = append(ta.lines[:ta.cursorLine+1],
			append([]string{after}, ta.lines[ta.cursorLine+1:]...)...)
		ta.cursorLine++
		ta.cursorCol = 0
		changed = true
	case input.KeyInsert:
		ta.overwrite = !ta.overwrite
	default:
		if ev.Rune != 0 {
			if ta.overwrite && ta.cursorCol < len(line) {
				line[ta.cursorCol] = ev.Rune
				ta.lines[ta.cursorLine] = string(line)
				ta.cursorCol++
			} else {
				line = append(line[:ta.cursorCol],
					append([]rune{ev.Rune}, line[ta.cursorCol:]...)...)
				ta.lines[ta.cursorLine] = string(line)
				ta.cursorCol++
			}
			changed = true
		}
	}

	ta.cursorVisible = true
	rows = ta.buildVisualRows()
	ta.clampOffsets(rows)
	if ta.wnd != nil {
		ta.wnd.Redraw()
	}
	if changed && ta.OnChanged != nil {
		ta.OnChanged(ta.Text())
	}
}

func (ta *TextArea) handleClick(ev *input.MouseEvent) {
	row := ev.Pos.Y
	if row < 0 || row >= ta.height {
		return
	}
	rows := ta.buildVisualRows()
	vr := ta.offsetTop + row
	if vr < 0 || vr >= len(rows) {
		return
	}

	xStart := 0
	if ta.showLineNumbers {
		xStart = ta.lineNumWidth() + 1
	}
	vc := ev.Pos.X - xStart
	if vc < 0 {
		vc = 0
	}

	ta.setCursorFromVisual(rows, vr, vc)
	ta.cursorVisible = true
	ta.clampOffsets(rows)

	if ta.wnd != nil {
		ta.wnd.Focus().SetFocus(ta)
	}
}

func (ta *TextArea) startBlink() {
	if ta.wnd == nil {
		return
	}
	ta.blinkOnce.Do(func() {
		go ta.blinkLoop()
	})
}

func (ta *TextArea) blinkLoop() {
	ticker := time.NewTicker(ta.blinkInterval)
	defer ticker.Stop()

	quit := ta.wnd.OnQuit()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			ta.wnd.Do(func() {
				if ta.focused && ta.blinkEnabled {
					ta.cursorVisible = !ta.cursorVisible
					ta.wnd.Redraw()
				}
			})
		}
	}
}
