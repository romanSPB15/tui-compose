package extra

import (
	"strings"
	"time"

	"github.com/romanSPB15/tui-compose/v3"
)

type BlinkLabel struct {
	label   *tui.Label
	visible bool
	ticker  *time.Ticker
	stopCh  chan struct{}
}

func NewBlinkLabel(len int) *BlinkLabel {
	return &BlinkLabel{
		label:   tui.NewDynamicLabel("", len),
		visible: true,
		stopCh:  make(chan struct{}),
	}
}

func (b *BlinkLabel) WithStyle(s tui.Style) *BlinkLabel {
	b.label.WithStyle(s)
	return b
}

func (b *BlinkLabel) WithText(txt string) *BlinkLabel {
	b.label.SetText(txt)
	return b
}

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
				tui.CurrentWindow().Do(func() {
					if tui.CurrentWindow().IsRunned() {
						tui.CurrentWindow().Redraw()
					}
				})
			}
		}
	}()
	return b
}

func (b *BlinkLabel) Stop() {
	close(b.stopCh)
}

func (b *BlinkLabel) InnerText() string {
	if b.visible {
		return b.label.InnerText()
	}
	return strings.Repeat(" ", b.MaxWidth())
}

func (b *BlinkLabel) MaxWidth() int  { return b.label.MaxWidth() }
func (b *BlinkLabel) MaxHeight() int { return b.label.MaxHeight() }

func (b *BlinkLabel) SetText(text string) {
	b.label.Text = text
}
