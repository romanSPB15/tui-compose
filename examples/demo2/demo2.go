package main

import (
	"image/png"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/builder"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func ReadLogoAsBrailleMatrix() [][]bool {
	f, err := os.Open("./examples/demo2/logo.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	mat := make([][]bool, h)
	for y := 0; y < h; y++ {
		mat[y] = make([]bool, w)
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (r + g + b) / 3
			mat[y][x] = gray < 32768
		}
	}
	return mat
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

func main() {
	w := tui.NewWindow()
	w.SetTitle("TUI Compose")
	img := tui.NewImage()
	img = img.LoadBraille(ReadLogoAsBrailleMatrix(), tui.Style(0))

	buf := builder.Builder{}
	for i := range len(img[0]) {
		for j := range len(img) {
			img[j][i].Style.Fg = makeRGBANSIFg(uint8(255-float64(i)*1.5), 0, uint8(i), &buf)
		}
	}

	ln := int(math.Round(math.Pi * 16))

	sparkData := make([]int, ln)
	for i := range sparkData {
		sparkData[i] = int((math.Sin(float64(i)/4) + 1) * 50)
	}
	spark := extra.NewSparkline().WithValues(sparkData[:19]).WithHeight(3).AutoScale().WithBarStyle(func(i, v int) tui.Style {
		if v > 90 {
			return tui.FrRed
		}
		if v > 60 {
			return tui.FrYellow
		}
		return tui.FrGreen
	})

	go func() {
		offset := 0
		for {
			select {
			case <-w.OnQuit():
				return
			default:
				w.Commit(func() {
					for i := range len(img[0]) {
						for j := range len(img) {
							r := uint8(128 + 127*math.Sin(float64(i+offset)/40))
							b := uint8(128 + 127*math.Cos(float64(i+offset)/40))
							img[j][i].Style.Fg = makeRGBANSIFg(r, 0, b, &buf)
						}
					}

					newVals := sparkData[(offset % (ln / 2)) : (offset%(ln/2))+ln/2]
					spark.WithValues(newVals)
				})
				time.Sleep(time.Second / 30)
				offset++
			}
		}
	}()

	v := 0.0
	g := tui.NewGauge(25).Blocks().WithValue(v)

	content := tui.NewVBox(
		img,
		tui.NewFrame(
			tui.NewHBox(
				tui.NewVBox(
					tui.NewHyperlink("GitHub", "https://github.com/romanSPB15/tui-compose").WithStyle(tui.FrBrightBlue|tui.Underline),
					tui.NewHyperlink("Site", "https://romanspb15.github.io/tui-compose/").WithStyle(tui.FrBrightBlue|tui.Underline),
				).WithGap(1),
				tui.NewVBox(
					g,
					tui.NewHBox(
						tui.NewStaticLabel("Loading...").WithStyle(tui.FrBrightBlack),
						extra.NewSpinner(extra.SpinnerBrailleReverse).WithStyle(tui.FrYellow).Start(time.Second/14),
					).WithGap(2),
				).WithGap(1),
				tui.NewButton("Click Me!", func() {
					g.WithValue(v)
					if v < 1 {
						v += 0.05
					} else {
						g.WithOnStyle(tui.FrGreen)
						g.WithLabel("OK")
						w.Redraw()
					}
				}).WithPaddings(3, 1).WithStyle(tui.BgYellow|tui.Bold),
				tui.NewButton("Quit", w.Quit).WithPaddings(5, 1).WithStyle(tui.BgRed|tui.Bold),
				spark,
			).WithGap(5),
		).Heavy().WithBorderStyle(tui.FrCyan).WithPaddings(0, 2),
	)

	w.SetContent(
		tui.NewFrame(content).
			Rounded().
			WithTitle(tui.Title{Text: "TUI Compose", Style: tui.Bold | tui.FrGreen, Pos: tui.TitleTopCenter}).
			WithTitle(tui.Title{Text: "Golang terminal UIs", Style: tui.Bold | tui.FrBrightBlue, Pos: tui.TitleBottomLeft}).
			WithTitle(tui.Title{Text: "v3.5.0", Style: tui.FrBrightMagenta, Pos: tui.TitleBottomRight}).
			WithBorderStyle(tui.FrBrightBlack),
	)
	w.Run()
}
