package extra

import (
	"time"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/builder"
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
	style tui.Style
}

func NewSpinner(typ int) *Spinner {
	return &Spinner{typ: typ}
}

func (bc *Spinner) MaxWidth() int {
	switch bc.typ {
	case 2:
		return 3
	default:
		return 1
	}
}

func (bc *Spinner) MaxHeight() int {
	return 1
}

func (bc *Spinner) InnerText() string {
	b := builder.New(10)

	if bc.style != 0 {
		b.WriteString(bc.style.String())
	}
	switch bc.typ {
	case 0:
		switch bc.i {
		case 0:
			b.WriteByte('_')
		case 1:
			b.WriteByte(' ')
		}
	case 1:
		switch bc.i {
		case 0:
			b.WriteRune('⣾')
		case 1:
			b.WriteRune('⣽')
		case 2:
			b.WriteRune('⣻')
		case 3:
			b.WriteRune('⢿')
		case 4:
			b.WriteRune('⡿')
		case 5:
			b.WriteRune('⣟')
		case 6:
			b.WriteRune('⣯')
		case 7:
			b.WriteRune('⣷')
		}
	case 2:
		switch bc.i {
		case 0:
			b.WriteString("∙∙∙")
		case 1:
			b.WriteString("●∙∙")
		case 2:
			b.WriteString("∙●∙")
		case 3:
			b.WriteString("∙∙●")
		}
	case 3:
		switch bc.i {
		case 0:
			b.WriteByte('|')
		case 1:
			b.WriteByte('/')
		case 2:
			b.WriteByte('-')
		case 3:
			b.WriteByte('\\')
		}
	case 4:
		switch bc.i {
		case 7:
			b.WriteRune('⣾')
		case 6:
			b.WriteRune('⣽')
		case 5:
			b.WriteRune('⣻')
		case 4:
			b.WriteRune('⢿')
		case 3:
			b.WriteRune('⡿')
		case 2:
			b.WriteRune('⣟')
		case 1:
			b.WriteRune('⣯')
		case 0:
			b.WriteRune('⣷')
		}
	}

	if bc.style != 0 {
		b.WriteString(tui.Reset.String())
	}

	return b.String()
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
	bc.style = s
	return bc
}
