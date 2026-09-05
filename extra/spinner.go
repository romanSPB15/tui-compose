package extra

import (
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

const (
	SpinnerUnderscore = iota
	SpinnerBraille
	SpinnerDots
	SpinnerLine
	SpinnerBrailleReverse
)

type Spinner struct {
	typ   int
	i     int
	style cell.Style
}

func NewSpinner(typ int) *Spinner {
	return &Spinner{typ: typ}
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

func (bc *Spinner) Start(f time.Duration) *Spinner {
	go func() {
		ticker := time.NewTicker(f)
		for {
			select {
			case <-tui.CurrentWindow().OnQuit():
				ticker.Stop()
				return
			case <-ticker.C:
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
				tui.CurrentWindow().Do(func() {
					if tui.CurrentWindow().IsRunned() {
						tui.CurrentWindow().Redraw()
					}
				})
			}
		}
	}()
	return bc
}

func (bc *Spinner) WithStyle(s tui.Style) *Spinner {
	bc.style = tui.ConvertToCellStyle(s)
	return bc
}
