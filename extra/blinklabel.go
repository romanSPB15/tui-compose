package extra

import (
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// BlinkLabel — это виджет мигающей метки.
// Периодически скрывает и показывает текст.
// Добавлено в TUI v3.3.0.
type BlinkLabel struct {
	label   *tui.Label
	visible bool
	ticker  *time.Ticker
	stopCh  chan struct{}
	wnd     tui.Window
}

// NewBlinkLabel создаёт мигающую метку с зарезервированной шириной len.
// Добавлено в TUI v3.3.0.
func NewBlinkLabel(len int) *BlinkLabel {
	return &BlinkLabel{
		label:   tui.NewDynamicLabel("", len),
		visible: true,
		stopCh:  make(chan struct{}),
	}
}

// WithStyle применяет стиль к тексту.
// Добавлено в TUI v3.3.0.
func (b *BlinkLabel) WithStyle(s tui.Style) *BlinkLabel {
	b.label.WithStyle(s)
	return b
}

// WithText устанавливает текст метки.
// Добавлено в TUI v3.3.0.
func (b *BlinkLabel) WithText(txt string) *BlinkLabel {
	b.label.SetText(txt)
	return b
}

// Start запускает мигание с заданным интервалом.
// Добавлено в TUI v3.3.0.
func (b *BlinkLabel) Start(interval time.Duration) *BlinkLabel {
	b.ticker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-b.stopCh:
				b.ticker.Stop()
				return
			case <-b.ticker.C:
				b.visible = !b.visible
				if b.wnd != nil {
					b.wnd.Redraw()
				}
			}
		}
	}()
	return b
}

// Stop останавливает мигание.
// Добавлено в TUI v3.3.0.
func (b *BlinkLabel) Stop() {
	close(b.stopCh)
}

// Render реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (b *BlinkLabel) Render(buf [][]cell.Cell) {
	if b.visible {
		b.label.Render(buf)
	}
}

// Width реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (b *BlinkLabel) Width() int { return b.label.Width() }

// Height реализует интерфейс Widget.
// Добавлено в TUI v4.0.0.
func (b *BlinkLabel) Height() int { return b.label.Height() }

// SetText устанавливает текст метки.
// Добавлено в TUI v3.3.0.
func (b *BlinkLabel) SetText(text string) {
	b.label.Text = text
}
