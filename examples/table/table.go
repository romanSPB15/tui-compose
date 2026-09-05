package main

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("TUI Compose v3.3 Demo")

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
	})

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
	}).WithHorSeparator(extra.EverywhereHorSeparator)

	// Таблица 3 — без горизонтальных разделителей
	table3 := extra.NewTable([][]extra.TableCell{
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
	}).WithHorSeparator(extra.NoHorSeparator)

	// Оборачиваем каждую таблицу в рамку с заголовком
	frame1 := tui.NewFrame(table1).
		Rounded().
		WithTitle(tui.Title{Text: "Таблица 1 (по умолчанию)"})

	frame2 := tui.NewFrame(table2).
		Rounded().
		WithTitle(tui.Title{Text: "Таблица 2 (разделители везде)"}).
		WithTitle(tui.Title{Text: "Визуальное разделение лучше", Style: tui.FrGreen, Pos: tui.TitleBottomCenter})

	frame3 := tui.NewFrame(table3).
		Rounded().
		WithTitle(tui.Title{Text: "Таблица 3 (без разделителей)"}).
		WithTitle(tui.Title{Text: "Компактней", Style: tui.FrGreen, Pos: tui.TitleBottomCenter})

	// Размещаем рамки в горизонтальном ряду
	wnd.SetContent(tui.NewHBox(frame1, frame2, frame3))

	wnd.Run()
}
