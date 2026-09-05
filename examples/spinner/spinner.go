package main

import (
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()

	wnd.SetContent(tui.NewFrame(tui.NewVBox(
		tui.NewHBox(
			extra.NewSpinner(extra.SpinnerUnderscore).Start(time.Second/2).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerBraille).Start(time.Second/2),
			extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/2),
			extra.NewSpinner(extra.SpinnerDots).Start(time.Second/2).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerLine).Start(time.Second/2),
		),
		tui.NewHBox(
			extra.NewSpinner(extra.SpinnerUnderscore).Start(time.Second/3),
			extra.NewSpinner(extra.SpinnerBraille).Start(time.Second/3),
			extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/3),
			extra.NewSpinner(extra.SpinnerDots).Start(time.Second/3).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerLine).Start(time.Second/3),
		),
		tui.NewHBox(
			extra.NewSpinner(extra.SpinnerUnderscore).Start(time.Second/5),
			extra.NewSpinner(extra.SpinnerBraille).Start(time.Second/5),
			extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/5),
			extra.NewSpinner(extra.SpinnerDots).Start(time.Second/5),
			extra.NewSpinner(extra.SpinnerLine).Start(time.Second/5).WithStyle(tui.FrBrightGreen),
		),
		tui.NewHBox(
			extra.NewSpinner(extra.SpinnerUnderscore).Start(time.Second/10),
			extra.NewSpinner(extra.SpinnerBraille).Start(time.Second/10).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/10).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerDots).Start(time.Second/10),
			extra.NewSpinner(extra.SpinnerLine).Start(time.Second/10),
		),
		tui.NewHBox(
			extra.NewSpinner(extra.SpinnerUnderscore).Start(time.Second/15),
			extra.NewSpinner(extra.SpinnerBraille).Start(time.Second/15).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/15).WithStyle(tui.FrBrightGreen),
			extra.NewSpinner(extra.SpinnerDots).Start(time.Second/15),
			extra.NewSpinner(extra.SpinnerLine).Start(time.Second/15),
		),
	)).WithTitle(tui.Title{Text: "Spinner", Pos: tui.TitleTopCenter, Style: tui.FrYellow}).Rounded())

	wnd.Run()
}
