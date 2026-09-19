package extra

import (
	"sync"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// Типы спиннеров.
// Передаются в Spinner при создании.
const (
	SpinnerUnderscore = iota
	SpinnerBraille
	SpinnerDots
	SpinnerLine
	SpinnerBrailleReverse
)

// Spinner — виджет спиннера.
// Добавлено в TUI v3.3.0.
//
// s := extra.NewSpinner(extra.SpinnerBrailleReverse)
// s.Start(time.Second/15)
type Spinner struct {
	typ       int
	i         int
	style     cell.Style
	wnd       tui.Window
	wndReady  chan struct{}
	wndOnce   sync.Once
	startOnce sync.Once
}

func (sp *Spinner) Send(ev tui.Event) {
	switch ev := ev.(type) {
	case *tui.WindowEvent:
		sp.wnd = ev.Window
		sp.wndOnce.Do(func() { close(sp.wndReady) })
	}
}

// NewSpinner создаёт спиннер с указанным типом.
func NewSpinner(typ int) *Spinner {
	return &Spinner{
		typ:      typ,
		wndReady: make(chan struct{}),
	}
}

func (bc *Spinner) Width() int {
	switch bc.typ {
	case 2:
		return 3
	default:
		return 1
	}
}

func (bc *Spinner) Height() int {
	return 1
}

func (bc *Spinner) Render(buf [][]cell.Cell) {
	switch bc.typ {
	case 0:
		switch bc.i {
		case 0:
			buf[0][0] = cell.Cell{Char: '_', Style: bc.style}
		case 1:
			buf[0][0] = cell.Cell{Char: ' ', Style: bc.style}
		}
	case 1:
		switch bc.i {
		case 0:
			buf[0][0] = cell.Cell{Char: '⣾', Style: bc.style}
		case 1:
			buf[0][0] = cell.Cell{Char: '⣽', Style: bc.style}
		case 2:
			buf[0][0] = cell.Cell{Char: '⣻', Style: bc.style}
		case 3:
			buf[0][0] = cell.Cell{Char: '⢿', Style: bc.style}
		case 4:
			buf[0][0] = cell.Cell{Char: '⡿', Style: bc.style}
		case 5:
			buf[0][0] = cell.Cell{Char: '⣟', Style: bc.style}
		case 6:
			buf[0][0] = cell.Cell{Char: '⣯', Style: bc.style}
		case 7:
			buf[0][0] = cell.Cell{Char: '⣷', Style: bc.style}
		}
	case 2:
		switch bc.i {
		case 0:
			buf[0][0] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][1] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][2] = cell.Cell{Char: '∙', Style: bc.style}
		case 1:
			buf[0][0] = cell.Cell{Char: '●', Style: bc.style}
			buf[0][1] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][2] = cell.Cell{Char: '∙', Style: bc.style}
		case 2:
			buf[0][0] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][1] = cell.Cell{Char: '●', Style: bc.style}
			buf[0][2] = cell.Cell{Char: '∙', Style: bc.style}
		case 3:
			buf[0][0] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][1] = cell.Cell{Char: '∙', Style: bc.style}
			buf[0][2] = cell.Cell{Char: '●', Style: bc.style}
		}
	case 3:
		switch bc.i {
		case 0:
			buf[0][0] = cell.Cell{Char: '|', Style: bc.style}
		case 1:
			buf[0][0] = cell.Cell{Char: '/', Style: bc.style}
		case 2:
			buf[0][0] = cell.Cell{Char: '-', Style: bc.style}
		case 3:
			buf[0][0] = cell.Cell{Char: '\\', Style: bc.style}
		}
	case 4:
		switch bc.i {
		case 7:
			buf[0][0] = cell.Cell{Char: '⣾', Style: bc.style}
		case 6:
			buf[0][0] = cell.Cell{Char: '⣽', Style: bc.style}
		case 5:
			buf[0][0] = cell.Cell{Char: '⣻', Style: bc.style}
		case 4:
			buf[0][0] = cell.Cell{Char: '⢿', Style: bc.style}
		case 3:
			buf[0][0] = cell.Cell{Char: '⡿', Style: bc.style}
		case 2:
			buf[0][0] = cell.Cell{Char: '⣟', Style: bc.style}
		case 1:
			buf[0][0] = cell.Cell{Char: '⣯', Style: bc.style}
		case 0:
			buf[0][0] = cell.Cell{Char: '⣷', Style: bc.style}
		}
	}
}

// Start запускает анимацию.
// Анимация автоматически останвливается при выходе из приложения.
func (bc *Spinner) Start(f time.Duration) *Spinner {
	bc.startOnce.Do(func() {
		go func() {
			<-bc.wndReady
			ticker := time.NewTicker(f)
			defer ticker.Stop()
			for {
				select {
				case <-bc.wnd.OnQuit():
					return
				case <-ticker.C:
					bc.wnd.Commit(func() {
						bc.i++
						switch bc.typ {
						case 0:
							if bc.i > 1 {
								bc.i = 0
							}
						case 1, 4:
							if bc.i > 7 {
								bc.i = 0
							}
						case 2, 3:
							if bc.i > 3 {
								bc.i = 0
							}
						}
					})
				}
			}
		}()
	})
	return bc
}

// WithStyle устанавливает стиль для спиннера.
func (bc *Spinner) WithStyle(s tui.Style) *Spinner {
	bc.style = tui.ConvertToCellStyle(s)
	return bc
}
