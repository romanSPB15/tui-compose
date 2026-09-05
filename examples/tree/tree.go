package main

import (
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/extra"
)

func main() {
	wnd := tui.NewWindow()

	nodes := []extra.TreeNode{
		// Корень
		{Label: "TUI Compose", Depth: 0, Style: tui.Bold | tui.FrBlue},

		{Label: "tui/", Depth: 1, Style: tui.Bold | tui.FrRed},
		{Label: "Label", Depth: 2},
		{Label: "Button", Depth: 2},
		{Label: "Checkbox", Depth: 2},
		{Label: "InputField", Depth: 2},
		{Label: "Frame", Depth: 2},
		{Label: "ColorProgress", Depth: 2},
		{Label: "TextProgress", Depth: 2},
		{Label: "Canvas", Depth: 2},
		{Label: "CanvasRGB", Depth: 2},

		{Label: "extra/", Depth: 1, Style: tui.Bold | tui.FrRed},
		{Label: "Accordion", Depth: 2},
		{Label: "BarChart", Depth: 2},
		{Label: "BlinkLabel", Depth: 2},
		{Label: "LineChart", Depth: 2},
		{Label: "PageIndicator", Depth: 2},
		{Label: "PieChart", Depth: 2},
		{Label: "Sparkline", Depth: 2},
		{Label: "Spinner", Depth: 2},
		{Label: "Table", Depth: 2},
		{Label: "Tabs", Depth: 2},
		{Label: "TextView", Depth: 2},
		{Label: "Tree", Depth: 2},
	}

	last := &nodes[len(nodes)-1] // Tree

	go func() {
		tick := time.NewTicker(time.Second / 2)
		defer tick.Stop()
		for {
			select {
			case <-wnd.OnQuit():
				return
			case <-tick.C:
				wnd.Commit(func() {
					if last.Style == 0 {
						last.Style = tui.FrBlue | tui.BgBrightYellow
					} else {
						last.Style = 0
					}
				})
			}
		}
	}()

	// все четыре стиля дерева
	wnd.SetContent(tui.NewFrame(tui.NewHBox(
		extra.NewTree(nodes),
		extra.NewTree(nodes).Rounded(),
		extra.NewTree(nodes).ASCII(),
		extra.NewTree(nodes).Heavy(),
	)))

	wnd.Run()
}
