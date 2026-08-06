package main

import (
	"fmt"
	"time"

	tui "github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/input"
)

func main() {
	w := tui.NewWindow()
	w.SetTitle("TUI Compose")

	timeLabel := tui.NewDynamicLabel("Exit in 3", 20)

	goodbye := tui.NewPage(tui.NewVBox(
		tui.NewStaticLabel("Goodbye!").WithStyle(tui.FrYellow),
		timeLabel,
	))

	box := tui.NewVBox(
		tui.NewHBox(tui.NewStaticLabel("First Name:"), tui.NewInputField(20)),
		tui.NewHBox(tui.NewStaticLabel("Last Name:"), tui.NewInputField(20)),
		tui.NewCheck("I agree to the terms of use").WithStyle(tui.Italic|tui.FrRed),
		tui.NewButton("Submit", func() {
			w.Commit(goodbye.Open)

			go func() {
				for i := 3; i > 0; i-- {
					w.Commit(func() {
						timeLabel.SetText(fmt.Sprintf("Exit in %d", i))
					})

					time.Sleep(time.Second)
				}
				w.Quit()
			}()
		}).WithStyle(tui.Italic|tui.BgBrightCyan),
	)

	form := tui.NewPage(box)

	cnv := tui.NewCanvas(11, 11)

	started := false

	w.RegisterKeyHandler(func(ke *input.KeyboardEvent) {
		if started {
			return
		}
		started = true
		go func() {
			// анимация креста
			for i := range 11 {
				w.Commit(func() {
					cnv.Draw(i, i, tui.Black+1+tui.Color((i+1)/2))
				})
				time.Sleep(time.Second / 10)
			}
			for i := range 11 {
				w.Commit(func() {
					cnv.Draw(10-i, i, tui.Black+1+tui.Color((i+1)/2))
				})
				time.Sleep(time.Second / 10)
			}
			form.Open()
		}()
	})

	w.SetContent(cnv)

	w.Run()
}
