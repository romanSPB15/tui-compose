package main

import (
	"fmt"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/tea"
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
	return tui.NewStaticLabel(fmt.Sprintf("Count: %d", m.count))
}

func main() {
	p := tea.NewProgram(model{})
	p.Run()
}
