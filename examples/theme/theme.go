package main

import (
	"time"

	tui "github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose — Theme")

	label := tui.NewStaticLabel("Hello, TUI Compose!")
	button := tui.NewButton("Click me", nil)
	check1 := tui.NewCheck("Option 1")
	check2 := tui.NewCheck("Option 2")
	check3 := tui.NewCheck("Option 3")
	field := tui.NewInputField(20).WithPlaceholder("Type here...")
	link := tui.NewHyperlink("GitHub", "https://github.com/romanSPB15/tui-compose")
	spinner := extra.NewSpinner(extra.SpinnerBrailleReverse)
	chart := extra.NewBarChart().WithValues([]int{10, 25, 40, 55, 30, 15, 5})

	wnd.SetStyleFunc(func(w tui.Widget) {
		switch x := w.(type) {
		case *tui.Button:
			if x.GetText() == "GitHub" {
				x.WithStyle(tui.FrBlue | tui.Underline).
					WithFocusedStyle(tui.FrMagenta | tui.Underline)
			} else {
				x.WithStyle(tui.BgBlue|tui.FrWhite).
					WithFocusedStyle(tui.BgBrightBlue|tui.FrWhite|tui.Bold).
					WithPaddings(3, 1)
			}
		case *tui.InputField:
			x.WithStyle(tui.BgBlack | tui.FrWhite).
				WithFocusedStyle(tui.BgBlue | tui.FrBrightWhite)
		case *tui.Label:
			x.WithStyle(tui.Italic)
		case *tui.Check:
			x.WithStyle(tui.FrBrightCyan).
				WithCheckedStyle(tui.FrBrightGreen | tui.Bold)
		case *extra.BarChart:
			x.WithBarStyle(func(i, v int) tui.Style {
				if v >= 40 {
					return tui.FrRed
				}
				if v > 20 {
					return tui.FrYellow
				}
				return tui.FrGreen
			}).WithDataHeight(10).WithBarWidth(2).AutoScale()
		case *extra.Spinner:
			x.WithStyle(tui.FrBrightMagenta).Start(time.Second / 15)
		}
	})

	wnd.SetContent(tui.NewFrame(tui.NewVBox(
		label,
		tui.NewHBox(button, tui.NewVBox(
			check1, check2, check3,
		)),
		field,
		link,
		spinner,
		chart,
	).WithGap(1)).
		Rounded().
		WithTitle(tui.Title{
			Text:  "TUI Compose SetStyleFunc",
			Pos:   tui.TitleTopCenter,
			Style: tui.Bold | tui.FrYellow,
		}).
		WithTitle(tui.Title{
			Text: "v4",
			Pos:  tui.TitleBottomRight,
		}))

	wnd.Run()
}
