package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	ta := extra.NewTextArea(80, 20).
		WithLineNumbers(true).
		WithLineNumberStyle(tui.FrBrightBlack).
		WithText(`package main

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
`).
		WithOnChanged(func(s string) {})
	w := tui.NewWindow()
	w.SetContent(tui.NewFrame(ta))
	w.Run()
}
