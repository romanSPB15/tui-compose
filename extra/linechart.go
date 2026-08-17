package extra

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/cell"
)

type AxisRunes struct {
	Hor, Ver rune // линии, наример ─ и │
	Corner   rune // нижний левый угол, например └
}

type Series struct {
	Values     []int
	LineStyle  tui.Style                // если 0 → используется глобальный
	PointStyle func(i, v int) tui.Style // если nil → используется глобальный
}

type LineChart struct {
	PointRune     rune // символ точек
	LineRune      rune // символ линий
	Data          []Series
	PointDistance int // расстояние между точками
	div           float64
	Height        int // высота поля графика

	DisplayPoints bool // отображение точек значений

	YLabels        []int     // вертикальные подписи
	XLabels        []string  // горизонтальные подписи
	AxisStyle      tui.Style // стиль осей
	AxisLabelStyle tui.Style // стиль подписей
	AxisRunes      AxisRunes // символы осей
}

func NewLineChart() *LineChart {
	return &LineChart{
		PointRune:     '●',
		LineRune:      '·',
		Height:        20,
		div:           1,
		PointDistance: 5,
		AxisRunes: AxisRunes{
			Hor:    '─',
			Ver:    '│',
			Corner: '└',
		},
	}
}

func (bc *LineChart) MaxWidth() int {
	mx := 0
	for _, v := range bc.Data {
		dataWidth := (len(v.Values)-1)*bc.PointDistance + 1
		if mx < dataWidth {
			mx = dataWidth
		}
	}
	if len(bc.YLabels) > 0 {
		return mx + 4
	}
	return mx
}

func (bc *LineChart) MaxHeight() int {
	if len(bc.YLabels) > 0 {
		return bc.Height + 1
	}
	return bc.Height
}

// алгоритм Брезенхема
func (bc *LineChart) drawLine(x1, y1, x2, y2 int, cells [][]cell.Cell, c cell.Cell) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	x, y := x1, y1
	for {
		if x >= 0 && x < len(cells[0]) && y >= 0 && y < len(cells) {
			cells[y][x] = c
		}
		if x == x2 && y == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func lenNumber(n int) (i int) {
	i = 1
	for n > 9 {
		i++
		n /= 10
	}
	return
}

func (bc *LineChart) InnerText() string {
	if len(bc.Data) == 0 {
		return ""
	}

	maxVal := 0
	for _, v2 := range bc.Data {
		for _, v := range v2.Values {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	offsetX := 0
	if len(bc.YLabels) > 0 {
		maxLabelLen := 0
		for _, v := range bc.YLabels {
			l := lenNumber(v)
			if l > maxLabelLen {
				maxLabelLen = l
			}
		}
		offsetX = maxLabelLen + 1
	}

	w := bc.MaxWidth()
	h := bc.MaxHeight()

	cells := make([][]cell.Cell, h)
	for y := range cells {
		cells[y] = make([]cell.Cell, w)
		for x := range cells[y] {
			cells[y][x] = cell.Cell{Char: ' '}
		}
	}

	axisStyle := convertToCellStyle(bc.AxisStyle)
	var axisLabelStyle cell.Style
	if bc.AxisLabelStyle != 0 {
		axisLabelStyle = tui.ConvertToCellStyle(bc.AxisLabelStyle)
	} else {
		axisLabelStyle = axisStyle
	}

	if bc.div == 0 {
		bc.AutoScale()
	}

	if len(bc.YLabels) > 0 {
		for y := 0; y < h; y++ {
			cells[y][offsetX-1] = cell.Cell{Char: bc.AxisRunes.Ver, Style: axisStyle}
		}
		for _, v := range bc.YLabels {
			y := h - int(float64(v)/bc.div) - 1

			txt := reverse(strconv.Itoa(v))

			if y < 0 {
				y = 0
			}

			for j, r := range []rune(txt) {
				x := offsetX - 2 - j
				if x >= 0 {
					cells[y][x] = cell.Cell{Char: r, Style: axisLabelStyle}
				}
			}
		}
	}

	if len(bc.XLabels) > 0 {
		for x := offsetX; x < w; x++ {
			cells[h-2][x] = cell.Cell{Char: bc.AxisRunes.Hor, Style: axisStyle}
		}
		for i, v := range bc.XLabels {
			x := offsetX + 1 + i*bc.PointDistance - len(v)/2
			if x < 0 {
				x = 0
			}
			end := x + len(v)
			if end > w-1 {
				x -= end - w + 1
			}
			for j, r := range []rune(v) {
				cells[h-1][x+j] = cell.Cell{Char: r, Style: axisLabelStyle}
			}
		}
	}

	if len(bc.XLabels) > 0 && len(bc.YLabels) > 0 {
		cells[h-2][offsetX-1] = cell.Cell{Char: bc.AxisRunes.Corner, Style: axisStyle}
		cells[h-1][offsetX-1] = cell.Cell{Char: ' ', Style: axisStyle}
	}

	for _, s := range bc.Data {
		lineStyle := convertToCellStyle(s.LineStyle)
		for i := 0; i < len(s.Values)-1; i++ {
			x1 := offsetX + i*bc.PointDistance
			y1 := h - 1 - int(float64(s.Values[i])/bc.div)
			x2 := offsetX + (i+1)*bc.PointDistance
			y2 := h - 1 - int(float64(s.Values[i+1])/bc.div)
			if len(bc.XLabels) > 0 {
				y1 -= 2
				y2 -= 2
			}
			if y1 < 0 {
				y1 = 0
			}
			if y2 < 0 {
				y2 = 0
			}
			bc.drawLine(x1, y1, x2, y2, cells, cell.Cell{Char: bc.LineRune, Style: lineStyle})
		}

		if bc.DisplayPoints {
			for i, v := range s.Values {
				x := offsetX + i*bc.PointDistance
				y := h - 3 - int(float64(v)/bc.div)
				if y < 0 {
					y = 0
				}
				if y < h && x < w {
					var st cell.Style
					if s.PointStyle != nil {
						st = tui.ConvertToCellStyle(s.PointStyle(i, v))
					}
					cells[y][x] = cell.Cell{Char: bc.PointRune, Style: st}
				}
			}
		}
	}

	return cell.ToString(cells)
}

// Default — стандартные Unicode-символы (─ │ └ ·)
func (lc *LineChart) WithDefaultAxis() *LineChart {
	lc.AxisRunes = AxisRunes{
		Hor:    '─',
		Ver:    '│',
		Corner: '└',
	}
	lc.PointRune = '●'
	lc.LineRune = '·'
	return lc
}

// ASCII — ASCII-совместимые символы
func (lc *LineChart) WithASCIIAxis() *LineChart {
	lc.AxisRunes = AxisRunes{
		Hor:    '-',
		Ver:    '|',
		Corner: '+',
	}
	lc.LineRune = '+'
	lc.PointRune = '*'
	return lc
}

// Rounded — скруглённый угол
func (lc *LineChart) WithRoundedAxis() *LineChart {
	lc.AxisRunes = AxisRunes{
		Hor:    '─',
		Ver:    '│',
		Corner: '╰',
	}
	lc.PointRune = '●'
	lc.LineRune = '·'
	return lc
}

// WithXLabels устанавливает подписи горизонтальной оси.
func (lc *LineChart) WithXLabels(lbl []string) *LineChart {
	lc.XLabels = lbl
	return lc
}

// WithYLabels устанавливает подписи вертикальной оси вручную.
func (lc *LineChart) WithYLabels(lbl []int) *LineChart {
	lc.YLabels = lbl
	return lc
}

// WithScale устанавливает масштаб вручную.
func (lc *LineChart) WithScale(s float64) *LineChart {
	lc.div = s
	return lc
}

// AutoScale устанавливает автоматический масштаб.
func (lc *LineChart) AutoScale() *LineChart {
	maxVal := 0
	for _, v2 := range lc.Data {
		for _, v := range v2.Values {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}
	lc.div = float64(maxVal) / (float64(lc.Height) * 0.8)
	return lc
}

// WithData устанавливает данные.
func (lc *LineChart) WithData(d []Series) *LineChart {
	lc.Data = d
	return lc
}

func (lc *LineChart) GenerateYLabels(count int, nice bool) *LineChart {
	if len(lc.YLabels) > 0 {
		return lc
	}
	if count < 2 {
		count = 2 // хотя бы две подписи
	}

	maxVal := 0
	for _, v2 := range lc.Data {
		for _, v := range v2.Values {
			if v > maxVal {
				maxVal = v
			}
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}
	step := maxVal / (count)

	if step == 0 {
		step = 1
	}
	for y := step; y <= maxVal; y += step {
		if y != 0 {
			if nice {
				lc.YLabels = append(lc.YLabels, roundToNice(y))
			} else {
				lc.YLabels = append(lc.YLabels, y)
			}

		}
	}

	return lc
}

func roundToNice(x int) int {
	goodList := []int{5, 25, 75, 125, 175}
	for _, v := range goodList {
		if abs(x-v) < v/10 {
			return v
		}
	}
	if x <= 0 {
		return 1
	}
	if x < 10 {
		if x == 5 {
			return x
		}
		return (x / 2) * 2
	}
	if x < 100 {
		return (x / 10) * 10
	}
	if x < 1000 {
		return (x / 50) * 50
	}
	if x < 10000 {
		return (x / 500) * 500
	}
	return (x / 1000) * 1000
}

// WithHeight устанавливает высоту графика (количество строк), без учёта горизонтальной оси.
func (lc *LineChart) WithHeight(h int) *LineChart {
	lc.Height = h
	return lc
}

// WithPointDistance устанавливает расстояние между точками по горизонтали.
func (lc *LineChart) WithPointDistance(d int) *LineChart {
	lc.PointDistance = d
	return lc
}

// WithDisplayPoints включает/отключает отображение точек значений.
func (lc *LineChart) WithDisplayPoints(b bool) *LineChart {
	lc.DisplayPoints = b
	return lc
}
