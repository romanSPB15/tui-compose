//go:build !no_widgets && !no_canvas

package tui

import (
	"image"
	"strconv"
	"strings"
)

type PixelSize int

const (
	PixelTwoSymbol PixelSize = iota // "  "
	PixelOneSymbol                  // "▀"
)

// Canvas — это многострочный виджет, на котором можно "рисовать" цветные пиксели.
type Canvas struct {
	width, height int
	pole          [][]Color
	idx           int
	PixelSize     PixelSize
	dirty         bool
	cached        string
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
	}
}

// Draw() устанавливает указанный цвет в указанном месте Canvas.
func (c *Canvas) Draw(x, y int, clr Color) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[y][x] = clr
	c.dirty = true
}

// Draw() устанавливает указанный цвет в указанном месте Canvas, и перерисовывает.
func (c *Canvas) DrawAndRender(x, y int, clr Color) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[y][x] = clr
	c.dirty = true
	if currentWindow != nil {
		currentWindow.Redraw()
	}
}

// InnerText() реализует интерфейс Widget
func (c *Canvas) InnerText() string {
	if !c.dirty && c.cached != "" {
		return c.cached
	}

	var b strings.Builder
	b.Grow(c.width * c.height * 6)

	if c.PixelSize == PixelTwoSymbol {
		lastClr := Color(-1)
		for y := 0; y < c.height; y++ {
			for x := 0; x < c.width; x++ {
				clr := c.pole[y][x]
				if lastClr != clr {
					if clr == NoColor {
						b.WriteString("\033[0m")
					} else {
						// Фоновый цвет: clr + 10
						b.WriteString("\033[")
						b.WriteString(strconv.Itoa(int(clr + 10)))
						b.WriteString("m")
					}
					lastClr = clr
				}
				b.WriteString("  ")
			}
			b.WriteString("\r\n")
		}
	} else {
		lastBg := Color(-1)
		lastFr := Color(-1)
		z := Color(0)

		b.WriteString("\033[30m")
		for y := 0; y < c.height/2; y++ {
			for x := 0; x < c.width; x++ {
				var bg Color
				if y*2+1 == c.height {
					bg = z
				} else {
					bg = c.pole[y*2+1][x]
				}

				fr := c.pole[y*2][x]
				if lastBg != bg {
					if bg == z {
						b.WriteString("\033[40m")
					} else {
						b.WriteString("\033[")
						b.WriteString(strconv.Itoa(int(bg + 10)))
						b.WriteString("m")
					}
					lastBg = bg
				}
				if lastFr != fr {
					if fr == z {
						b.WriteString("\033[30m")
					} else {
						b.WriteString("\033[")
						b.WriteString(strconv.Itoa(int(fr)))
						b.WriteString("m")
					}
					lastFr = fr
				}
				b.WriteString("▀")
			}
			b.WriteString("\r\n")
		}
	}

	b.WriteString("\033[0m")

	c.cached = b.String()
	c.dirty = false
	return c.cached
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

func (cnv *Canvas) Get(x, y int) Color {
	if x < 0 || x >= cnv.width || y < 0 || y >= cnv.height {
		return NoColor
	}
	return cnv.pole[y][x]
}

// Canvas — это многострочный виджет, на котором можно "рисовать" RGB-пиксели. Требуется терминал с True Color.
type CanvasRGB struct {
	width, height int
	pole          [][]ColorRGB
	PixelSize     PixelSize
	dirty         bool
	cached        string
}

// NewCanvas() создаёт виждет Canvas.
func NewCanvasRGB(width, height int) *CanvasRGB {
	p := make([][]ColorRGB, height)
	for i := range height {
		p[i] = make([]ColorRGB, width)
	}
	return &CanvasRGB{
		pole:      p,
		width:     width,
		height:    height,
		PixelSize: PixelOneSymbol,
	}
}

// Draw() устанавливает указанный цвет в указанном месте Canvas.
func (c *CanvasRGB) Draw(x, y int, clr ColorRGB) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[y][x] = clr
	c.dirty = true
}

// Draw() устанавливает указанный цвет в указанном месте Canvas, и перерисовывает.
func (c *CanvasRGB) DrawAndRender(x, y int, clr ColorRGB) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	c.pole[y][x] = clr
	c.dirty = true
	if currentWindow != nil {
		currentWindow.Redraw()
	}
}

// InnerText() реализует интерфейс Widget
func (c *CanvasRGB) InnerText() string {
	if !c.dirty && c.cached != "" {
		return c.cached
	}

	var b strings.Builder

	b.Grow(c.width * c.height * 10)

	if c.PixelSize == PixelTwoSymbol {
		lastClr := ColorRGB{}
		z := ColorRGB{}
		for y := 0; y < c.height; y++ {
			for x := 0; x < c.width; x++ {
				clr := c.pole[y][x]
				if lastClr != clr {
					if clr == z {
						b.WriteString("\033[0m")
					} else {
						b.WriteString("\033[48;2;")
						b.WriteString(strconv.Itoa(int(clr.R)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(clr.G)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(clr.B)))
						b.WriteString("m")
					}
					lastClr = clr
				}
				b.WriteString("  ")
			}
			b.WriteString("\r\n")
		}
	} else {
		lastBg := ColorRGB{}
		lastFr := ColorRGB{}
		z := ColorRGB{}

		b.WriteString("\033[30m")
		for y := 0; y < c.height/2; y++ {
			for x := 0; x < c.width; x++ {
				var bg ColorRGB
				if y*2+1 == c.height {
					bg = ColorRGB{}
				} else {
					bg = c.pole[y*2+1][x]
				}

				fr := c.pole[y*2][x]
				if lastBg != bg {
					if bg == z {
						b.WriteString("\033[40m")
					} else {
						b.WriteString("\033[48;2;")
						b.WriteString(strconv.Itoa(int(bg.R)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(bg.G)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(bg.B)))
						b.WriteString("m")
					}
					lastBg = bg
				}
				if lastFr != fr {
					if fr == z {
						b.WriteString("\033[30m")
					} else {
						b.WriteString("\033[38;2;")
						b.WriteString(strconv.Itoa(int(fr.R)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(fr.G)))
						b.WriteString(";")
						b.WriteString(strconv.Itoa(int(fr.B)))
						b.WriteString("m")
					}
					lastFr = fr
				}
				b.WriteString("▀")
			}
			b.WriteString("\r\n")
		}
	}

	b.WriteString("\033[0m")

	c.cached = b.String()
	c.dirty = false

	return c.cached
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

func (cnv *CanvasRGB) Load(img image.Image) {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := cnv.Width(), cnv.Height()

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			x0 := x * srcW / dstW
			x1 := (x + 1) * srcW / dstW
			y0 := y * srcH / dstH
			y1 := (y + 1) * srcH / dstH

			var r, g, b uint64
			count := 0
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					col := img.At(sx+bounds.Min.X, sy+bounds.Min.Y)
					rr, gg, bb, _ := col.RGBA()
					r += uint64(rr >> 8)
					g += uint64(gg >> 8)
					b += uint64(bb >> 8)
					count++
				}
			}
			if count > 0 {
				avgR := uint8(r / uint64(count))
				avgG := uint8(g / uint64(count))
				avgB := uint8(b / uint64(count))
				cnv.Draw(x, y, ColorRGB{avgR, avgG, avgB})
			}
		}
	}
	cnv.dirty = true
}

func (cnv *CanvasRGB) Get(x, y int) ColorRGB {
	if x < 0 || x >= cnv.width || y < 0 || y >= cnv.height {
		return ColorRGB{}
	}
	return cnv.pole[y][x]
}
