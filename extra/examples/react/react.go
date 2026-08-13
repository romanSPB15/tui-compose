package main

import (
	"fmt"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/react"
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
				}).WithStyle(tui.BgRed),
				tui.NewButton("-", func() {
					if state.Count > 0 {
						app.Mutate(func(s *State) { s.Count-- })
					}
				}).WithStyle(tui.BgBlue),
			),
			tui.NewButton("Выход", func() {
				app.Quit()
			}).WithStyle(tui.FrBrightBlack),
		)).WithTitle(tui.Title{Text: "TUI Compose React"})
	})

	app.SetTitle("Моё приложение")
	app.Run()
}
