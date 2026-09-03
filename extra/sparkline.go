package extra

import (
	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/cell"
)

type Sparkline struct {
	Fill8, Fill7, Fill6, Fill5, Fill4, Fill3, Fill2, Fill1 rune

	values   []int
	BarStyle func(i, v int) tui.Style
	div      float64
	Height   int
}

func NewSparkline() *Sparkline {
	return &Sparkline{
		Fill8:  '█',
		Fill7:  '▇',
		Fill6:  '▆',
		Fill5:  '▅',
		Fill4:  '▄',
		Fill3:  '▃',
		Fill2:  '▂',
		Fill1:  '▁',
		Height: 1,
		div:    1,
	}
}
func (bc *Sparkline) WithValues(v []int) *Sparkline {
	bc.values = v
	if bc.div == 0 {
		bc.recalcDiv()
	}
	return bc
}

func (bc *Sparkline) AutoScale() *Sparkline {
	bc.recalcDiv()
	return bc
}

func (bc *Sparkline) WithScale(div float64) *Sparkline {
	bc.div = div
	return bc
}

func (bc *Sparkline) Width() int {
	return len(bc.values)
}

func (bc *Sparkline) Height() int {
	return bc.Height
}

func (bc *Sparkline) InnerText() string {
	w := bc.Width()
	h := bc.Height()
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
			s = tui.ConvertToCellStyle(bc.BarStyle(i, v))
		}

		vDivided := int(float64(v) / bc.div)
		vDivided8 := vDivided / 8

		for z := vDivided8 - 1; z >= 0; z-- {
			if z > bc.Height {
				continue
			}

			cells[h-z-1][i] = cell.Cell{Char: bc.Fill8, Style: s}
		}
		yLast := h - (vDivided8) - 1
		switch vDivided % 8 {
		case 1:
			cells[yLast][i] = cell.Cell{Char: bc.Fill1, Style: s}
		case 2:
			cells[yLast][i] = cell.Cell{Char: bc.Fill2, Style: s}
		case 3:
			cells[yLast][i] = cell.Cell{Char: bc.Fill3, Style: s}
		case 4:
			cells[yLast][i] = cell.Cell{Char: bc.Fill4, Style: s}
		case 5:
			cells[yLast][i] = cell.Cell{Char: bc.Fill5, Style: s}
		case 6:
			cells[yLast][i] = cell.Cell{Char: bc.Fill6, Style: s}
		case 7:
			cells[yLast][i] = cell.Cell{Char: bc.Fill7, Style: s}
		}
	}

	return cell.ToString(cells)
}

func (s *Sparkline) Unicode() *Sparkline {
	s.Fill8, s.Fill7, s.Fill6, s.Fill5, s.Fill4, s.Fill3, s.Fill2, s.Fill1 =
		'█', '▇', '▆', '▅', '▄', '▃', '▂', '▁'
	return s
}

func (s *Sparkline) ASCII() *Sparkline {
	s.Fill8 = '#'
	s.Fill7 = '#'
	s.Fill6 = '='
	s.Fill5 = '='
	s.Fill4 = '='
	s.Fill3 = '='
	s.Fill2 = '_'
	s.Fill1 = '_'

	s.Height = 3
	return s
}

// WithHeight устанавливает высоту спарклайна (количество строк).
// Автоматически пересчитывает масштаб, если значения уже установлены.
func (s *Sparkline) WithHeight(h int) *Sparkline {
	s.Height = h
	if len(s.values) > 0 {
		s.recalcDiv()
	}
	return s
}

// WithBarStyle устанавливает функцию стилизации столбцов.
func (s *Sparkline) WithBarStyle(fn func(i, v int) tui.Style) *Sparkline {
	s.BarStyle = fn
	return s
}

func (s *Sparkline) recalcDiv() {
	if len(s.values) == 0 {
		s.div = 1
		return
	}
	mx := 0
	for _, v := range s.values {
		if v > mx {
			mx = v
		}
	}
	if mx == 0 {
		mx = 1
	}
	s.div = float64(mx) / float64(s.Height-s.Height/8) / 8
	if s.div < 1 {
		s.div = 1
	}
}
