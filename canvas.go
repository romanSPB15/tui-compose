//go:build !no_widgets && !no_canvas

package tui

import "fmt"

// Canvas — это многострочный виджет, на котором можно "рисовать".
// В символах Canvas в 2 раза шире чем указано при создании, чтобы пиксели были квадратные а не прямоугольные.
type Canvas struct {
	width, height int
	pole          [][]Color
	idx           int
	wnd           Window
}

// NewCanvas() создаёт виждет Canvas.ы
func NewCanvas(width, height int) *Canvas {
	p := make([][]Color, height)
	for i := range height {
		p[i] = make([]Color, width)
	}
	return &Canvas{
		pole:   p,
		width:  width,
		height: height,
		idx:    -1,
	}
}

// Draw() устанавливает указанный цвет в указанном месте Canvas.
func (c *Canvas) Draw(x, y int, clr Color) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[x][y] = clr
}

// Draw() устанавливает указанный цвет в указанном месте Canvas, и перерисовывает.
func (c *Canvas) DrawAndRender(x, y int, clr Color) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[x][y] = clr
	if c.wnd != nil && c.idx != -1 {
		c.wnd.RedrawWidget(c.idx)
	}
}

// InnerText() реализует интерфейс Widget
func (c *Canvas) InnerText() (res string) {
	lastClr := Color(-1)
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			clr := c.pole[x][y]
			if lastClr != clr {
				if clr == 0 {
					res += "\033[0m"
				} else {
					res += fmt.Sprintf("\033[%dm", clr+10)
				}
				lastClr = clr
			}
			res += "  " // 2 пробела чтобы придать пикселям более квадратную форму
		}
		res += "\r\n"
	}
	res += "\033[0m"
	return
}

// MaxLength() реализует интерфейс Widget
func (c *Canvas) MaxLength() int {
	return (c.width*2 + 2) * c.height
}

// DisplayMode() реализует интерфейс Widget
func (*Canvas) DisplayMode() DisplayMode {
	return DisplayBlock
}

// SetIndex() реализует интерфейс Widget
func (c *Canvas) SetIndex(idx int) {
	c.idx = idx
}
