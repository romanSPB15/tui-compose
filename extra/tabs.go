package extra

import (
	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
)

// TabPosition — позиция панели с заголовками вкладок.
// Добавлено в TUI v3.3.0.
type TabPosition int

const (
	TabsTop = iota
	TabsBottom
)

type tabsTopPanel struct {
	t            *Tabs
	focused      bool
	hoveredIndex int
	wnd          tui.Window
}

func (tp *tabsTopPanel) Send(ev tui.Event) {
	switch e := ev.(type) {
	case *tui.WindowEvent:
		tp.wnd = e.Window
	case *tui.CheckFocusableEvent:
		e.Result = true
	case *tui.FocusEvent:
		tp.focused = e.Focused
	case *tui.MouseHoverEvent:
		if !e.Entered {
			if tp.hoveredIndex != -1 {
				tp.hoveredIndex = -1
				if tp.wnd != nil {
					tp.wnd.Redraw()
				}
			}
		}
	case *acell.KeyboardEvent:
		switch e.Key {
		case acell.KeyArrowRight, acell.KeyEnter, acell.KeyPgDown:
			if tp.t.current < len(tp.t.tabs)-1 {
				tp.t.current++
				if tp.wnd != nil {
					tp.wnd.Index()
					tp.wnd.Redraw()
				}
			}
		case acell.KeyArrowLeft, acell.KeyBackspace, acell.KeyPgUp:
			if tp.t.current > 0 {
				tp.t.current--
				if tp.wnd != nil {
					tp.wnd.Index()
					tp.wnd.Redraw()
				}
			}
		case acell.KeyHome:
			if tp.t.current != 0 {
				tp.t.current = 0
				if tp.wnd != nil {
					tp.wnd.Index()
					tp.wnd.Redraw()
				}
			}
		case acell.KeyEnd:
			if tp.t.current != len(tp.t.tabs)-1 {
				tp.t.current = len(tp.t.tabs) - 1
				if tp.wnd != nil {
					tp.wnd.Index()
					tp.wnd.Redraw()
				}
			}
		}
	case *acell.MouseEvent:
		switch e.Action {
		case acell.MouseMove:
			idx := tp.tabAt(e.Pos.X)
			if idx != tp.hoveredIndex {
				tp.hoveredIndex = idx
				if tp.wnd != nil {
					tp.wnd.Redraw()
				}
			}
		case acell.MousePress:
			idx := tp.tabAt(e.Pos.X)
			if idx >= 0 {
				tp.t.current = idx
				if tp.wnd != nil {
					tp.wnd.Index()
					tp.wnd.Redraw()
				}
			}
		}
	}
}

func (tp *tabsTopPanel) tabAt(x int) int {
	w := 0
	for i, v := range tp.t.tabs {
		titleLen := len([]rune(v.Title))
		if x >= w && x < w+titleLen {
			return i
		}
		w += titleLen + 1
	}
	return -1
}

func (tp *tabsTopPanel) Render(buf [][]acell.Cell) {
	x := 0
	for i, v := range tp.t.tabs {
		s := tui.ConvertToCellStyle(v.TitleStyle)
		switch {
		case tp.t.current == i:
			s = tui.ConvertToCellStyle(tp.t.selected)
		case i == tp.hoveredIndex && tp.t.hover != 0:
			s = tui.ConvertToCellStyle(tp.t.hover)
		}
		runes := []rune(v.Title)
		for j, r := range runes {
			buf[0][j+x] = acell.Cell{Char: r, Style: s}
		}
		if tp.focused && i != len(tp.t.tabs)-1 {
			buf[0][x+len(runes)] = acell.Cell{Char: '_'}
		}
		x += len(runes) + 1
	}
}

func (tp *tabsTopPanel) Width() int {
	w := 0
	for _, v := range tp.t.tabs {
		w += len([]rune(v.Title)) + 1
	}
	return w - 1
}

func (tp *tabsTopPanel) Height() int {
	return 1
}

// Tab — вкладка Tabs.
// Title — заголовок вкладки, отображаемый на панели.
// TitleStyle — стиль заголовка.
// Content — виджет содержимого.
// Добавлено в TUI v3.3.0.
type Tab struct {
	Content    tui.Widget
	Title      string
	TitleStyle tui.Style
}

// Tabs — контейнер вкладок.
// Отображает панель заголовков и содержимое активной вкладки.
// Добавлено в TUI v3.3.0.
type Tabs struct {
	tabs     []Tab
	current  int
	topPanel tabsTopPanel
	selected tui.Style
	hover    tui.Style
	tp       TabPosition
}

// NewTabs создаёт контейнер вкладок из указанного списка.
// По умолчанию активна первая вкладка, панель заголовков сверху.
// Добавлено в TUI v3.3.0.
func NewTabs(t []Tab) *Tabs {
	tabs := &Tabs{
		tabs:     t,
		selected: tui.BgBrightWhite,
		hover:    tui.BgBrightBlack,
		tp:       TabsTop,
	}
	tabs.topPanel = tabsTopPanel{t: tabs, hoveredIndex: -1}
	return tabs
}

func (acc *Tabs) Render([][]acell.Cell) {}

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

// WithSelectedStyle устанавливает стиль заголовка выбранной вкладки.
// Добавлено в TUI v3.3.0.
func (acc *Tabs) WithSelectedStyle(s tui.Style) *Tabs {
	acc.selected = s
	return acc
}

// WithHoverStyle устанавливает стиль заголовка при наведении.
func (acc *Tabs) WithHoverStyle(s tui.Style) *Tabs {
	acc.hover = s
	return acc
}

// WithCurrent выбирает активную вкладку по индексу.
// Значения вне допустимого диапазона игнорируются.
// Добавлено в TUI v3.3.0.
func (acc *Tabs) WithCurrent(i int) *Tabs {
	if i > 0 && i < len(acc.tabs) {
		acc.current = i
	}
	return acc
}

// WithTabPosition задаёт позицию панели заголовков: сверху или снизу.
// Добавлено в TUI v3.3.0.
func (acc *Tabs) WithTabPosition(tp TabPosition) *Tabs {
	acc.tp = tp
	return acc
}

var _ tui.EventHandler = (*tabsTopPanel)(nil)
