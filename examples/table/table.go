package main

import (
	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose Table Demo")

	// Таблица 1 — без разделителей (по умолчанию)
	table1 := extra.NewTable([][]extra.TableCell{
		{
			{Text: "Service", Style: tui.Bold},
			{Text: "CPU", Style: tui.Bold},
			{Text: "Memory", Style: tui.Bold},
			{Text: "Status", Style: tui.Bold},
		},
		{
			{Text: "grafana"},
			{Text: "0%"},
			{Text: "0MB"},
			{Text: "Stopped", Style: tui.FrRed},
		},
		{
			{Text: "postgres-db"},
			{Text: "5%"},
			{Text: "50MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "collector"},
			{Text: "10%"},
			{Text: "100MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "tui-dashboard"},
			{Text: "1%"},
			{Text: "14MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
	}).WithBorder(extra.BorderAll)

	// Таблица 2 — с разделителями между всеми строками
	table2 := extra.NewTable([][]extra.TableCell{
		{
			{Text: "Service", Style: tui.Bold},
			{Text: "CPU", Style: tui.Bold},
			{Text: "Memory", Style: tui.Bold},
			{Text: "Status", Style: tui.Bold},
		},
		{
			{Text: "grafana"},
			{Text: "0%"},
			{Text: "0MB"},
			{Text: "Stopped", Style: tui.FrRed},
		},
		{
			{Text: "postgres-db"},
			{Text: "5%"},
			{Text: "50MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "collector"},
			{Text: "10%"},
			{Text: "100MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "tui-dashboard"},
			{Text: "1%"},
			{Text: "14MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
	}).WithBorder(extra.BorderInner).WithRowStyle(func(i int) acell.Style {
		if i%2 == 0 {
			return acell.Style{}
		}
		return acell.Style{Bg: acell.BgGrey4}
	})

	// Таблица 3 — без горизонтальных разделителей
	table3 := extra.NewTable([][]extra.TableCell{
		{
			{Text: "Service"},
			{Text: "CPU"},
			{Text: "Memory"},
			{Text: "Status"},
		},
		{
			{Text: "grafana"},
			{Text: "0%"},
			{Text: "0MB"},
			{Text: "Stopped", Style: tui.FrRed},
		},
		{
			{Text: "postgres-db"},
			{Text: "5%"},
			{Text: "50MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "collector"},
			{Text: "10%"},
			{Text: "100MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
		{
			{Text: "tui-dashboard"},
			{Text: "1%"},
			{Text: "14MB"},
			{Text: "Runned", Style: tui.FrBrightGreen},
		},
	}).WithBorder(extra.BorderNone).WithRowStyle(func(i int) acell.Style {
		if i == 0 {
			return acell.Style{Args: acell.Bold, Fg: acell.FgYellow}
		}
		if i%2 == 0 {
			return acell.Style{}
		}
		return acell.Style{Bg: acell.BgGrey4}
	})

	// Оборачиваем каждую таблицу в рамку с заголовком
	frame1 := tui.NewFrame(table1).
		Rounded()

	frame2 := tui.NewFrame(table2).
		Rounded()

	frame3 := tui.NewFrame(table3).
		Rounded()

	// Размещаем рамки в горизонтальном ряду
	wnd.SetContent(tui.NewFlex(
		tui.FlexItem(frame1, 1),
		tui.FlexItem(frame2, 1),
		tui.FlexItem(frame3, 1),
	).Horizontal().WithGap(1))

	wnd.Run()
}
