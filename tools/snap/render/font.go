package render

import (
	_ "embed"
	"image"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

//go:embed fonts/CascadiaCode-Regular.ttf
var fontData []byte

var (
	face       font.Face
	charWidth  int
	charHeight int
)

func init() {
	ttf, err := opentype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	face, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	m := face.Metrics()
	charWidth = m.Height.Ceil() / 2
	charHeight = m.Height.Ceil()
	if charWidth < 1 {
		charWidth = 8
	}
	if charHeight < 1 {
		charHeight = 14
	}
}

func renderRune(img *image.RGBA, x, y int, r rune, style cell.Style) {
	fg, bg := getColor(style)

	if bg != nil {
		for dy := 0; dy < charHeight; dy++ {
			for dx := 0; dx < charWidth; dx++ {
				img.Set(x+dx, y+dy, bg)
			}
		}
	}

	ascent := face.Metrics().Ascent.Ceil()
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y + ascent)},
	}
	d.DrawString(string(r))

	if style.Args&cell.Underline != 0 {
		uy := y + ascent + 2
		for dx := 0; dx < charWidth; dx++ {
			img.Set(x+dx, uy, fg)
		}
	}
}

func getColor(style cell.Style) (fg, bg color.Color) {
	fg = defaultTextForeground
	bg = defaultTextBackground

	if style.Fg != "" {
		switch {
		case strings.HasPrefix(style.Fg, "38;2;"):
			p := strings.Split(style.Fg, ";")
			if len(p) == 5 {
				r, _ := strconv.Atoi(p[2])
				g, _ := strconv.Atoi(p[3])
				b, _ := strconv.Atoi(p[4])
				fg = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			}
		case strings.HasPrefix(style.Fg, "38;5;"):
			p := strings.Split(style.Fg, ";")
			if len(p) == 3 {
				i, _ := strconv.Atoi(p[2])
				fg = ansi256ToRGB(i)
			}
		default:
			if c, ok := ansi16Fg[style.Fg]; ok {
				fg = c
			}
		}
	}

	if style.Bg != "" {
		switch {
		case strings.HasPrefix(style.Bg, "48;2;"):
			p := strings.Split(style.Bg, ";")
			if len(p) == 5 {
				r, _ := strconv.Atoi(p[2])
				g, _ := strconv.Atoi(p[3])
				b, _ := strconv.Atoi(p[4])
				bg = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
			}
		case strings.HasPrefix(style.Bg, "48;5;"):
			p := strings.Split(style.Bg, ";")
			if len(p) == 3 {
				i, _ := strconv.Atoi(p[2])
				bg = ansi256ToRGB(i)
			}
		default:
			if c, ok := ansi16Bg[style.Bg]; ok {
				bg = c
			}
		}
	}

	if style.Args&cell.Reverse != 0 {
		fg, bg = bg, fg
	}
	return
}

var ansi16Fg = map[string]color.Color{
	"30": color.RGBA{0, 0, 0, 255},
	"31": color.RGBA{255, 0, 0, 255},
	"32": color.RGBA{0, 255, 0, 255},
	"33": color.RGBA{255, 255, 0, 255},
	"34": color.RGBA{50, 50, 255, 255},
	"35": color.RGBA{255, 0, 255, 255},
	"36": color.RGBA{0, 255, 255, 255},
	"37": color.RGBA{255, 255, 255, 255},
	"90": color.RGBA{128, 128, 128, 255},
	"91": color.RGBA{255, 128, 128, 255},
	"92": color.RGBA{128, 255, 128, 255},
	"93": color.RGBA{255, 255, 128, 255},
	"94": color.RGBA{128, 128, 255, 255},
	"95": color.RGBA{255, 128, 255, 255},
	"96": color.RGBA{128, 255, 255, 255},
	"97": color.RGBA{255, 255, 255, 255},
}

var ansi16Bg = map[string]color.Color{
	"40":  color.RGBA{0, 0, 0, 255},
	"41":  color.RGBA{128, 0, 0, 255},
	"42":  color.RGBA{0, 128, 0, 255},
	"43":  color.RGBA{128, 128, 0, 255},
	"44":  color.RGBA{0, 0, 128, 255},
	"45":  color.RGBA{128, 0, 128, 255},
	"46":  color.RGBA{0, 128, 128, 255},
	"47":  color.RGBA{192, 192, 192, 255},
	"100": color.RGBA{128, 128, 128, 255},
	"101": color.RGBA{255, 0, 0, 255},
	"102": color.RGBA{0, 255, 0, 255},
	"103": color.RGBA{255, 255, 0, 255},
	"104": color.RGBA{0, 0, 255, 255},
	"105": color.RGBA{255, 0, 255, 255},
	"106": color.RGBA{0, 200, 200, 255},
	"107": color.RGBA{200, 200, 200, 255},
}

func ansi256ToRGB(idx int) color.Color {
	if idx < 16 {
		l := []color.RGBA{
			{0, 0, 0, 255}, {128, 0, 0, 255}, {0, 128, 0, 255}, {128, 128, 0, 255},
			{0, 0, 128, 255}, {128, 0, 128, 255}, {0, 128, 128, 255}, {192, 192, 192, 255},
			{128, 128, 128, 255}, {255, 0, 0, 255}, {0, 255, 0, 255}, {255, 255, 0, 255},
			{0, 0, 255, 255}, {255, 0, 255, 255}, {0, 255, 255, 255}, {255, 255, 255, 255},
		}
		if idx < len(l) {
			return l[idx]
		}
		return color.Black
	}
	if idx < 232 {
		r := (idx - 16) / 36
		g := ((idx - 16) % 36) / 6
		b := (idx - 16) % 6
		return color.RGBA{uint8(r * 51), uint8(g * 51), uint8(b * 51), 255}
	}
	v := uint8((idx-232)*10 + 8)
	return color.RGBA{v, v, v, 255}
}
