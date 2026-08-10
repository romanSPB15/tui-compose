package extra

import (
	"fmt"

	"github.com/romanSPB15/tui-compose/v3"
)

type Accordion struct {
	content   tui.Widget
	label     *tui.Button
	opened    bool
	CloseRune rune
	OpenRune  rune
	text      string
}

func NewAccordion(label string, content tui.Widget) *Accordion {
	acc := &Accordion{
		CloseRune: '▶',
		OpenRune:  '▼',
		text:      label,
		opened:    true,
		content:   content,
	}
	acc.label = tui.NewButton("not inited", func() {
		if acc.opened {
			acc.opened = false
			acc.label.WithText(fmt.Sprintf("%c %s", acc.CloseRune, acc.text))
		} else {
			acc.opened = true
			acc.label.WithText(fmt.Sprintf("%c %s", acc.OpenRune, acc.text))
		}
		tui.CurrentWindow().Do(func() {
			tui.CurrentWindow().Index()
			tui.CurrentWindow().Redraw()
		})
	})

	acc.label.OnClick()

	return acc
}

func (acc *Accordion) InnerText() string {
	return ""
}

func (acc *Accordion) Child() []tui.Widget {
	if acc.opened {
		return []tui.Widget{acc.label, acc.content}
	}
	return []tui.Widget{acc.label}
}

func (acc *Accordion) Pos(i int) tui.Pos {
	if i == 1 {
		return tui.Pos{Line: acc.label.MaxHeight(), Col: 0}
	}
	return tui.Pos{Line: 0, Col: 0}
}

func (acc *Accordion) MaxWidth() int {
	return max(acc.content.MaxWidth(), acc.label.MaxWidth())
}

func (acc *Accordion) MaxHeight() int {
	h := acc.label.MaxHeight()
	if acc.opened {
		h += acc.content.MaxHeight()
	}
	return h
}
