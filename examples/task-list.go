package main

import (
	"github.com/romanSPB15/tui-compose/v3"
)

func main() {
	a := tui.NewWindow()

	a.SetTitle("TUI Compose")

	label := tui.NewDynamicLabel("", 30).ColorizeForeground(tui.Red)

	count := 0

	checkIfAll := func() {
		if count == 3 {
			a.Do(func() {
				label.Text = "Все дела сделаны"
				a.Redraw()
			})
		}
	}

	onChange := func(b bool) {
		if b {
			count++
		}
		checkIfAll()
	}

	check := tui.NewCheck("Дело 1")

	check2 := tui.NewCheck("Дело 2")

	check3 := tui.NewCheck("Дело 3")

	check.OnChanged, check2.OnChanged, check3.OnChanged = onChange, onChange, onChange

	box := tui.NewVBox(check, check2, check3, label)
	a.SetContent(box)

	a.Run()
}
