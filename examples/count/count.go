package main

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v4"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("Моё приложение")

	label := tui.NewStaticLabel("0").WithStyle(tui.FrCyan)

	btnQuit := tui.NewButton("Выход", func() {
		wnd.Quit()
	}).WithStyle(tui.BgBrightBlack).WithPaddings(3, 0)

	v := 0

	btnAdd := tui.NewButton("+", func() {
		v++
		wnd.Commit(func() {
			label.SetText(strconv.Itoa(v))
		})
	}).WithStyle(tui.BgRed)

	btnSub := tui.NewButton("-", func() {
		if v > 0 {
			v--
			wnd.Commit(func() {
				label.SetText(strconv.Itoa(v))
			})
		}
	}).WithStyle(tui.BgBlue)

	box := tui.NewVBox(label, tui.NewHBox(btnAdd, btnSub), btnQuit).WithGap(1)
	wnd.SetContent(tui.NewFrame(box).
		WithTitle(tui.Title{Text: "TUI Compose", Pos: tui.TitleTopCenter, Style: tui.FrYellow | tui.Bold}).
		WithPaddings(0, 2).Rounded())

	wnd.Run()
}
