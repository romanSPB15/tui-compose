package main

import (
	"fmt"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
)

func main() {
	wnd := tui.NewWindow()

	first := tui.NewStaticLabel("TUI Compose")
	second := tui.NewStaticLabel("Красный текст").WithStyle(tui.FrRed)
	third := tui.NewStaticLabel("Жирный текст с синим фоном").WithStyle(tui.Bold | tui.BgBlue)
	fourth := tui.NewDynamicLabel("Счёт: 0", 10).WithStyle(tui.Italic)

	go func() {
		i := 1

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-wnd.OnQuit(): // обработка выхода
				return
			case <-ticker.C:
				wnd.Commit(func() { // выполняем в UI-горутине
					fourth.SetText(fmt.Sprintf("Счёт: %d", i))
					i++
				})
			}
		}
	}()

	wnd.SetContent(tui.NewVBox(first, second, third, fourth))
	wnd.Run()
}
