package main

import (
	tui "github.com/romanSPB15/tui-compose/v3"
)

func main() {
	w := tui.NewWindow()
	w.SetTitle("TUI Compose - Frames")

	content := func() tui.Widget {
		return tui.NewStaticLabel("TUI Compose Frame")
	}

	fr1 := tui.NewFrame(content())
	fr2 := tui.NewFrame(content()).WithTitle(tui.Title{Text: "Заголовок"})
	fr3 := tui.NewFrame(content()).
		WithTitle(tui.Title{Text: "Заголовок", Pos: tui.TitleTopCenter, Style: tui.FrRed}).
		WithTitle(tui.Title{Text: "6 позиций", Pos: tui.TitleBottomRight, Style: tui.FrYellow | tui.Bold})

	fr4 := tui.NewFrame(content()).Rounded().WithTitle(tui.Title{Text: "Скругления", Pos: tui.TitleTopCenter, Style: tui.FrYellow | tui.Bold})
	fr5 := tui.NewFrame(content()).Heavy().WithTitle(tui.Title{Text: "Толстая рамка", Pos: tui.TitleTopCenter})
	fr6 := tui.NewFrame(content()).Double().WithTitle(tui.Title{Text: "Двойная рамка", Pos: tui.TitleTopCenter})

	fr7 := tui.NewFrame(content()).ASCII().WithTitle(tui.Title{Text: "ASCII", Pos: tui.TitleTopCenter, Style: tui.FrBrightMagenta | tui.Bold})
	fr8 := tui.NewFrame(content()).BevelASCII().WithTitle(tui.Title{Text: "Фаски", Pos: tui.TitleTopCenter})
	fr9 := tui.NewFrame(content()).Custom('O', 'O', 'O', 'O', '─', '│').WithTitle(tui.Title{Text: "Custom", Pos: tui.TitleTopCenter, Style: tui.FrGreen | tui.Bold})

	w.SetContent(
		tui.NewFrame(
			tui.NewVBox(
				tui.NewHBox(fr1, fr2, fr3),
				tui.NewHBox(fr4, fr5, fr6),
				tui.NewHBox(fr7, fr8, fr9),
			)).Dashed(),
	)

	w.Run()
}
