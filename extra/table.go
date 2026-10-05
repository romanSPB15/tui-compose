package extra

import (
	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
)

// TableCellAlign — выравнивание текста в ячейке таблицы.
type TableCellAlign int

const (
	AlignLeft   = iota
	AlignCenter = iota
	AlignRight  = iota
)

// BorderMode — режим отрисовки границ таблицы.
type BorderMode int

const (
	// BorderNone — без границ, только текст ячеек.
	BorderNone BorderMode = iota
	// BorderOuter — только внешняя рамка.
	BorderOuter
	// BorderInner — только внутренние разделители между строками.
	BorderInner
	// BorderAll — внешняя рамка и разделители между строками.
	BorderAll
)

// TableCell — одна ячейка таблицы.
type TableCell struct {
	Style tui.Style
	Text  string
	Align TableCellAlign
}

// TableStyle — символы рамки.
type TableStyle struct {
	Hor, Ver                                       rune
	TopRight, TopLeft, BottomRight, BottomLeft     rune
	Plus, PlusLeft, PlusRight, PlusBottom, PlusTop rune
}

// Table — виджет таблицы.
type Table struct {
	data      [][]TableCell
	widths    []int
	border    BorderMode
	s         TableStyle
	rowStyleF func(i int) acell.Style
	colStyleF func(i int) acell.Style
}

// NewTable создаёт таблицу с указанными ячейками.
func NewTable(cells [][]TableCell) *Table {
	t := &Table{data: cells, border: BorderAll}
	if len(cells) > 0 {
		t.widths = make([]int, len(cells[0]))
	}
	t.recalcWidths()
	t.Default()
	return t
}

func (t *Table) recalcWidths() {
	for i := range t.widths {
		t.widths[i] = 0
	}
	for _, row := range t.data {
		for j, cell := range row {
			if j >= len(t.widths) {
				break
			}
			if n := len([]rune(cell.Text)); n > t.widths[j] {
				t.widths[j] = n
			}
		}
	}
}

// Default устанавливает одинарные линии.
func (t *Table) Default() *Table {
	t.s = TableStyle{
		Hor: '─', Ver: '│',
		TopLeft: '┌', TopRight: '┐',
		BottomLeft: '└', BottomRight: '┘',
		Plus: '┼', PlusLeft: '├', PlusRight: '┤',
		PlusTop: '┬', PlusBottom: '┴',
	}
	return t
}

// Rounded устанавливает скруглённые углы.
func (t *Table) Rounded() *Table {
	t.s = TableStyle{
		Hor: '─', Ver: '│',
		TopLeft: '╭', TopRight: '╮',
		BottomLeft: '╰', BottomRight: '╯',
		Plus: '┼', PlusLeft: '├', PlusRight: '┤',
		PlusTop: '┬', PlusBottom: '┴',
	}
	return t
}

// ASCII устанавливает ASCII-символы.
func (t *Table) ASCII() *Table {
	t.s = TableStyle{
		Hor: '-', Ver: '|',
		TopLeft: '+', TopRight: '+',
		BottomLeft: '+', BottomRight: '+',
		Plus: '+', PlusLeft: '+', PlusRight: '+',
		PlusTop: '+', PlusBottom: '+',
	}
	return t
}

// WithStyle задаёт произвольный стиль рамки.
func (t *Table) WithStyle(s TableStyle) *Table {
	t.s = s
	return t
}

// WithBorder задаёт режим границ.
func (t *Table) WithBorder(b BorderMode) *Table {
	t.border = b
	return t
}

// WithRowStyle задаёт стиль ячеек строки i, мержится с их собственным стилем.
func (t *Table) WithRowStyle(fn func(i int) acell.Style) *Table {
	t.rowStyleF = fn
	return t
}

// WithColStyle задаёт стиль ячеек колонки i, мержится с их собственным стилем.
func (t *Table) WithColStyle(fn func(i int) acell.Style) *Table {
	t.colStyleF = fn
	return t
}

// WithData заменяет содержимое таблицы и пересчитывает ширины колонок.
func (t *Table) WithData(tc [][]TableCell) *Table {
	t.data = tc
	if len(tc) > 0 {
		if len(t.widths) < len(tc[0]) {
			for range len(tc[0]) - len(t.widths) {
				t.widths = append(t.widths, 0)
			}
		}
		t.recalcWidths()
	}
	return t
}

// Width возвращает ширину таблицы в ячейках.
func (t *Table) Width() int {
	if len(t.widths) == 0 {
		return 0
	}
	w := 0
	for _, x := range t.widths {
		w += x
	}
	cols := len(t.widths)
	switch t.border {
	case BorderNone:
		return w + 2*cols
	case BorderOuter:
		return w + 2*cols + 2
	default:
		return w + 3*cols + 1
	}
}

// Height возвращает высоту таблицы в ячейках.
func (t *Table) Height() int {
	rows := len(t.data)
	switch t.border {
	case BorderNone:
		return rows
	case BorderOuter:
		return rows + 2
	case BorderInner:
		return 2*rows - 1
	case BorderAll:
		return 2*rows + 1
	}
	return rows
}

// Render рисует таблицу в буфер.
func (t *Table) Render(buf [][]acell.Cell) {
	if len(t.data) == 0 || len(t.widths) == 0 {
		return
	}
	if len(buf) == 0 || len(buf[0]) == 0 {
		return
	}

	cols := len(t.widths)
	rows := len(t.data)

	hasInnerV := t.border == BorderInner || t.border == BorderAll
	hasOuter := t.border == BorderOuter || t.border == BorderAll
	hasInnerH := t.border == BorderInner || t.border == BorderAll
	hasVer := t.border != BorderNone

	contentX := make([]int, cols)
	totalW := 0
	switch t.border {
	case BorderNone:
		x := 0
		for i := 0; i < cols; i++ {
			contentX[i] = x + 1
			x += t.widths[i] + 2
		}
		totalW = x
	case BorderOuter:
		x := 1
		for i := 0; i < cols; i++ {
			contentX[i] = x + 1
			x += t.widths[i] + 2
		}
		totalW = x + 1
	default:
		x := 1
		for i := 0; i < cols; i++ {
			contentX[i] = x + 1
			x += t.widths[i] + 3
		}
		totalW = x
	}

	totalH := t.Height()

	rowY := make([]int, rows)
	switch t.border {
	case BorderNone:
		for i := range rowY {
			rowY[i] = i
		}
	case BorderOuter:
		for i := range rowY {
			rowY[i] = i + 1
		}
	case BorderInner:
		for i := range rowY {
			rowY[i] = 2 * i
		}
	case BorderAll:
		for i := range rowY {
			rowY[i] = 2*i + 1
		}
	}
	_ = totalH

	set := func(y, x int, c acell.Cell) {
		buf[y][x] = c
	}

	for i, row := range t.data {
		y := rowY[i]
		for j, cell := range row {
			if j >= cols {
				break
			}
			st := tui.ConvertToCellStyle(cell.Style)
			if t.rowStyleF != nil {
				st = st.Merge(t.rowStyleF(i))
			}
			if t.colStyleF != nil {
				st = st.Merge(t.colStyleF(j))
			}
			runes := []rune(cell.Text)
			n := len(runes)
			colW := t.widths[j]
			if n > colW {
				n = colW
			}
			var lPad int
			switch cell.Align {
			case AlignLeft:
				lPad = 0
			case AlignRight:
				lPad = colW - n
			case AlignCenter:
				lPad = (colW - n) / 2
			}

			startX := contentX[j] - 1
			for k := 0; k < colW+2; k++ {
				set(y, startX+k, acell.Cell{Char: ' ', Style: st})
			}

			cx := contentX[j] + lPad
			for k := 0; k < n; k++ {
				set(y, cx+k, acell.Cell{Char: runes[k], Style: st})
			}
		}
	}

	if hasOuter {
		s := acell.Style{}
		set(0, 0, acell.Cell{Char: t.s.TopLeft, Style: s})
		x := 1
		for j := 0; j < cols; j++ {
			for k := 0; k < t.widths[j]+2; k++ {
				set(0, x, acell.Cell{Char: t.s.Hor, Style: s})
				x++
			}
			if hasInnerV && j < cols-1 {
				set(0, x, acell.Cell{Char: t.s.PlusTop, Style: s})
				x++
			}
		}
		set(0, x, acell.Cell{Char: t.s.TopRight, Style: s})

		botY := totalH - 1
		set(botY, 0, acell.Cell{Char: t.s.BottomLeft, Style: s})
		x = 1
		for j := 0; j < cols; j++ {
			for k := 0; k < t.widths[j]+2; k++ {
				set(botY, x, acell.Cell{Char: t.s.Hor, Style: s})
				x++
			}
			if hasInnerV && j < cols-1 {
				set(botY, x, acell.Cell{Char: t.s.PlusBottom, Style: s})
				x++
			}
		}
		set(botY, x, acell.Cell{Char: t.s.BottomRight, Style: s})
	}

	if hasVer {
		for i := 0; i < rows; i++ {
			s := acell.Style{}
			set(rowY[i], 0, acell.Cell{Char: t.s.Ver, Style: s})
			set(rowY[i], totalW-1, acell.Cell{Char: t.s.Ver, Style: s})
		}
	}

	if hasInnerV {
		for i := 0; i < rows; i++ {
			s := acell.Style{}
			x := 1
			for j := 0; j < cols-1; j++ {
				x += t.widths[j] + 2
				set(rowY[i], x, acell.Cell{Char: t.s.Ver, Style: s})
				x++
			}
		}
	}

	if hasInnerH {
		for i := 0; i < rows-1; i++ {
			lineY := rowY[i] + 1
			s := acell.Style{}
			set(lineY, 0, acell.Cell{Char: t.s.PlusLeft, Style: s})
			x := 1
			for j := 0; j < cols; j++ {
				for k := 0; k < t.widths[j]+2; k++ {
					set(lineY, x, acell.Cell{Char: t.s.Hor, Style: s})
					x++
				}
				if hasInnerV && j < cols-1 {
					set(lineY, x, acell.Cell{Char: t.s.Plus, Style: s})
					x++
				}
			}
			set(lineY, x, acell.Cell{Char: t.s.PlusRight, Style: s})
		}
	}
}
