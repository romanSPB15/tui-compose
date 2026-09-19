package extra

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// Sparkline — виджет мини-графика (спарклайна).
// Отображает последовательность значений столбцами высотой до DataHeight
// строк, используя частичные блочные символы для сглаживания.
// Добавлено в TUI v3.3.0.
//
// s := extra.NewSparkline().WithValues([]int{1, 5, 3, 8}).WithHeight(3).AutoScale()
type Sparkline struct {
	Fill8, Fill7, Fill6, Fill5, Fill4, Fill3, Fill2, Fill1 rune

	values     []int
	BarStyle   func(i, v int) tui.Style
	div        float64
	DataHeight int
}

// NewSparkline создаёт спарклайн с Unicode-заполнителями по умолчанию,
// высотой в одну строку и единичным масштабом.
func NewSparkline() *Sparkline {
	return &Sparkline{
		Fill8:      '█',
		Fill7:      '▇',
		Fill6:      '▆',
		Fill5:      '▅',
		Fill4:      '▄',
		Fill3:      '▃',
		Fill2:      '▂',
		Fill1:      '▁',
		DataHeight: 1,
		div:        1,
	}
}

// WithValues устанавливает данные графика.
// Если масштаб ещё не был задан явно, пересчитывает его автоматически.
func (bc *Sparkline) WithValues(v []int) *Sparkline {
	bc.values = v
	if bc.div == 0 {
		bc.recalcDiv()
	}
	return bc
}

// AutoScale пересчитывает автоматический масштаб.
func (bc *Sparkline) AutoScale() *Sparkline {
	bc.recalcDiv()
	return bc
}

// WithScale устанавливает масштаб вручную.
func (bc *Sparkline) WithScale(div float64) *Sparkline {
	bc.div = div
	return bc
}

// Width возвращает ширину спарклайна — по количеству значений.
func (bc *Sparkline) Width() int {
	return len(bc.values)
}

// Height возвращает высоту спарклайна в строках.
func (bc *Sparkline) Height() int {
	return bc.DataHeight
}

// Render рисует спарклайн в буфер.
func (bc *Sparkline) Render(cells [][]cell.Cell) {
	h := bc.Height()

	for i, v := range bc.values {
		var s cell.Style
		if bc.BarStyle != nil {
			s = tui.ConvertToCellStyle(bc.BarStyle(i, v))
		}

		vDivided := int(float64(v) / bc.div)
		vDivided8 := vDivided / 8

		for z := vDivided8 - 1; z >= 0; z-- {
			if z > bc.DataHeight {
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
}

// Unicode переключает на блочные символы Unicode.
func (s *Sparkline) Unicode() *Sparkline {
	s.Fill8, s.Fill7, s.Fill6, s.Fill5, s.Fill4, s.Fill3, s.Fill2, s.Fill1 =
		'█', '▇', '▆', '▅', '▄', '▃', '▂', '▁'
	return s
}

// ASCII переключает на ASCII-заполнители.
func (s *Sparkline) ASCII() *Sparkline {
	s.Fill8 = '#'
	s.Fill7 = '#'
	s.Fill6 = '='
	s.Fill5 = '='
	s.Fill4 = '='
	s.Fill3 = '='
	s.Fill2 = '_'
	s.Fill1 = '_'

	return s
}

// WithHeight устанавливает высоту спарклайна (количество строк).
// Автоматически пересчитывает масштаб, если значения уже установлены.
func (s *Sparkline) WithHeight(h int) *Sparkline {
	s.DataHeight = h
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

// recalcDiv пересчитывает делитель масштаба так, чтобы максимум значений
// умещался в DataHeight строк с учётом дробной части блочных символов.
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
	s.div = float64(mx) / float64(s.DataHeight-s.DataHeight/8) / 8
	if s.div < 1 {
		s.div = 1
	}
}
