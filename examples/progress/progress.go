package main

import (
	tui "github.com/romanSPB15/tui-compose/v3"
)

func main() {
	w := tui.NewWindow()
	w.SetTitle("TUI Compose - Frames")

	pr1 := tui.NewColorProgress(20, tui.Blue, tui.BrightBlack)
	pr1.SetValue(0.5)

	pr2 := tui.NewColorProgress(20, tui.Green, tui.BrightBlack)
	pr2.SetValue(0.75)

	pr3 := tui.NewTextProgress(20, '#', '-')
	pr3.SetValue(0.3)

	pr4 := tui.NewTextProgress(20, '#', '_')
	pr4.SetValue(0.6)

	w.SetContent(tui.NewVBox(pr1, pr3, pr2, pr4))

	w.Run()
}
