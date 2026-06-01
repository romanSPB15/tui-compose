package main

import (
	"github.com/eiannone/keyboard"
	"github.com/romanSPB15/go-tui"
)

func main() {
	a := tui.NewApp()
	cnv := tui.NewCanvas(10, 10)

	x, y := 0, 0
	cnv.Draw(0, 0, tui.White)

	a.AddKeyHandler(keyboard.KeyArrowDown, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) // Замазываем следы
		y++
		cnv.DrawAndRender(x, y, tui.White)
	})
	a.AddKeyHandler(keyboard.KeyArrowUp, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		y--
		cnv.DrawAndRender(x, y, tui.White)
	})
	a.AddKeyHandler(keyboard.KeyArrowLeft, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		x--
		cnv.DrawAndRender(x, y, tui.White)
	})
	a.AddKeyHandler(keyboard.KeyArrowRight, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) //
		x++
		cnv.DrawAndRender(x, y, tui.White)
	})

	a.AddWidgets(cnv)

	a.Run()
}
