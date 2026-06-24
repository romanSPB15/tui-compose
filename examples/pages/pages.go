package main

import (
	"github.com/romanSPB15/tui-compose/v3"
)

func main() {
	a := tui.NewWindow()

	pg2 := tui.NewPage(tui.NewStaticLabel("Это страница 2!")).SetTitle("Page 2")

	lbl := tui.NewStaticLabel("Страница 1")
	btn := tui.NewButton("Перейти на страницу 2", func() {
		pg2.Open()
	})

	box := tui.NewVBox(lbl, btn)

	pgMain := tui.NewPage(box).SetTitle("Page 1")

	pgMain.Open()

	a.Run()
}
