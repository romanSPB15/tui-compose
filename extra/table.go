package extra

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

type TableCellAlign int

const (
	AlignLeft   = iota
	AlignCenter = iota
	AlignRight  = iota
)

type HorSeparator int

const (
	NoHorSeparator = iota
	BetweenHorSeparator
	EverywhereHorSeparator
)

type TableCell struct {
	Style tui.Style
	Text  string
	Align TableCellAlign
}

type TableStyle struct {
	Hor, Ver                                       rune // — │
	TopRight, TopLeft, BottomRight, BottomLeft     rune // ┐ ┌ ┘ └
	Plus, PlusLeft, PlusRight, PlusBottom, PlusTop rune // ┼ ├ ┤ ┴ ┬
}

type Table struct {
	data   [][]TableCell
	widths []int
	hor    HorSeparator
	s      TableStyle
}

func NewTable(cells [][]TableCell) *Table {
	if len(cells) == 0 {
		return (&Table{
			data: cells,
			hor:  BetweenHorSeparator,
		}).Default()
	}
	tbl := &Table{
		data:   cells,
		widths: make([]int, len(cells[0])),
		hor:    BetweenHorSeparator,
	}

	for _, v2 := range cells { // по строкам
		for x, v := range v2 { // по ячейкам
			rlen := len([]rune(v.Text))
			if tbl.widths[x] < rlen {
				tbl.widths[x] = rlen
			}
		}
	}

	tbl.Default()

	return tbl
}

func (tbl *Table) Default() *Table {
	tbl.s = TableStyle{
		Hor: '─', Ver: '│',
		TopLeft: '┌', TopRight: '┐',
		BottomLeft: '└', BottomRight: '┘',
		Plus: '┼', PlusLeft: '├', PlusRight: '┤',
		PlusTop: '┬', PlusBottom: '┴',
	}
	return tbl
}

func (tbl *Table) Rounded() *Table {
	tbl.s = TableStyle{
		Hor: '─', Ver: '│',
		TopLeft: '╭', TopRight: '╮',
		BottomLeft: '╰', BottomRight: '╯',
		Plus: '┼', PlusLeft: '├', PlusRight: '┤',
		PlusTop: '┬', PlusBottom: '┴',
	}
	return tbl
}

func (tbl *Table) ASCII() *Table {
	tbl.s = TableStyle{
		Hor: '-', Ver: '|',
		TopLeft: '+', TopRight: '+',
		BottomLeft: '+', BottomRight: '+',
		Plus: '+', PlusLeft: '+', PlusRight: '+',
		PlusTop: '+', PlusBottom: '+',
	}
	return tbl
}

func (pc *Table) Width() int {
	res := 1
	for _, v := range pc.widths {
		res += v + 3
	}
	return res
}

func (pc *Table) Height() int {
	if pc.hor == NoHorSeparator {
		return len(pc.data)
	}
	if pc.hor == BetweenHorSeparator {
		return len(pc.data)*2 - 1
	}
	return len(pc.data)*2 + 1
}

func (pc *Table) WithData(tc [][]TableCell) *Table {
	pc.data = tc
	if len(pc.data) != 0 {
		if len(pc.widths) < len(pc.data[0]) {
			for range len(pc.data[0]) - len(pc.widths) {
				pc.widths = append(pc.widths, 0)
			}
		}

		for _, v2 := range pc.data {
			for x, v := range v2 {
				rlen := len([]rune(v.Text))
				if pc.widths[x] < rlen {
					pc.widths[x] = rlen
				}
			}
		}
	}

	if tui.CurrentWindow() != nil {
		tui.CurrentWindow().Do(tui.CurrentWindow().Redraw)
	}
	return pc
}

func (t *Table) Render(buf [][]cell.Cell) {
	if len(t.data) == 0 {
		return
	}
	h := t.Height()

	style := func(s tui.Style) cell.Style {
		return tui.ConvertToCellStyle(s)
	}

	if t.hor == EverywhereHorSeparator {
		buf[0][0] = cell.Cell{Char: t.s.TopLeft, Style: cell.Style{}}
		x := 1
		for col, width := range t.widths {
			for i := 0; i < width+2; i++ {
				buf[0][x] = cell.Cell{Char: t.s.Hor, Style: cell.Style{}}
				x++
			}
			if col < len(t.widths)-1 {
				buf[0][x] = cell.Cell{Char: t.s.PlusTop, Style: cell.Style{}}
				x++
			}
		}
		buf[0][x] = cell.Cell{Char: t.s.TopRight, Style: cell.Style{}}
	}

	rowY := 0
	if t.hor == EverywhereHorSeparator {
		rowY = 1
	}

	for rowIdx, row := range t.data {
		x := 0

		buf[rowY][x] = cell.Cell{Char: t.s.Ver, Style: cell.Style{}}
		x++

		for colIdx, cellData := range row {
			colWidth := t.widths[colIdx]
			text := cellData.Text
			textRunes := []rune(text)
			textLen := len(textRunes)

			leftPad, rightPad := 0, 0
			switch cellData.Align {
			case AlignLeft:
				rightPad = colWidth - textLen
			case AlignRight:
				leftPad = colWidth - textLen
			case AlignCenter:
				leftPad = (colWidth - textLen) / 2
				rightPad = colWidth - textLen - leftPad
			}

			buf[rowY][x] = cell.Cell{Char: ' ', Style: style(cellData.Style)}
			x++

			for i := 0; i < leftPad; i++ {
				buf[rowY][x] = cell.Cell{Char: ' ', Style: style(cellData.Style)}
				x++
			}

			for _, r := range textRunes {
				buf[rowY][x] = cell.Cell{Char: r, Style: style(cellData.Style)}
				x++
			}

			for i := 0; i < rightPad; i++ {
				buf[rowY][x] = cell.Cell{Char: ' ', Style: style(cellData.Style)}
				x++
			}

			buf[rowY][x] = cell.Cell{Char: ' ', Style: style(cellData.Style)}
			x++

			if colIdx < len(row)-1 {
				buf[rowY][x] = cell.Cell{Char: t.s.Ver, Style: cell.Style{}}
				x++
			} else {
				buf[rowY][x] = cell.Cell{Char: t.s.Ver, Style: cell.Style{}}
				x++
			}
		}

		rowY++

		if rowIdx < len(t.data)-1 && t.hor >= BetweenHorSeparator {
			buf[rowY][0] = cell.Cell{Char: t.s.PlusLeft, Style: cell.Style{}}
			x := 1
			for col, width := range t.widths {
				for i := 0; i < width+2; i++ {
					buf[rowY][x] = cell.Cell{Char: t.s.Hor, Style: cell.Style{}}
					x++
				}
				if col < len(t.widths)-1 {
					buf[rowY][x] = cell.Cell{Char: t.s.Plus, Style: cell.Style{}}
					x++
				}
			}
			buf[rowY][x] = cell.Cell{Char: t.s.PlusRight, Style: cell.Style{}}
			rowY++
		}
	}

	if t.hor == EverywhereHorSeparator {
		buf[h-1][0] = cell.Cell{Char: t.s.BottomLeft, Style: cell.Style{}}
		x := 1
		for col, width := range t.widths {
			for i := 0; i < width+2; i++ {
				buf[h-1][x] = cell.Cell{Char: t.s.Hor, Style: cell.Style{}}
				x++
			}
			if col < len(t.widths)-1 {
				buf[h-1][x] = cell.Cell{Char: t.s.PlusBottom, Style: cell.Style{}}
				x++
			}
		}
		buf[h-1][x] = cell.Cell{Char: t.s.BottomRight, Style: cell.Style{}}
	}
}

func (t *Table) WithStyle(s TableStyle) *Table {
	t.s = s
	return t
}

func (t *Table) WithHorSeparator(h HorSeparator) *Table {
	t.hor = h
	return t
}
