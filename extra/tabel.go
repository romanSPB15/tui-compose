package extra

import (
	"strings"

	"github.com/romanSPB15/tui-compose/v3"
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

func (pc *Table) MaxWidth() int {
	res := 1
	for _, v := range pc.widths {
		res += v + 3
	}
	return res
}

func (pc *Table) MaxHeight() int {
	if pc.hor == NoHorSeparator {
		return len(pc.data)
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

func (pc *Table) InnerText() string {
	if len(pc.data) == 0 {
		return ""
	}

	var b strings.Builder

	// разделитель сверху
	if pc.hor == EverywhereHorSeparator {

		b.WriteRune(pc.s.TopLeft)

		for i, w := range pc.widths {
			for range w + 2 {
				b.WriteRune(pc.s.Hor)
			}
			if i < len(pc.widths)-1 {
				b.WriteRune(pc.s.PlusTop)
			}
		}

		b.WriteRune(pc.s.TopRight)
		b.WriteString("\n")
	}

	for i, row := range pc.data {
		b.WriteRune(pc.s.Ver)
		b.WriteRune(' ')

		for i, c := range row {
			b.WriteString(c.Style.String())

			l := len([]rune(c.Text))
			w := pc.widths[i]

			if c.Align == AlignRight {
				if l < w {
					b.WriteString(strings.Repeat(" ", w-l))
				}
			}

			if c.Align == AlignCenter {
				if l < w {
					b.WriteString(strings.Repeat(" ", (w-l)/2))
				}
			}

			b.WriteString(c.Text)

			if c.Align == AlignLeft {
				if l < w {
					b.WriteString(strings.Repeat(" ", w-l))
				}
			}

			if c.Align == AlignCenter {
				if l < w {
					b.WriteString(strings.Repeat(" ", (w-l)/2))
					if (w-l)%2 == 1 {
						b.WriteRune(' ')
					}
				}
			}

			if c.Style != 0 {
				b.WriteString("\033[0m")
			}

			b.WriteRune(' ')
			b.WriteRune(pc.s.Ver)
			b.WriteRune(' ')
		}
		if i < len(pc.data)-1 {
			b.WriteString("\n")
			if pc.hor >= BetweenHorSeparator {
				b.WriteRune(pc.s.PlusLeft)

				for i, w := range pc.widths {
					for range w + 2 {
						b.WriteRune(pc.s.Hor)
					}
					if i < len(pc.widths)-1 {
						b.WriteRune(pc.s.Plus)
					}
				}

				b.WriteRune(pc.s.PlusRight)
				b.WriteString("\n")
			}
		}
	}

	// разделитель снизу
	if pc.hor == EverywhereHorSeparator {
		b.WriteString("\n")
		b.WriteRune(pc.s.BottomLeft)

		for i, w := range pc.widths {
			for range w + 2 {
				b.WriteRune(pc.s.Hor)
			}
			if i < len(pc.widths)-1 {
				b.WriteRune(pc.s.PlusBottom)
			}
		}

		b.WriteRune(pc.s.BottomRight)
	}

	return b.String()
}

func (t *Table) WithStyle(s TableStyle) *Table {
	t.s = s
	return t
}

func (t *Table) WithHorSeparator(h HorSeparator) *Table {
	t.hor = h
	return t
}
