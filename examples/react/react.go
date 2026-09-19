package main

import (
	"fmt"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/react"
)

type State struct {
	Count int
}

func main() {
	app := react.New(State{}, func(app *react.App[State], state State) tui.Widget {
		return tui.NewFrame(tui.NewVBox(
			tui.NewStaticLabel(fmt.Sprintf("Count: %d", state.Count)).
				WithStyle(tui.FrCyan|tui.Bold),
			tui.NewHBox(
				tui.NewButton("+", func() {
					app.Mutate(func(s *State) { s.Count++ })
				}).WithStyle(tui.BgRed).WithPaddings(0, 0),
				tui.NewButton("-", func() {
					if state.Count > 0 {
						app.Mutate(func(s *State) { s.Count-- })
					}
				}).WithStyle(tui.BgBlue).WithPaddings(0, 0),
			),
			tui.NewButton("Выход", func() {
				app.Quit()
			}).WithStyle(tui.BgBrightBlack),
		).WithGap(1)).WithTitle(tui.Title{Text: "TUI Compose", Pos: tui.TitleTopCenter, Style: tui.FrBlue | tui.Bold}).Double().WithPaddings(0, 3)
	})

	app.SetTitle("Моё приложение")
	app.Run()
}
