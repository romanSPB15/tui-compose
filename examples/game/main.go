package main

import (
	"github.com/eiannone/keyboard"
	"github.com/romanSPB15/go-tui"
)

func main() {
	w := tui.NewWindow()
	cnv := tui.NewCanvas(10, 10)

	x, y := 0, 0
	cnv.Draw(0, 0, tui.White)

	// Назначаем обработчики
	w.RegisterKeyHandler(keyboard.KeyArrowDown, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) // Замазываем следы
		y++
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(keyboard.KeyArrowUp, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		y--
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(keyboard.KeyArrowLeft, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		x--
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(keyboard.KeyArrowRight, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) //
		x++
		cnv.DrawAndRender(x, y, tui.White)
	})

	w.AddWidgets(cnv)

	w.Run()
}
