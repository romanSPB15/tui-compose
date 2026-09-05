package extra

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

type PageIndicator struct {
	total    int
	current  int
	active   rune
	inactive rune
	style    cell.Style
}

func NewPageIndicator(total int) *PageIndicator {
	return &PageIndicator{
		total:    total,
		active:   '●',
		inactive: '∙',
	}
}

func (p *PageIndicator) WithActive(r rune) *PageIndicator {
	p.active = r
	return p
}

func (p *PageIndicator) WithInactive(r rune) *PageIndicator {
	p.inactive = r
	return p
}

func (p *PageIndicator) WithStyle(s tui.Style) *PageIndicator {
	p.style = tui.ConvertToCellStyle(s)
	return p
}

func (p *PageIndicator) SetCurrent(i int) *PageIndicator {
	p.current = i
	if p.current < 0 {
		p.current = 0
	}
	if p.current >= p.total {
		p.current = p.total - 1
	}
	tui.CurrentWindow().Do(func() {
		if tui.CurrentWindow().IsRunned() {
			tui.CurrentWindow().Redraw()
		}
	})

	return p
}

func (p *PageIndicator) Width() int {
	return p.total
}

func (p *PageIndicator) Height() int {
	return 1
}

func (p *PageIndicator) Render(buf [][]cell.Cell) {
	for i := 0; i < p.total; i++ {
		var ch rune
		var style cell.Style
		if i == p.current {
			ch = p.active
			style = p.style
		} else {
			ch = p.inactive
			style = cell.Style{}
		}
		buf[0][i] = cell.Cell{Char: ch, Style: style}
	}
}
