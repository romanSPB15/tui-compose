package tui

import (
	"image"
	"image/color"
	"strconv"

	"github.com/romanSPB15/tui-compose/v3/builder"
	"github.com/romanSPB15/tui-compose/v3/cell"
)

type Image [][]cell.Cell

func NewImage() Image {
	return nil
}

func (iw Image) Render(buf [][]cell.Cell) {
	copy(buf, iw)
}

func (iw Image) Width() int {
	if len(iw) == 0 {
		return 0
	}
	return len(iw[0])
}

func (iw Image) Height() int {
	return len(iw)
}

type LoadMode uint16

const (
	Palette16Color LoadMode = 1 << iota
	PaletteTrueColor
	HalfSymbol LoadMode = 1<<iota + 4
	OneSymbol
	TwoSymbol
)

func makeRGBANSIBg(r, g, b uint8, buf *builder.Builder) string {
	buf.Reset()
	buf.WriteString("48;2;")
	buf.WriteString(strconv.Itoa(int(r)))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(int(g)))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(int(b)))
	return buf.StringCopy()
}

func makeRGBANSIFg(r, g, b uint8, buf *builder.Builder) string {
	buf.Reset()
	buf.WriteString("38;2;")
	buf.WriteString(strconv.Itoa(int(r)))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(int(g)))
	buf.WriteByte(';')
	buf.WriteString(strconv.Itoa(int(b)))
	return buf.StringCopy()
}

func nearestANSI16Fg(r, g, b uint8) string {
	// таблица RGB для 16 цветов (индексы 0..15)
	palette := [16][3]uint8{
		{0, 0, 0},       // 0 black
		{128, 0, 0},     // 1 red
		{0, 128, 0},     // 2 green
		{128, 128, 0},   // 3 yellow
		{0, 0, 128},     // 4 blue
		{128, 0, 128},   // 5 magenta
		{0, 128, 128},   // 6 cyan
		{192, 192, 192}, // 7 white
		{128, 128, 128}, // 8 bright black
		{255, 0, 0},     // 9 bright red
		{0, 255, 0},     // 10 bright green
		{255, 255, 0},   // 11 bright yellow
		{0, 0, 255},     // 12 bright blue
		{255, 0, 255},   // 13 bright magenta
		{0, 255, 255},   // 14 bright cyan
		{255, 255, 255}, // 15 bright white
	}

	bestIdx := 0
	bestDist := int(^uint(0) >> 1) // max int
	for i, c := range palette {
		dr := int(r) - int(c[0])
		dg := int(g) - int(c[1])
		db := int(b) - int(c[2])
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}

	// ANSI код: 30..37 для индексов 0..7, 90..97 для 8..15
	if bestIdx < 8 {
		return strconv.Itoa(30 + bestIdx)
	}
	return strconv.Itoa(90 + (bestIdx - 8))
}

func nearestANSI16Bg(r, g, b uint8) string {
	// таблица RGB для 16 цветов (индексы 0..15)
	palette := [16][3]uint8{
		{0, 0, 0},       // 0 black
		{128, 0, 0},     // 1 red
		{0, 128, 0},     // 2 green
		{128, 128, 0},   // 3 yellow
		{0, 0, 128},     // 4 blue
		{128, 0, 128},   // 5 magenta
		{0, 128, 128},   // 6 cyan
		{192, 192, 192}, // 7 white
		{128, 128, 128}, // 8 bright black
		{255, 0, 0},     // 9 bright red
		{0, 255, 0},     // 10 bright green
		{255, 255, 0},   // 11 bright yellow
		{0, 0, 255},     // 12 bright blue
		{255, 0, 255},   // 13 bright magenta
		{0, 255, 255},   // 14 bright cyan
		{255, 255, 255}, // 15 bright white
	}

	bestIdx := 0
	bestDist := int(^uint(0) >> 1) // max int
	for i, c := range palette {
		dr := int(r) - int(c[0])
		dg := int(g) - int(c[1])
		db := int(b) - int(c[2])
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = i
		}
	}

	// ANSI код: 30..37 для индексов 0..7, 90..97 для 8..15
	if bestIdx < 8 {
		return strconv.Itoa(40 + bestIdx)
	}
	return strconv.Itoa(100 + (bestIdx - 8))
}

func (iw Image) LoadImage(img image.Image, pal LoadMode) Image {
	sb := &builder.Builder{}
	switch pal {
	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// RGB(True Color)

	case PaletteTrueColor | HalfSymbol:
		iw = make(Image, img.Bounds().Dy()/2)
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx())
		}
		for y := range img.Bounds().Dy() / 2 {
			for x := range img.Bounds().Dx() {
				rT, gT, bT, _ := img.At(x, y*2).RGBA()
				rT = rT / 256
				gT = gT / 256
				bT = bT / 256
				if rT > 255 {
					rT = 255
				}
				if gT > 255 {
					gT = 255
				}
				if bT > 255 {
					bT = 255
				}

				rB, gB, bB, _ := img.At(x, y*2+1).RGBA()
				rB = rB / 256
				gB = gB / 256
				bB = bB / 256
				if rB > 255 {
					rB = 255
				}
				if gB > 255 {
					gB = 255
				}
				if bB > 255 {
					bB = 255
				}

				iw[y][x].Char = '▀'
				iw[y][x].Style = cell.Style{
					Bg: makeRGBANSIBg(uint8(rB), uint8(gB), uint8(bB), sb),
					Fg: makeRGBANSIFg(uint8(rT), uint8(gT), uint8(bT), sb),
				}
			}
		}
	case PaletteTrueColor | OneSymbol:
		iw = make(Image, img.Bounds().Dy())
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx())
		}
		for y := range img.Bounds().Dy() {
			for x := range img.Bounds().Dx() {
				r, g, b, _ := img.At(x, y).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}

				iw[y][x].Char = ' '
				iw[y][x].Style = cell.Style{
					Bg: makeRGBANSIBg(uint8(r), uint8(g), uint8(b), sb),
				}
			}
		}
	case PaletteTrueColor | TwoSymbol:
		iw = make(Image, img.Bounds().Dy())
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx()*2)
		}
		for y := range img.Bounds().Dy() {
			for x := range img.Bounds().Dx() {
				r, g, b, _ := img.At(x, y).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}

				bg := makeRGBANSIBg(uint8(r), uint8(g), uint8(b), sb)

				iw[y][x*2].Char = ' '
				iw[y][x*2].Style = cell.Style{
					Bg: bg,
				}

				iw[y][x*2+1].Char = ' '
				iw[y][x*2+1].Style = cell.Style{
					Bg: bg,
				}
			}

		}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 16 цветов

	case Palette16Color | HalfSymbol:
		iw = make(Image, img.Bounds().Dy()/2)
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx())
		}
		for y := range img.Bounds().Dy() / 2 {
			for x := range img.Bounds().Dx() {
				rT, gT, bT, _ := img.At(x, y*2).RGBA()
				rT = rT / 256
				gT = gT / 256
				bT = bT / 256
				if rT > 255 {
					rT = 255
				}
				if gT > 255 {
					gT = 255
				}
				if bT > 255 {
					bT = 255
				}

				rB, gB, bB, _ := img.At(x, y*2+1).RGBA()
				rB = rB / 256
				gB = gB / 256
				bB = bB / 256
				if rB > 255 {
					rB = 255
				}
				if gB > 255 {
					gB = 255
				}
				if bB > 255 {
					bB = 255
				}

				iw[y][x].Char = '▀'
				iw[y][x].Style = cell.Style{
					Bg: nearestANSI16Bg(uint8(rB), uint8(gB), uint8(bB)),
					Fg: nearestANSI16Fg(uint8(rT), uint8(gT), uint8(bT)),
				}
			}
		}

	case Palette16Color | OneSymbol:
		iw = make(Image, img.Bounds().Dy())
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx())
		}
		for y := range img.Bounds().Dy() {
			for x := range img.Bounds().Dx() {
				r, g, b, _ := img.At(x, y).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}

				iw[y][x].Char = ' '
				iw[y][x].Style = cell.Style{
					Bg: nearestANSI16Bg(uint8(r), uint8(g), uint8(b)),
				}
			}
		}
	case Palette16Color | TwoSymbol:
		iw = make(Image, img.Bounds().Dy())
		for i := range iw {
			iw[i] = make([]cell.Cell, img.Bounds().Dx()*2)
		}
		for y := range img.Bounds().Dy() {
			for x := range img.Bounds().Dx() {
				r, g, b, _ := img.At(x, y).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}

				bg := nearestANSI16Bg(uint8(r), uint8(g), uint8(b))

				iw[y][x*2].Char = ' '
				iw[y][x*2].Style = cell.Style{
					Bg: bg,
				}

				iw[y][x*2+1].Char = ' '
				iw[y][x*2+1].Style = cell.Style{
					Bg: bg,
				}
			}
		}

	}
	return iw
}

func DownscaleImage(img image.Image, newWidth, newHeight int) *image.RGBA {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	if newWidth <= 0 || newHeight <= 0 {
		return nil
	}
	if srcW == 0 || srcH == 0 {
		return image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Вычисляем область исходного изображения, соответствующую этому пикселю
			x0 := x * srcW / newWidth
			x1 := (x + 1) * srcW / newWidth
			y0 := y * srcH / newHeight
			y1 := (y + 1) * srcH / newHeight

			var r, g, b, a uint64
			count := 0
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					col := img.At(sx+bounds.Min.X, sy+bounds.Min.Y)
					rr, gg, bb, aa := col.RGBA()
					r += uint64(rr >> 8)
					g += uint64(gg >> 8)
					b += uint64(bb >> 8)
					a += uint64(aa >> 8)
					count++
				}
			}
			if count > 0 {
				avgR := uint8(r / uint64(count))
				avgG := uint8(g / uint64(count))
				avgB := uint8(b / uint64(count))
				avgA := uint8(a / uint64(count))
				dst.SetRGBA(x, y, color.RGBA{avgR, avgG, avgB, avgA})
			}
		}
	}
	return dst
}

func ScaleToHeight(img image.Image, targetHeight int) *image.RGBA {
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	coff := float64(dy) / float64(targetHeight)
	newW := int(float64(dx) / coff)
	newH := targetHeight
	return DownscaleImage(img, newW, newH)
}

func ScaleToWidth(img image.Image, targetWidth int) *image.RGBA {
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
	coff := float64(dx) / float64(targetWidth)
	newW := targetWidth
	newH := int(float64(dy) / coff)
	return DownscaleImage(img, newW, newH)
}

// brailleChar возвращает символ Брайля для блока 4 строки × 2 колонки.
// bits — матрица bool размером 4×2 (rows=4, cols=2).
func brailleChar(bits [4][2]bool) rune {
	code := 0x2800 // база для символов Брайля

	if bits[0][0] {
		code |= 1 << 0
	}
	if bits[1][0] {
		code |= 1 << 1
	}
	if bits[2][0] {
		code |= 1 << 2
	}
	if bits[0][1] {
		code |= 1 << 3
	}
	if bits[1][1] {
		code |= 1 << 4
	}
	if bits[2][1] {
		code |= 1 << 5
	}
	if bits[3][0] {
		code |= 1 << 6
	}
	if bits[3][1] {
		code |= 1 << 7
	}

	return rune(code)
}

// LoadBraille загружает матрицу bool как изображение в стиле Брайля.
func (iw Image) LoadBraille(data [][]bool, style Style) Image {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil
	}
	rows := len(data) / 4
	cols := len(data[0]) / 2
	if rows == 0 || cols == 0 {
		return nil
	}

	data = data[:rows*4]
	for i := range data {
		data[i] = data[i][:cols*2]
	}

	s := ConvertToCellStyle(style)

	iw = make(Image, rows)
	for y := 0; y < rows; y++ {
		iw[y] = make([]cell.Cell, cols)
		for x := 0; x < cols; x++ {
			var block [4][2]bool
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					block[dy][dx] = data[y*4+dy][x*2+dx]
				}
			}
			iw[y][x] = cell.Cell{
				Char:  brailleChar(block),
				Style: s,
			}
		}
	}
	return iw
}
