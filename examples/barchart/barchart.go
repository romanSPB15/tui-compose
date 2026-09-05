package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - BarChart")

	// данные
	data1 := []int{10, 25, 40, 55, 30, 15, 5}
	data2 := []int{14, 20, 30, 60, 35, 18, 6}
	data3 := []int{140, 200, 400, 600, 350}

	chart := extra.NewBarChart().
		WithValues(data1)

	chart2 := extra.NewBarChart().
		WithValues(data2).
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
	wnd.SetContent(tui.NewHBox(tui.NewFrame(chart).WithTitle(tui.Title{Text: "Standart"}),
		tui.NewFrame(chart2).WithTitle(tui.Title{Text: "BarStyle"}),
		tui.NewVBox(tui.NewFrame(chart3).WithTitle(tui.Title{Text: "Narrow & Low"}),
			tui.NewFrame(chart4).WithTitle(tui.Title{Text: "Colored Labels"})),
		tui.NewFrame(chart5).WithTitle(tui.Title{Text: "Wide & Custom symbols"}),
	))

	wnd.Run()
}
