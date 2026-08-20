package main

import (
	"image/png"
	"log"
	"os"

	"github.com/romanSPB15/tui-compose/v3"
)

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

	var rows []tui.Widget
	for i := 0; i < len(modes); i += 3 {
		rowWidgets := make([]tui.Widget, 0, 3)
		for j := 0; j < 3 && i+j < len(modes); j++ {
			m := modes[i+j]
			imageWidget := tui.NewImage().LoadImage(img, m.mode)
			frame := tui.NewFrame(imageWidget).
				Rounded().
				WithTitle(tui.Title{
					Text:  m.name,
					Pos:   tui.TitleTopCenter,
					Style: tui.FrYellow,
				})
			rowWidgets = append(rowWidgets, frame)
		}
		rows = append(rows, tui.NewHBox(rowWidgets...).WithGap(2))
	}

	content := tui.NewVBox(rows...).WithGap(1)
	wnd.SetContent(content)
	wnd.Run()
}
