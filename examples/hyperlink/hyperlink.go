package main

import "github.com/romanSPB15/tui-compose/v3"

func main() {
	wnd := tui.NewWindow()

	wnd.SetContent(tui.NewFrame(tui.NewHyperlink("Open tui-compose GitHub", "https://github.com/romanSPB15/tui-compose").WithStyle(tui.FrBlue | tui.Underline)).Rounded())
	wnd.Run()
}
