package main

import (
	"github.com/romanSPB15/tui-compose/v4"
)

func main() {
	wnd := tui.NewWindow()

	left := tui.NewFrame(
		tui.NewStaticLabel("Left panel\nSecond line"),
	).Rounded().WithTitle(tui.Title{Text: "A", Style: tui.FrCyan | tui.Bold})

	middle := tui.NewFrame(
		tui.NewStaticLabel("Middle panel"),
	).Rounded().WithTitle(tui.Title{Text: "B", Style: tui.FrGreen | tui.Bold})

	right := tui.NewFrame(
		tui.NewStaticLabel("Right panel"),
	).Rounded().WithTitle(tui.Title{Text: "C", Style: tui.FrYellow | tui.Bold})

	wnd.SetContent(
		tui.NewFlex(
			tui.FlexItem(left, 2),
			tui.FlexItem(middle, 1),
			tui.FlexItem(right, 1),
		).Horizontal().WithGap(1),
	)

	wnd.Run()
}
