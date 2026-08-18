package extra

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/cell"
)

func convertToCellStyle(s tui.Style) cell.Style {
	var cs cell.Style
	fg := int(s & 0x1F)
	if fg != 0 {
		if fg <= 8 {
			cs.Fg = strconv.Itoa(fg + 29)
		} else {
			cs.Fg = strconv.Itoa(fg + 81)
		}
	}
	bg := int((s >> 5) & 0x1F)
	if bg != 0 {
		if bg <= 8 {
			cs.Bg = strconv.Itoa(bg + 39)
		} else {
			cs.Bg = strconv.Itoa(bg + 91)
		}
	}

	if s&tui.Bold != 0 {
		cs.Args |= cell.Bold
	}
	if s&tui.Italic != 0 {
		cs.Args |= cell.Italic
	}
	if s&tui.Underline != 0 {
		cs.Args |= cell.Underline
	}
	if s&tui.Blink != 0 {
		cs.Args |= cell.Blink
	}
	if s&tui.Reverse != 0 {
		cs.Args |= cell.Reverse
	}
	if s&tui.Reset != 0 {
		cs.Args |= cell.Reset
	}
	return cs
}
