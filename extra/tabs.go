package extra

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

type TabPosition int

const (
	TabsTop = iota
	TabsBottom
)

type tabsTopPanel struct {
	t       *Tabs
	focused bool
}

func (tp *tabsTopPanel) OnFocus() {
	tp.focused = true
}

func (tp *tabsTopPanel) OnBlur() {
	tp.focused = false
}

func (tp *tabsTopPanel) OnKeyPress(ev *input.KeyboardEvent) {
	switch ev.Key {
	case input.KeyArrowRight, input.KeyEnter, input.KeyPgDown:
		if tp.t.current < len(tp.t.tabs)-1 {
			tp.t.current++
			tui.CurrentWindow().Commit(func() {
				tui.CurrentWindow().Index()
				tui.CurrentWindow().Redraw()
			})
		}
	case input.KeyArrowLeft, input.KeyBackspace, input.KeyPgUp:
		if tp.t.current > 0 {
			tp.t.current--
			tui.CurrentWindow().Commit(func() {
				tui.CurrentWindow().Index()
				tui.CurrentWindow().Redraw()
			})
		}
	case input.KeyHome:
		if tp.t.current != 0 {
			tp.t.current = 0
			tui.CurrentWindow().Commit(func() {
				tui.CurrentWindow().Index()
				tui.CurrentWindow().Redraw()
			})
		}
	case input.KeyEnd:
		if tp.t.current != len(tp.t.tabs)-1 {
			tp.t.current = len(tp.t.tabs) - 1
			tui.CurrentWindow().Commit(func() {
				tui.CurrentWindow().Index()
				tui.CurrentWindow().Redraw()
			})
		}
	}
}

func (acc *tabsTopPanel) Render(buf [][]cell.Cell) {
	x := 0
	for i, v := range acc.t.tabs {
		s := tui.ConvertToCellStyle(v.TitleStyle)
		if acc.t.current == i {
			s = tui.ConvertToCellStyle(acc.t.selected)
		}
		runes := []rune(v.Title)
		for i, v := range runes {
			buf[0][i+x] = cell.Cell{Char: v, Style: s}
		}
		if acc.focused {
			buf[0][x+len(runes)] = cell.Cell{}
		}
		x += len(runes) + 1
	}
}

func (acc *tabsTopPanel) Width() int {
	w := 0
	for _, v := range acc.t.tabs {
		w += len(v.Title) + 1
	}
	return w - 1
}

func (acc *tabsTopPanel) Height() int {
	return 1
}

func (acc *tabsTopPanel) OnClickAt(x, _ int) {
	w := 0
	for i, v := range acc.t.tabs {
		if x > w && x < w+len(v.Title) {
			acc.t.current = i
			tui.CurrentWindow().Commit(func() {
				tui.CurrentWindow().Index()
				tui.CurrentWindow().Redraw()
			})
		}
		w += len(v.Title) + 1
	}
}

// Tab — вкладка Tabs.
type Tab struct {
	Content    tui.Widget
	Title      string
	TitleStyle tui.Style
}

// Tabs — контейнер вкладок.
type Tabs struct {
	tabs     []Tab
	current  int
	topPanel tabsTopPanel
	selected tui.Style
	tp       TabPosition
}

func NewTabs(t []Tab) *Tabs {
	tabs := &Tabs{
		tabs:     t,
		selected: tui.BgBrightWhite,
		tp:       TabsTop,
	}
	tabs.topPanel = tabsTopPanel{t: tabs}
	return tabs
}

func (acc *Tabs) Render([][]cell.Cell) {}

func (acc *Tabs) Child() []tui.Widget {
	return []tui.Widget{&acc.topPanel, acc.tabs[acc.current].Content}
}

func (acc *Tabs) Pos(i int) tui.Pos {
	if acc.tp == TabsBottom {
		if i == 1 {
			return tui.Pos{0, 0}
		}
		return tui.Pos{acc.tabs[acc.current].Content.Height(), 0}
	}
	if i == 1 {
		return tui.Pos{acc.topPanel.Height(), 0}
	}
	return tui.Pos{0, 0}
}

func (acc *Tabs) Width() int {
	return max(acc.tabs[acc.current].Content.Width(), acc.topPanel.Width())
}

func (acc *Tabs) Height() int {
	return acc.tabs[acc.current].Content.Height() + 1
}

// WithSelectedStyle устанавливает стиль для заголовка выбранной вкладки.
func (acc *Tabs) WithSelectedStyle(s tui.Style) *Tabs {
	acc.selected = s
	return acc
}

// WithCurrent выбирает текущую вкладку.
func (acc *Tabs) WithCurrent(i int) *Tabs {
	if i > 0 && i < len(acc.tabs) {
		acc.current = i
	}
	return acc
}

// WithTabPosition выбирает расположение панели с кнопками.
func (acc *Tabs) WithTabPosition(tp TabPosition) *Tabs {
	acc.tp = tp
	return acc
}
