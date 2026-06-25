//go:build !no_widgets && !no_canvas

package tui

import (
	"fmt"
)

type PixelSize int

const (
	PixelTwoSymbol PixelSize = iota // "  "
	PixelOneSymbol                  // "▀"
)

// Canvas — это многострочный виджет, на котором можно "рисовать" цветные пиксели.
// В символах Canvas в 2 раза шире чем указано при создании, чтобы пиксели были квадратные а не прямоугольные.
type Canvas struct {
	width, height int
	pole          [][]Color
	idx           int
	PixelSize     PixelSize
}

// NewCanvas() создаёт виждет Canvas.
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
	if currentWindow != nil && c.idx != -1 {
		currentWindow.Redraw()
	}
}

// InnerText() реализует интерфейс Widget
func (c *Canvas) InnerText() (res string) {
	if c.PixelSize == PixelTwoSymbol {
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
	} else {
		lastBg, lastFr := Color(-1), Color(-1)
		z := Color(0)

		res += "\033[30m"
		for y := 0; y < c.height/2; y++ {
			for x := 0; x < c.width; x++ {
				var bg Color
				if y*2+1 == c.height {
					bg = z
				} else {
					bg = c.pole[x][y*2+1]
				}

				fr := c.pole[x][y*2]
				if lastBg != bg {
					if bg == z {
						res += "\033[40m"
					} else {
						res += fmt.Sprintf("\033[%dm", bg+10) //background
					}
					lastBg = bg
				}
				if lastFr != fr {
					if fr == z {
						res += "\033[30m"
					} else {
						res += fmt.Sprintf("\033[%dm", fr+10) //foreground
					}
					lastFr = fr
				}
				res += "▀"
			}
			res += "\r\n"
		}
	}

	res += "\033[0m"
	return
}

// MaxWidth() реализует интерфейс Widget
func (c *Canvas) MaxWidth() int {
	if c.PixelSize == PixelTwoSymbol {
		return (c.width * 2)
	}
	return c.width
}

// MaxHeight() реализует интерфейс Widget
func (c *Canvas) MaxHeight() int {
	if c.PixelSize == PixelTwoSymbol {
		return c.height
	}
	return c.height / 2
}

func (c *Canvas) Width() int {
	return c.width
}

func (c *Canvas) Height() int {
	return c.height
}

// Canvas — это многострочный виджет, на котором можно "рисовать" RGB-пиксели. Требуется терминал с True Color.
type CanvasRGB struct {
	width, height int
	pole          [][]ColorRGB
	idx           int
	PixelSize     PixelSize
}

// NewCanvas() создаёт виждет Canvas.ы
func NewCanvasRGB(width, height int) *CanvasRGB {
	p := make([][]ColorRGB, height)
	for i := range height {
		p[i] = make([]ColorRGB, width)
	}
	return &CanvasRGB{
		pole:      p,
		width:     width,
		height:    height,
		idx:       -1,
		PixelSize: PixelOneSymbol,
	}
}

// Draw() устанавливает указанный цвет в указанном месте Canvas.
func (c *CanvasRGB) Draw(x, y int, clr ColorRGB) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[x][y] = clr
}

// Draw() устанавливает указанный цвет в указанном месте Canvas, и перерисовывает.
func (c *CanvasRGB) DrawAndRender(x, y int, clr ColorRGB) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[x][y] = clr
	if currentWindow != nil && c.idx != -1 {
		currentWindow.Redraw()
	}
}

// InnerText() реализует интерфейс Widget
func (c *CanvasRGB) InnerText() (res string) {
	if c.PixelSize == PixelTwoSymbol {
		lastClr := ColorRGB{}
		z := ColorRGB{}
		for y := 0; y < c.height; y++ {
			for x := 0; x < c.width; x++ {
				clr := c.pole[x][y]
				if lastClr != clr {
					if clr == z {
						res += "\033[0m"
					} else {
						res += fmt.Sprintf("\033[48;2;%d;%d;%dm", clr.R, clr.G, clr.B)
					}
					lastClr = clr
				}
				res += "  " // 2 пробела чтобы придать пикселям более квадратную форму
			}
			res += "\r\n"
		}
	} else {
		lastBg := ColorRGB{}
		lastFr := ColorRGB{}
		z := ColorRGB{}

		res += "\033[30m"
		for y := 0; y < c.height/2; y++ {
			for x := 0; x < c.width; x++ {
				var bg ColorRGB
				if y*2+1 == c.height {
					bg = ColorRGB{}
				} else {
					bg = c.pole[x][y*2+1]
				}

				fr := c.pole[x][y*2]
				if lastBg != bg {
					if bg == z {
						res += "\033[40m"
					} else {
						res += fmt.Sprintf("\033[48;2;%d;%d;%dm", bg.R, bg.G, bg.B)
					}
					lastBg = bg
				}
				if lastFr != fr {
					if fr == z {
						res += "\033[30m"
					} else {
						res += fmt.Sprintf("\033[38;2;%d;%d;%dm", fr.R, fr.G, fr.B)
					}
					lastFr = fr
				}
				res += "▀"
			}
			res += "\r\n"
		}
	}

	res += "\033[0m"
	return
}

// MaxWidth() реализует интерфейс Widget
func (c *CanvasRGB) MaxWidth() int {
	if c.PixelSize == PixelTwoSymbol {
		return (c.width * 2)
	}
	return c.width
}

// MaxHeight() реализует интерфейс Widget
func (c *CanvasRGB) MaxHeight() int {
	if c.PixelSize == PixelTwoSymbol {
		return c.height
	}
	return c.height / 2
}

func (c *CanvasRGB) Width() int {
	return c.width
}

func (c *CanvasRGB) Height() int {
	return c.height
}

func init() {
	var _ Widget = (*Canvas)(nil)
	var _ Widget = (*CanvasRGB)(nil)
}
