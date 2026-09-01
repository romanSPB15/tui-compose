package main

import (
	"image/png"
	"log"
	"os"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func makeBraile() [][]bool {
	w := 80
	h := 80
	res := make([][]bool, h)
	for i := range res {
		res[i] = make([]bool, w)
	}

	for y := range h {
		y2 := y
		if y2 > h/2 {
			y2 = h - y2
		}
		for x := w/2 - y2; x < w/2+y2; x++ {
			res[y][x] = true
		}
	}

	for y := range h {
		y2 := y
		if y2 > h/2 {
			y2 = h - y2
		}
		y2 /= 2
		for x := w/2 - y2; x < w/2+y2; x++ {
			res[y][x] = false
		}
	}
	return res
}

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.4 - Image Modes")

	file, err := os.Open("examples/image/image.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	img = tui.ScaleToWidth(img, 20)

	modes := []struct {
		name string
		mode tui.LoadMode
	}{
		{"TrueColor + HalfSymbol", tui.PaletteTrueColor | tui.HalfSymbol},
		{"TrueColor + OneSymbol", tui.PaletteTrueColor | tui.OneSymbol},
		{"TrueColor + TwoSymbol", tui.PaletteTrueColor | tui.TwoSymbol},
		{"16Color + HalfSymbol", tui.Palette16Color | tui.HalfSymbol},
		{"16Color + OneSymbol", tui.Palette16Color | tui.OneSymbol},
		{"16Color + TwoSymbol", tui.Palette16Color | tui.TwoSymbol},
	}

	var tabs []extra.Tab
	for i := 0; i < len(modes); i++ {
		m := modes[i]
		imageWidget := tui.NewImage().LoadImage(img, m.mode)

		tabs = append(tabs, extra.Tab{
			Title: "  " + m.name + "  ",
			Content: tui.NewVBox(
				tui.NewStaticLabel("——————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————————"),
				imageWidget),
			TitleStyle: tui.BgBlue,
		})
	}

	tabs = append(tabs, extra.Tab{
		Title:      "  Braile  ",
		TitleStyle: tui.BgBlue,
		Content:    tui.NewVBox(tui.NewImage().LoadBraille(makeBraile(), 0)),
	})

	wnd.SetContent(tui.NewFrame(extra.NewTabs(tabs)))
	wnd.Run()
}
