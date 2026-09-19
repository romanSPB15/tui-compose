package main

import (
	"fmt"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/tea"
)

type model struct {
	count int
}

func (m model) Init() tea.Cmd {
	return tea.After(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case time.Time:
		m.count++
		return m, tea.After(time.Second, func(t time.Time) tea.Msg {
			return t
		})
	}
	return m, nil
}

func (m model) View() tui.Widget {
	return tui.NewFrame(tui.NewVBox(
		tui.NewStaticLabel(fmt.Sprintf("Count: %d", m.count)).WithStyle(tui.FrBlue),
		tui.NewButton("Quit", func() {
			tea.Send(tea.Quit())
		}).WithStyle(tui.BgRed),
	).WithGap(1)).Heavy().WithTitle(tui.Title{Text: "TUI Compose", Pos: tui.TitleTopRight}).
		WithTitle(tui.Title{Text: "ELM", Pos: tui.TitleBottomRight}).
		WithPaddings(0, 2)
}

func main() {
	p := tea.NewProgram(model{})
	p.Run()
}
