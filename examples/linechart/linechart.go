package main

import (
	"math"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 - LineChart")

	// Данные: синусоида и косинусоида
	const n = 16
	sinVals := make([]int, n)
	cosVals := make([]int, n)
	for i := 0; i < n; i++ {
		sinVals[i] = int(50 + 50*math.Sin(float64(i)*2*math.Pi/float64(n)))
		cosVals[i] = int(50 + 50*math.Cos(float64(i)*2*math.Pi/float64(n)))
	}

	chart1 := extra.NewLineChart().
		WithData([]extra.Series{
			{Values: sinVals, LineStyle: tui.FrBrightCyan},
		}).
		WithXLabels([]string{"0", "5", "10", "15", "20"}).
		WithYLabels([]int{0, 25, 50, 75, 100}).
		WithDefaultAxis().
		WithHeight(12).
		AutoScale(). // AutoScale обязательно после WithHeight и WithData!!!
		WithPointDistance(4).
		WithDisplayPoints(true)

	chart2 := extra.NewLineChart().
		WithData([]extra.Series{
			{Values: sinVals, LineStyle: tui.FrRed},
			{Values: cosVals, LineStyle: tui.FrYellow},
		}).
		WithXLabels([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14"}).
		GenerateYLabels(4, true).
		WithRoundedAxis().
		WithHeight(12).
		WithScale(10).
		WithPointDistance(3)

	chart3 := extra.NewLineChart().
		WithData([]extra.Series{
			{Values: sinVals, LineStyle: tui.FrGreen},
		}).
		WithASCIIAxis().
		WithHeight(10).
		AutoScale().
		WithPointDistance(2).
		GenerateYLabels(3, false)

	wnd.SetContent(tui.NewFrame(tui.NewVBox(
		tui.NewHBox(
			tui.NewFrame(chart1).
				WithTitle(tui.Title{Text: "DisplayPoints", Pos: tui.TitleTopCenter, Style: tui.FrCyan}).
				WithTitle(tui.Title{Text: "PointDistance=4", Pos: tui.TitleBottomCenter, Style: tui.FrRed}),
			tui.NewFrame(chart2).
				WithTitle(tui.Title{Text: "Two series, Rounded", Pos: tui.TitleTopCenter, Style: tui.FrYellow}).
				WithTitle(tui.Title{Text: "PointDistance=3", Pos: tui.TitleBottomCenter, Style: tui.FrRed}),
		),
		tui.NewFrame(chart3).
			WithTitle(tui.Title{Text: "ASCII-style, no horizontal axis", Pos: tui.TitleTopCenter, Style: tui.FrGreen}).
			WithTitle(tui.Title{Text: "PointDistance=2", Pos: tui.TitleBottomCenter, Style: tui.FrRed}),
	)).Double().WithTitle(tui.Title{Text: "TUI Compose v3.3 - LineChart", Pos: tui.TitleTopCenter, Style: tui.Bold | tui.FrBlue}))

	wnd.Run()
}
