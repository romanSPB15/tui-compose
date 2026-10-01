package extra

import (
	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
)

// PageIndicator — это виджет индикатора страниц.
// Отображает список символов, где активный символ отличается от остальных.
// Добавлено в TUI v3.3.0.
type PageIndicator struct {
	total    int
	current  int
	active   rune
	inactive rune
	style    acell.Style
}

// NewPageIndicator создаёт индикатор с указанным количеством страниц.
// Символы по умолчанию: '●' — активная, '∙' — неактивная.
// Добавлено в TUI v3.3.0.
func NewPageIndicator(total int) *PageIndicator {
	return &PageIndicator{
		total:    total,
		active:   '●',
		inactive: '∙',
	}
}

// WithActive устанавливает символ активной страницы.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) WithActive(r rune) *PageIndicator {
	p.active = r
	return p
}

// WithInactive устанавливает символ неактивной страницы.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) WithInactive(r rune) *PageIndicator {
	p.inactive = r
	return p
}

// WithStyle устанавливает стиль активной страницы.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) WithStyle(s tui.Style) *PageIndicator {
	p.style = tui.ConvertToCellStyle(s)
	return p
}

// SetCurrent устанавливает текущую страницу по индексу.
// Значение автоматически ограничивается диапазоном [0, total-1].
// Автоматически перерисовывает окно.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) SetCurrent(i int) *PageIndicator {
	p.current = i
	if p.current < 0 {
		p.current = 0
	}
	if p.current >= p.total {
		p.current = p.total - 1
	}

	return p
}

// Width реализует интерфейс Widget.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) Width() int {
	return p.total
}

// Height реализует интерфейс Widget.
// Добавлено в TUI v3.3.0.
func (p *PageIndicator) Height() int {
	return 1
}

// Render реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (p *PageIndicator) Render(buf [][]acell.Cell) {
	for i := 0; i < p.total; i++ {
		var ch rune
		var style acell.Style
		if i == p.current {
			ch = p.active
			style = p.style
		} else {
			ch = p.inactive
			style = acell.Style{}
		}
		buf[0][i] = acell.Cell{Char: ch, Style: style}
	}
}
