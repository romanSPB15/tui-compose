package extra

import (
	"time"

	"github.com/romanSPB15/tui-compose/v3"
)

type Spinner struct {
	style int
	i     int
}

func NewSpinner(style int) *Spinner {
	return &Spinner{style: style}
}

func (bc *Spinner) MaxWidth() int {
	switch bc.style {
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
	switch bc.style {
	case 0:
		switch bc.i {
		case 0:
			return "_"
		case 1:
			return " "
		}
	case 1:
		switch bc.i {
		case 0:
			return "⣾"
		case 1:
			return "⣽"
		case 2:
			return "⣻"
		case 3:
			return "⢿"
		case 4:
			return "⡿"
		case 5:
			return "⣟"
		case 6:
			return "⣯"
		case 7:
			return "⣷"
		}
	case 2:
		switch bc.i {
		case 0:
			return "∙∙∙"
		case 1:
			return "●∙∙"
		case 2:
			return "∙●∙"
		case 3:
			return "∙∙●"
		}
	case 3:
		switch bc.i {
		case 0:
			return "|"
		case 1:
			return "/"
		case 2:
			return "-"
		case 3:
			return "\\"
		}
	}

	return ""
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
				switch bc.style {
				case 0:
					if bc.i > 1 {
						bc.i = 0
					}
				case 1:
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
