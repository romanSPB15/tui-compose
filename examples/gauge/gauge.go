package main

import "github.com/romanSPB15/tui-compose/v4"

func main() {
	wnd := tui.NewWindow()
	wnd.SetContent(tui.NewVBox(
		tui.NewGauge(20).WithValue(0.6),
	))
	wnd.Run()
}
