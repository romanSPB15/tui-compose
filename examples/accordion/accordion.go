package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - Accordion")

	accordion := extra.NewAccordion(
		"Accordion. Click to open",
		tui.NewStaticLabel("<you content here>").WithStyle(tui.FrCyan),
	)

	wnd.SetContent(tui.NewFrame(accordion).
		Rounded().
		WithTitle(tui.Title{Text: "Accordion Demo", Pos: tui.TitleTopCenter, Style: tui.FrMagenta}))

	wnd.Run()
}
