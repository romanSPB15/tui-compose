package main

import (
	"github.com/romanSPB15/tui-compose/v3"
)

func main() {
	wnd := tui.NewWindow()

	name := tui.NewInputField(30).WithPlaceholder("Введите имя...").WithStyle(tui.BgBrightBlue)
	com := tui.NewInputField(30).WithPlaceholder("Введите комментарий...").WithStyle(tui.BgBrightBlue)

	pro := tui.NewCheck("Pro").WithStyle(tui.FrMagenta)

	ok, cancel := tui.NewButton("ОК", func() {
		nameValue := name.Text
		comValue := com.Text
		proValue := pro.State()
		_ = nameValue
		_ = comValue
		_ = proValue
		// ...
	}).WithStyle(tui.BgBrightCyan), tui.NewButton("Отмена", func() {
		// ...
	}).WithStyle(tui.BgBrightRed)

	fr := tui.NewFrame((tui.NewVBox(name, com, pro, tui.NewHBox(ok, cancel)))).WithTitle(tui.Title{Text: "Регистрация", Pos: tui.TitleTopCenter}).Rounded()

	wnd.SetContent(fr)
	wnd.Run()
}
