package main

import (
	"github.com/romanSPB15/tui-compose/v4"
)

func main() {
	wnd := tui.NewWindow()

	name := tui.NewInputField(30).WithPlaceholder("Введите имя...").WithStyle(tui.BgBlue).WithFocusedStyle(tui.BgBrightBlue)
	com := tui.NewInputField(30).WithPlaceholder("Введите комментарий...").WithStyle(tui.BgBlue).WithFocusedStyle(tui.BgBrightBlue)

	pro := tui.NewCheck("Pro").WithStyle(tui.FrMagenta)

	ok, cancel := tui.NewButton("ОК", func() {
		nameValue := name.Text
		comValue := com.Text
		proValue := pro.State()
		_ = nameValue
		_ = comValue
		_ = proValue
		// ...
	}).WithStyle(tui.BgCyan).WithPaddings(2, 1).WithFocusedStyle(tui.BgBrightCyan), tui.NewButton("Отмена", func() {
		// ...
	}).WithStyle(tui.BgRed).WithPaddings(2, 1).WithFocusedStyle(tui.BgBrightRed)

	fr := tui.NewFrame(tui.NewVBox(name, com, pro, tui.NewHBox(ok, cancel)).WithGap(1)).WithTitle(tui.Title{Text: "Регистрация", Pos: tui.TitleTopCenter}).Rounded()

	wnd.SetContent(fr)
	wnd.Run()
}
