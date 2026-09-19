package extra

import (
	"fmt"
	"math"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// PieData — сектор круговой диаграммы.
// Label — подпись для легенды.
// Value — значение (в любых единицах, важно только соотношение).
// Color — цвет сектора.
// Добавлено в TUI v3.3.0.
type PieData struct {
	Label string
	Value float64
	Color tui.Style
}

// PieChart — виджет круговой диаграммы.
// Добавлено в TUI v3.3.0.
type PieChart struct {
	data        []PieData
	radius      int
	showPercent bool
	showLegend  bool
	valueStyle  tui.Style
	legendStyle tui.Style
}

// NewPieChart создаёт круговую диаграмму с указанными данными.
// Радиус по умолчанию 10, проценты отображаются, легенда скрыта.
// Добавлено в TUI v3.3.0.
func NewPieChart(data []PieData) *PieChart {
	return &PieChart{
		data:        data,
		radius:      10,
		showPercent: true,
		showLegend:  false,
		valueStyle:  tui.FrDefault,
		legendStyle: tui.FrDefault,
	}
}

// WithRadius устанавливает радиус диаграммы.
// Добавлено в TUI v3.3.0.
func (pc *PieChart) WithRadius(r int) *PieChart {
	pc.radius = r
	return pc
}

// WithShowPercent включает или отключает отображение процентов.
// Добавлено в TUI v3.3.0.
func (pc *PieChart) WithShowPercent(b bool) *PieChart {
	pc.showPercent = b
	return pc
}

// WithShowLegend включает или отключает отображение легенды.
// Добавлено в TUI v3.3.0.
func (pc *PieChart) WithShowLegend(b bool) *PieChart {
	pc.showLegend = b
	return pc
}

// WithValueStyle устанавливает стиль текста процентов.
// Добавлено в TUI v3.3.0.
func (pc *PieChart) WithValueStyle(s tui.Style) *PieChart {
	pc.valueStyle = s
	return pc
}

// WithLegendStyle устанавливает стиль текста легенды.
// Добавлено в TUI v3.3.0.
func (pc *PieChart) WithLegendStyle(s tui.Style) *PieChart {
	pc.legendStyle = s
	return pc
}

// Width реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (pc *PieChart) Width() int {
	width := 2*pc.radius + 2
	if pc.showLegend {
		maxLabelLen := 0
		for _, d := range pc.data {
			if len(d.Label) > maxLabelLen {
				maxLabelLen = len(d.Label)
			}
		}
		legendWidth := maxLabelLen + 10
		if legendWidth > 0 {
			width += legendWidth
		}
	}
	return width
}

// Height реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (pc *PieChart) Height() int {
	return pc.radius
}

// Render реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (pc *PieChart) Render(cells [][]cell.Cell) {
	if len(pc.data) == 0 {
		return
	}

	// центр окружности

	cx := pc.radius
	cy := pc.radius / 2

	// PieData.Value может быть в любых еденицах, пересчитываем в угловой размер

	total := 0.0

	for _, d := range pc.data {
		total += d.Value
	}
	if total == 0 {
		return
	}

	angles := []float64{} // углы конца секторов
	angle := 0.0

	for _, d := range pc.data {
		angle += (d.Value / total) * (2 * math.Pi)
		angles = append(angles, angle)
	}

	w := pc.Width()
	h := pc.Height()

	// Получаем цвет пикселя в точке x, y

	getColor := func(x, y int) cell.Style {
		// Считаем угол от центра к этой точке

		dx := float64(x - cx)
		dy := float64(y-cy) * 2.2
		dist2 := dx*dx + dy*dy

		if dist2 > float64(pc.radius*pc.radius) {
			return cell.Style{}
		}

		angle := math.Atan2(dy, dx) // Atan2 вовзращает угол в радианах направления из 0, 0 к точке x, y
		// ! y сначала !

		if angle < 0 { // нормализуем угол
			angle += 2 * math.Pi
		}

		// ищем сектор

		for i, v := range angles {
			if angle <= v { // если угол попадает в сектор
				return tui.ConvertToCellStyle(pc.data[i].Color)
			}
		}
		return cell.Style{}
	}

	// Рисуем на матрице

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s := getColor(x, y)
			if s != (cell.Style{}) {
				cells[y][x] = cell.Cell{Char: '█', Style: s}
			}
		}
	}

	if pc.showPercent {
		drawText := func(cx, cy int, str string) {
			l := len(str)
			sx := cx - l/2 // начало текста
			for i, v := range []rune(str) {
				cells[cy][sx+i] = cell.Cell{Char: v, Style: tui.ConvertToCellStyle(pc.valueStyle)}
			}
		}

		before := 0.0

		for i, a := range angles {
			textRadius := pc.radius / 3 * 2

			am := (before+a)/2 - 0.1 // середина

			x := cx + int(math.Cos(am)*float64(textRadius))
			y := cy + int(math.Sin(am)*float64(textRadius)/2.2) // коррекция соотношения сторон

			drawText(x, y, fmt.Sprintf("%.0f%%", (pc.data[i].Value/total)*100))

			before = a
		}
	}

	if pc.showLegend {
		// рисуем легенду
		for j, v := range pc.data {
			// формат: ■ Label (25%)

			str := fmt.Sprintf("■ %s (%.0f%%)", v.Label, (pc.data[j].Value / total * 100))
			legendStartX := pc.radius*2 + 4

			for i, v := range []rune(str) {
				if i == 0 {
					cells[j][i+legendStartX] = cell.Cell{Char: v, Style: tui.ConvertToCellStyle(pc.data[j].Color)}
				} else {
					cells[j][i+legendStartX] = cell.Cell{Char: v, Style: tui.ConvertToCellStyle(pc.legendStyle)}
				}

			}
		}
	}
}
