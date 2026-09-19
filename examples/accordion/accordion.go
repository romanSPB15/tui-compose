package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose - Accordion Demo")

	accordion := extra.NewAccordion(
		"Accordion",
		tui.NewStaticLabel("<you content here>").WithStyle(tui.FrCyan),
	)

	wnd.SetContent(tui.NewFrame(accordion).
		Rounded().
		WithTitle(tui.Title{Text: "Accordion Demo", Pos: tui.TitleTopCenter, Style: tui.FrMagenta}))

	wnd.Run()
}
