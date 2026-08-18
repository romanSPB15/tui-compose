package main

import (
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 Demo")

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

	pr := tui.NewTextProgress(10, '#', '-')
	pr.SetValue(0.25)

	wnd.SetContent(tui.NewHBox(tui.NewVBox(
		tui.NewFrame(chart).Rounded().WithTitle(tui.Title{Text: "Memory, %", Style: tui.FrRed}).WithTitle(tui.Title{Text: "CPU, %", Style: tui.FrYellow, Pos: tui.TitleTopRight}),
		tui.NewHBox(
			tui.NewFrame(extra.NewBlinkLabel(6).WithText("Ошибка").WithStyle(tui.FrRed).Start(time.Second/2)).Heavy(),
			tui.NewFrame(tui.NewHBox(tui.NewStaticLabel("Загрузка"), extra.NewSpinner(extra.SpinnerBrailleReverse).Start(time.Second/10))).Double(),
			tui.NewFrame(tui.NewHBox(tui.NewStaticLabel("Прогресс: 25%").WithStyle(tui.FrCyan|tui.Italic), pr, extra.NewSpinner(2).Start(time.Second/2)).WithGap(3)).Double(),
		),
	), tui.NewFrame(extra.NewTree([]extra.TreeNode{
		{
			Label: "users-db",
			Depth: 0,
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
			Style: tui.Italic | tui.BgBrightCyan | tui.FrBlack,
		},
		{
			Label: "products",
			Depth: 1,
			Style: tui.FrBlue,
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
		},
	})), tui.NewFrame(extra.NewTabs([]extra.Tab{
		{Title: "Accordion", Content: extra.NewAccordion("Open Me!", extra.NewTable([][]extra.TableCell{
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
		}).Rounded().WithHorSeparator(extra.EverywhereHorSeparator))},
		{Title: "PieChart", Content: extra.NewPieChart([]extra.PieData{
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
	})).WithTitle(tui.Title{Text: "Tabs & Accordion"}).WithTitle(tui.Title{Text: "PieChart", Pos: tui.TitleBottomCenter})))

	go func() {
		i := 0
		j := len(data) / 2
		for {
			select {
			case <-wnd.OnQuit():
				return
			case <-time.Tick(time.Second / 5):
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
