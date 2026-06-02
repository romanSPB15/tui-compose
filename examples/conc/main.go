package main

import (
	"fmt"
	"time"

	"github.com/romanSPB15/go-tui"
)

func main() {
	w := tui.NewWindow()
	lbl := tui.NewDynamicLabel("0", 9).ColorizeBackground(tui.Cyan)
	go func() {
		i := 1
		for {
			time.Sleep(time.Second)
			w.Do(func() {
				lbl.Text = fmt.Sprint(i)
				i++
				w.RedrawWidget(0)
			})
		}
	}()
	w.AddWidgets(lbl)
	w.Run()
}
