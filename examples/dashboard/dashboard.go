package main

import (
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.4 Demo")

	// данные
	data := []int{
		0, 4, 0, 8, 13, 8, 17, 21, 26, 17, 39, 17, 30, 43, 47, 52, 60, 65, 60, 69, 65, 73, 82, 91, 95, 100,
		95, 91, 82, 73, 65, 69, 60, 65, 60, 52, 47, 43, 30, 17, 39, 17, 26, 21, 17, 8, 13, 8, 0, 4, 0,
	}

	chart := extra.NewLineChart()

	setData := func(i, j int) {
		chart.WithData([]extra.Series{
			{
				Values:    data[i : i+7],
				LineStyle: tui.FrBrightRed,
			},
			{
				Values:    data[j : j+7],
				LineStyle: tui.FrBrightYellow,
			},
		})
	}

	setData(0, len(data)/2)

	chart.PointDistance = 10

	chart.WithScale(5)
	chart.WithYLabels([]int{25, 50, 75, 100})
	chart.WithXLabels([]string{"1", "2", "3", "4", "5", "6", "7"})

	chart.AxisStyle = tui.FrBrightBlack

	wnd.SetContent(tui.NewFrame(tui.NewHBox(tui.NewVBox(
		tui.NewFrame(chart).Rounded().WithTitle(tui.Title{Text: "Memory, MB", Style: tui.FrRed}).WithTitle(tui.Title{Text: "CPU, %", Style: tui.FrYellow, Pos: tui.TitleTopRight}),
		tui.NewHBox(
			tui.NewFrame(extra.NewBlinkLabel(6).WithText("Ошибка").WithStyle(tui.FrRed).Start(time.Second/2)).Heavy(),
			tui.NewFrame(tui.NewHBox(tui.NewStaticLabel("Загрузка page.html..."),
				extra.NewSpinner(extra.SpinnerBrailleReverse).WithStyle(tui.FrBrightMagenta).Start(time.Second/10),
			)).Double(),
			tui.NewFrame(tui.NewHBox(
				tui.NewGauge(20).EmptySquares().WithValue(0.25),
				extra.NewSpinner(extra.SpinnerLine).Start(time.Second/3)).WithGap(2),
			).BevelASCII(),
		),
	), tui.NewFrame(extra.NewTree([]extra.TreeNode{
		{
			Label: "users-db",
			Depth: 0,
			Style: tui.Bold | tui.FrCyan,
		},
		{
			Label: "users",
			Depth: 1,
		},
		{
			Label: "Bob",
			Depth: 2,
		},
		{
			Label: "Tom",
			Depth: 2,
			Style: tui.Italic,
		},
		{
			Label: "products",
			Depth: 1,
		},
		{
			Label: "comments",
			Depth: 1,
		},
		{
			Label: "foo",
			Depth: 1,
		},
		{
			Label: "bar",
			Depth: 2,
			Style: tui.FrBlue,
		},
	})), tui.NewFrame(extra.NewTabs([]extra.Tab{
		{Title: "  Accordion  ", TitleStyle: tui.BgYellow, Content: tui.NewVBox(extra.NewAccordion("Service List", extra.NewTable([][]extra.TableCell{
			{
				{
					Text:  "Service",
					Style: tui.Bold,
				},
				{
					Text:  "CPU",
					Style: tui.Bold,
				},
				{
					Text:  "Memory",
					Style: tui.Bold,
				},
				{
					Text:  "Status",
					Style: tui.Bold,
				},
			},
			{
				{
					Text: "grafana",
				},
				{
					Text: "0%",
				},
				{
					Text: "0MB",
				},
				{
					Text:  "Stopped",
					Style: tui.FrRed,
				},
			},
			{
				{
					Text: "postgres-db",
				},
				{
					Text: "5%",
				},
				{
					Text: "50MB",
				},
				{
					Text:  "Runned",
					Style: tui.FrBrightGreen,
				},
			},
			{
				{
					Text: "collector",
				},
				{
					Text: "10%",
				},
				{
					Text: "100MB",
				},
				{
					Text:  "Runned",
					Style: tui.FrBrightGreen,
				},
			},
			{
				{
					Text: "tui-dashboard",
				},
				{
					Text: "1%",
				},
				{
					Text: "14MB",
				},
				{
					Text:  "Runned",
					Style: tui.FrBrightGreen,
				},
			},
		}).Rounded().WithHorSeparator(extra.EverywhereHorSeparator)),
			extra.NewAccordion("TUI Compose", tui.NewHyperlink("Visit tui-compose on GitHub", "https://github.com/romanSPB15/tui-compose").WithStyle(tui.FrBlue|tui.Underline)))},
		{Title: "  PieChart  ", TitleStyle: tui.BgBrightMagenta, Content: extra.NewPieChart([]extra.PieData{
			{
				Label: "10%",
				Value: 10,
				Color: tui.FrRed,
			},
			{
				Label: "20%",
				Value: 20,
				Color: tui.FrBrightMagenta,
			},
			{
				Label: "30%",
				Value: 50,
				Color: tui.FrYellow,
			},
			{
				Label: "50%",
				Value: 50,
				Color: tui.FrBrightCyan,
			},
		}).WithRadius(15)},
	})).WithTitle(tui.Title{Text: "Tabs & Accordion"}).WithTitle(tui.Title{Text: "PieChart", Pos: tui.TitleBottomCenter}))).Rounded().
		WithTitle(tui.Title{Text: "TUI Compose Dashboard", Pos: tui.TitleTopCenter, Style: tui.FrBlue | tui.Bold}).
		WithTitle(tui.Title{Text: "v3.4", Pos: tui.TitleBottomRight, Style: tui.FrBrightBlack | tui.Italic}).
		WithBorderStyle(tui.FrBrightBlack),
	)

	go func() {
		i := 0
		j := len(data) / 2
		for {
			select {
			case <-wnd.OnQuit():
				return
			case <-time.Tick(time.Second / 200):
				wnd.Commit(func() {
					if i == len(data)-7 {
						i = 0
					}
					if j == len(data)-7 {
						j = 0
					}
					setData(i, j)
					i++
					j++
				})
			}
		}
	}()

	wnd.Run()
}
