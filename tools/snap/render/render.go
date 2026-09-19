package render

import (
	"image"
	"image/color"
	"math"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

var (
	windowBackground      = color.RGBA{20, 20, 20, 255}
	defaultTextForeground = color.RGBA{255, 255, 255, 255}
	defaultTextBackground = windowBackground

	PaddingH = 32
	PaddingV = 32

	// MarginSize  = 40
	MarginColor = color.RGBA{107, 80, 255, 255}
	// RoundRadius = 15
	MarginSize  = 0
	RoundRadius = 0

	MacOS       = false
	macOSColors = [3]color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{255, 255, 0, 255},
		color.RGBA{0, 255, 0, 255},
	}
	macOSStartX   = 20
	macOSStartY   = 20
	macOSGap      = 24
	macOSDiameter = 12
)

func Frame(cells [][]cell.Cell) image.Image {
	if len(cells) == 0 || len(cells[0]) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 300, 300))
	}

	if len(cells) == 0 || len(cells[0]) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 300, 300))
	}

	endX, endY := -1, -1
	for y, row := range cells {
		for x := len(row) - 1; x >= 0; x-- {
			if row[x].Char != ' ' {
				if x > endX {
					endX = x
				}
				if y > endY {
					endY = y
				}
				break
			}
		}
	}

	if endX < 0 || endY < 0 {
		return image.NewRGBA(image.Rect(0, 0, 300, 300))
	}

	cells = cells[:endY+1]
	for i := range cells {
		if endX+1 <= len(cells[i]) {
			cells[i] = cells[i][:endX+1]
		}
	}
	if len(cells) == 0 || len(cells[0]) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 300, 300))
	}

	h, w := len(cells), len(cells[0])
	imgW := w*charWidth + PaddingH*2 + MarginSize*2
	if imgW%2 == 1 {
		imgW++
	}
	imgH := h*charHeight + PaddingV*2 + MarginSize*2
	if imgH%2 == 1 {
		imgH++
	}
	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))

	for y := 0; y < imgH; y++ {
		for x := 0; x < imgW; x++ {
			img.Set(x, y, MarginColor)
		}
	}

	left, top := MarginSize, MarginSize
	right, bottom := imgW-MarginSize-1, imgH-MarginSize-1
	r := RoundRadius
	r2 := r * r

	for y := top; y <= bottom; y++ {
		for x := left; x <= right; x++ {
			corner := false
			if x < left+r && y < top+r {
				dx, dy := x-(left+r), y-(top+r)
				corner = dx*dx+dy*dy > r2
			}
			if x > right-r && y < top+r {
				dx, dy := x-(right-r), y-(top+r)
				corner = dx*dx+dy*dy > r2
			}
			if x < left+r && y > bottom-r {
				dx, dy := x-(left+r), y-(bottom-r)
				corner = dx*dx+dy*dy > r2
			}
			if x > right-r && y > bottom-r {
				dx, dy := x-(right-r), y-(bottom-r)
				corner = dx*dx+dy*dy > r2
			}
			if !corner {
				img.Set(x, y, windowBackground)
			}
		}
	}

	for y, row := range cells {
		for x, c := range row {
			renderRune(img,
				x*charWidth+PaddingH+MarginSize,
				y*charHeight+PaddingV+MarginSize,
				c.Char, c.Style)
		}
	}

	if MacOS {
		drawCircle(img, macOSStartX+MarginSize, macOSStartY+MarginSize, macOSDiameter, macOSColors[0])
		drawCircle(img, macOSStartX+macOSGap+MarginSize, macOSStartY+MarginSize, macOSDiameter, macOSColors[1])
		drawCircle(img, macOSStartX+macOSGap*2+MarginSize, macOSStartY+MarginSize, macOSDiameter, macOSColors[2])
	}

	return img
}

func drawCircle(img *image.RGBA, cx, cy, diameter int, col color.Color) {
	if diameter <= 0 {
		return
	}
	rad := diameter / 2
	b := img.Bounds()
	for y := cy - rad; y <= cy+rad; y++ {
		if y < b.Min.Y || y >= b.Max.Y {
			continue
		}
		dy := float64(y - cy)
		half := int(math.Sqrt(float64(rad*rad) - dy*dy))
		x1, x2 := cx-half, cx+half
		if x1 < b.Min.X {
			x1 = b.Min.X
		}
		if x2 >= b.Max.X {
			x2 = b.Max.X - 1
		}
		for x := x1; x <= x2; x++ {
			img.Set(x, y, col)
		}
	}
}
