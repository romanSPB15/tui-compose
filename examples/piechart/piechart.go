package main

import (
	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - PieChart")

	data := []extra.PieData{
		{Label: "TUI Compose", Value: 50, Color: tui.FrBrightBlue},
		{Label: "Bubble Tea", Value: 25, Color: tui.FrBrightBlack},
		{Label: "tview", Value: 10, Color: tui.FrWhite},
		{Label: "termui", Value: 10, Color: tui.FrCyan},
		{Label: "other", Value: 5, Color: tui.FrRed},
	}

	pie := extra.NewPieChart(data).
		WithRadius(20).
		WithShowPercent(true).
		WithShowLegend(true).
		WithValueStyle(tui.FrBrightWhite | tui.Bold)

	wnd.SetContent(tui.NewFrame(pie).
		Rounded().
		WithTitle(tui.Title{Text: "TUI Compose vs others", Pos: tui.TitleTopCenter, Style: tui.FrYellow}))

	wnd.Run()
}
