package main

import (
	"github.com/romanSPB15/tui-compose/v2"
)

func main() {
	w := tui.NewWindow()
	cnv := tui.NewCanvas(10, 10)

	w.DisableFocusChange()

	x, y := 0, 0
	cnv.Draw(0, 0, tui.White)

	// Назначаем обработчики
	w.RegisterKeyHandler(tui.KeyArrowDown, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) // Замазываем следы
		y++
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(tui.KeyArrowUp, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		y--
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(tui.KeyArrowLeft, func() {
		cnv.DrawAndRender(x, y, tui.NoColor)
		x--
		cnv.DrawAndRender(x, y, tui.White)
	})
	w.RegisterKeyHandler(tui.KeyArrowRight, func() {
		cnv.DrawAndRender(x, y, tui.NoColor) //
		x++
		cnv.DrawAndRender(x, y, tui.White)
	})

	w.AddWidgets(cnv)

	w.Run()
}
