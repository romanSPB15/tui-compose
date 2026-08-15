package extra

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/cell"
)

type BarChart struct {
	FillRune  rune
	values    []int
	BarWidth  int
	Space     int
	BarStyle  func(i, v int) tui.Style
	TextStyle func(i, v int) tui.Style
	div       float32
	Height    int
}

func NewBarChart() *BarChart {
	return &BarChart{
		FillRune: '█',
		BarWidth: 3,
		Space:    1,
		Height:   20,
		div:      1,
	}
}

func (bc *BarChart) WithValues(v []int) *BarChart {
	bc.values = v
	mx := -1 << 31
	for _, v := range bc.values {
		if mx < v {
			mx = v
		}
	}
	bc.div = float32(mx) / float32(bc.Height) * 1.15
	if bc.div == 0 {
		bc.div = 1
	}
	return bc
}

func (bc *BarChart) MaxWidth() int {
	return (len(bc.values) * (bc.BarWidth + bc.Space)) - bc.Space
}

func (bc *BarChart) MaxHeight() int {
	return bc.Height
}

func (bc *BarChart) InnerText() string {
	w := bc.MaxWidth()
	h := bc.MaxHeight()
	cells := make([][]cell.Cell, h)
	for y := range cells {
		cells[y] = make([]cell.Cell, w)
		for x := range cells[y] {
			cells[y][x] = cell.Cell{Char: ' '}
		}
	}

	for i, v := range bc.values {
		var s cell.Style
		if bc.BarStyle != nil {
			s = convertToCellStyle(bc.BarStyle(i, v))
		}
		vDivided := int(float32(v) / bc.div)
		for z := vDivided; z >= 0; z-- {
			if z >= bc.Height {
				continue
			}
			for j := range bc.BarWidth {
				cells[h-z-1][i*(bc.BarWidth+bc.Space)+j] = cell.Cell{Char: bc.FillRune, Style: s}
			}
		}

		txt := strconv.Itoa(v)

		x := i*(bc.BarWidth+bc.Space) + (bc.BarWidth/2 - (len(txt) / 2))

		y := h - vDivided - 2
		if y < 0 {
			y = 0
		}

		if bc.TextStyle != nil {
			s = convertToCellStyle(bc.TextStyle(i, v))
		} else {
			s = cell.Style{}
		}

		for i, r := range []rune(txt) {
			cells[y][x+i] = cell.Cell{Char: r, Style: s}
		}
	}

	return cell.ToString(cells)
}

// WithBarWidth устанавливает ширину столбцов.
func (bc *BarChart) WithBarWidth(w int) *BarChart {
	bc.BarWidth = w
	return bc
}

// WithSpace устанавливает пропуски между столбцами.
func (bc *BarChart) WithSpace(s int) *BarChart {
	bc.Space = s
	return bc
}

// WithHeight устанавливает высоту графика в символах(примерно).
func (bc *BarChart) WithHeight(h int) *BarChart {
	bc.Height = h
	return bc
}

// WithBarStyle устанавливает стиль столбцов.
func (bc *BarChart) WithBarStyle(fn func(i, v int) tui.Style) *BarChart {
	bc.BarStyle = fn
	return bc
}

// WithTextStyle устанавливает стиль подписей.
func (bc *BarChart) WithTextStyle(fn func(i, v int) tui.Style) *BarChart {
	bc.TextStyle = fn
	return bc
}

// WithFillRune устанавливает символ-заполнитель столбцов.
func (bc *BarChart) WithFillRune(r rune) *BarChart {
	bc.FillRune = r
	return bc
}
