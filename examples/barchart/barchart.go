package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose - BarChart Demo")

	// данные
	data1 := []int{10, 25, 40, 55, 30, 15, 5}
	data2 := []int{14, 20, 30, 60, 35, 18, 6}
	data3 := []int{140, 200, 400, 600, 350}

	chart := extra.NewBarChart().
		WithValues(data1).
		AutoScale()

	chart2 := extra.NewBarChart().
		WithValues(data2).
		AutoScale().
		WithBarStyle(func(i, v int) tui.Style {
			if v > 35 {
				return tui.FrRed
			}
			if v > 25 {
				return tui.FrYellow
			}
			return tui.FrGreen
		})

	chart3 := extra.NewBarChart().
		WithDataHeight(9).
		WithValues(data1).
		AutoScale().
		WithBarStyle(func(i, v int) tui.Style {
			if v > 35 {
				return tui.FrRed
			}
			if v > 25 {
				return tui.FrYellow
			}
			return tui.FrGreen
		}).
		WithBarWidth(2)

	chart4 := extra.NewBarChart().
		WithDataHeight(9).
		WithValues(data1).
		WithScale(15).
		WithBarStyle(func(i, v int) tui.Style {
			if v > 35 {
				return tui.FrRed
			}
			if v > 25 {
				return tui.FrYellow
			}
			return tui.FrGreen
		}).
		WithBarWidth(2).
		WithTextStyle(func(i, v int) tui.Style {
			if v > 35 {
				return tui.FrBrightMagenta
			}
			if v > 25 {
				return tui.FrBrightYellow
			}
			return tui.FrCyan
		})

	chart5 := extra.NewBarChart().
		WithValues(data3).
		AutoScale().
		WithBarWidth(5).
		WithBarStyle(func(i, v int) tui.Style {
			if v > 400 {
				return tui.FrRed
			}
			if v > 300 {
				return tui.FrYellow
			}
			return tui.FrGreen
		}).
		WithFillRune('#')

	// Размещаем все диаграммы в окне
	wnd.SetContent(tui.NewFrame(tui.NewHBox(tui.NewFrame(chart).WithTitle(tui.Title{Text: "Default", Style: tui.FrBrightBlue | tui.Bold, Pos: tui.TitleTopCenter}),
		tui.NewFrame(chart2).WithTitle(tui.Title{Text: "With BarStyle", Style: tui.FrRed | tui.Bold, Pos: tui.TitleTopCenter}),
		tui.NewVBox(tui.NewFrame(chart3).WithTitle(tui.Title{Text: "Small", Style: tui.FrBrightGreen | tui.Bold, Pos: tui.TitleTopCenter}),
			tui.NewFrame(chart4).WithTitle(tui.Title{Text: "Colored labels", Style: tui.FrCyan, Pos: tui.TitleTopCenter}).
				WithTitle(tui.Title{Text: "Scale=15", Style: tui.FrBrightYellow, Pos: tui.TitleBottomCenter})),
		tui.NewFrame(chart5).ASCII().WithTitle(tui.Title{Text: "Wide, ASCII", Style: tui.FrBrightMagenta, Pos: tui.TitleTopCenter}),
	)).Double().WithTitle(tui.Title{Text: "TUI Compose BarChart Demo", Style: tui.FrBrightMagenta | tui.Bold, Pos: tui.TitleTopCenter}).WithPaddings(1, 3))

	wnd.Run()
}
