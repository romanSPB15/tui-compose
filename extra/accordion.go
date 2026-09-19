package extra

import (
	"fmt"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// Accordion — виджет раскрывающегося списка.
// Состоит из заголовка-кнопки и содержимого, которое можно скрыть или показать.
// Добавлено в TUI v3.3.0.
type Accordion struct {
	content   tui.Widget
	label     *tui.Button
	opened    bool
	CloseRune rune
	OpenRune  rune
	text      string
	wnd       tui.Window
}

func (btn *Accordion) Send(ev tui.Event) {
	switch ev := ev.(type) {
	case *tui.WindowEvent:
		btn.wnd = ev.Window
	}
}

// NewAccordion создаёт аккордеон с указанным заголовком и содержимым.
// По умолчанию аккордеон закрыт. Символы: '▼' — открыт, '▶' — закрыт.
// Добавлено в TUI v3.3.0.
func NewAccordion(label string, content tui.Widget) *Accordion {
	acc := &Accordion{
		CloseRune: '▶',
		OpenRune:  '▼',
		text:      label,
		opened:    true,
		content:   content,
	}
	acc.label = tui.NewButton("", func() {
		if acc.opened {
			acc.opened = false
			acc.label.WithText(fmt.Sprintf("%c %s", acc.CloseRune, acc.text))
		} else {
			acc.opened = true
			acc.label.WithText(fmt.Sprintf("%c %s", acc.OpenRune, acc.text))
		}
		if acc.wnd != nil {
			acc.wnd.Index()
			acc.wnd.Redraw()
		}
	}).WithPaddings(0, 0)

	acc.label.Send(&tui.MouseEvent{})

	return acc
}

func (acc *Accordion) Render([][]cell.Cell) {}

func (acc *Accordion) Child() []tui.Widget {
	if acc.opened {
		return []tui.Widget{acc.label, acc.content}
	}
	return []tui.Widget{acc.label}
}

func (acc *Accordion) Pos(i int) tui.Pos {
	if i == 1 {
		return tui.Pos{Line: acc.label.Height(), Col: 0}
	}
	return tui.Pos{Line: 0, Col: 0}
}

func (acc *Accordion) Width() int {
	return max(acc.content.Width(), acc.label.Width())
}

func (acc *Accordion) Height() int {
	h := acc.label.Height()
	if acc.opened {
		h += acc.content.Height()
	}
	return h
}
