package extra

import (
	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/builder"
)

type PageIndicator struct {
	total    int
	current  int
	active   rune
	inactive rune
	style    tui.Style
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
	p.style = s
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

func (p *PageIndicator) MaxWidth() int {
	return p.total
}

func (p *PageIndicator) MaxHeight() int {
	return 1
}

func (p *PageIndicator) InnerText() string {
	var builder builder.Builder
	for i := 0; i < p.total; i++ {
		var ch rune
		if i == p.current {
			ch = p.active
		} else {
			ch = p.inactive
		}
		if i == p.current {
			builder.WriteString(p.style.String())
		}
		builder.WriteRune(ch)
		if i == p.current {
			builder.WriteString("\033[0m")
		}
	}
	return builder.String()
}
