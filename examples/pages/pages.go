package main

import (
	"fmt"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
)

func main() {
	a := tui.NewWindow()

	var pgMain *tui.Page

	timer := tui.NewDynamicLabel("Выход через 3.00с", 20).WithStyle(tui.FrYellow)

	pg3 := tui.NewPage(timer)

	pg2 := tui.NewPage(tui.NewFrame(tui.NewVBox(
		tui.NewStaticLabel("Подтвердите действие").WithStyle(tui.Bold),
		tui.NewHBox(
			tui.NewButton("ОК", func() {
				pg3.Open(a)
				i := 300
				go func() {
					tick := time.NewTicker(time.Second / 50)
					for {
						select {
						case <-a.OnQuit():
							tick.Stop()
							return
						case <-tick.C:
							if i > 0 {
								i -= 2
								a.Commit(func() {
									timer.WithText(fmt.Sprintf("Выход через %d.%02dс", i/100, i%100))
								})
							} else {
								a.Quit()
							}
						}
					}
				}()
			}).WithStyle(tui.BgBlue).WithPaddings(3, 1),
			tui.NewButton("Отмена", func() {
				pgMain.Open(a)
			}).WithStyle(tui.BgRed).WithPaddings(3, 1),
		)).WithGap(1)).WithPaddings(1, 2).WithTitle(tui.Title{Text: "Диалог", Pos: tui.TitleTopCenter, Style: tui.Bold | tui.FrMagenta}).Rounded()).SetTitle("Page 2")

	btn := tui.NewButton("Удалить данные", func() {
		pg2.Open(a)
	}).WithStyle(tui.BgBrightRed).WithPaddings(2, 1)

	pgMain = tui.NewPage(tui.NewFrame(btn).Rounded().WithPaddings(2, 5).WithTitle(tui.Title{Text: "TUI Compose Page", Pos: tui.TitleTopCenter})).SetTitle("Page 1")

	pgMain.Open(a)

	a.Run()
}
