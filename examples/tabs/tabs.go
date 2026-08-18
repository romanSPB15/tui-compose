package main

import (
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - Tabs")

	wnd.SetContent(extra.NewTabs([]extra.Tab{
		{
			Title:      "First",
			TitleStyle: tui.FrYellow,
			Content:    tui.NewButton("First - exit Button", wnd.Quit).WithStyle(tui.BgBrightRed),
		},
		{
			Title:   "Second",
			Content: tui.NewCheck("Second - Check"),
		},
		{
			Title:   "Third",
			Content: tui.NewHBox(extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/15).WithStyle(tui.FrBrightMagenta), tui.NewStaticLabel("Loading forever...")),
		},
	}))

	wnd.Run()
}
