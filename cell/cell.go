package cell

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/romanSPB15/tui-compose/v3/ansi"
)

// BIURBlRe - Bold Italic Underline Reverse Blink

type Style struct {
	Fg   string
	Bg   string
	Args uint32
}

const (
	Bold = 1 << iota
	Italic
	Underline
	Reverse
	Blink
	Reset
)

// Cell представляет одну ячейку экрана.
type Cell struct {
	Char  rune
	Style Style
}

// ANSI возвращает последовательность для перехода от предыдущего стиля к текущему.
func (c Style) ANSI(last Style) string {
	if (last == Style{}) && (c != Style{}) {
		return "\x1b[0m"
	}

	var codes []string

	if c.Args&Bold != 0 && last.Args&Bold == 0 {
		codes = append(codes, "1")
	} else if c.Args&Bold == 0 && last.Args&Bold != 0 {
		codes = append(codes, "22")
	}

	if c.Args&Italic != 0 && last.Args&Italic == 0 {
		codes = append(codes, "3")
	} else if c.Args&Italic == 0 && last.Args&Italic != 0 {
		codes = append(codes, "23")
	}

	if c.Args&Underline != 0 && last.Args&Underline == 0 {
		codes = append(codes, "4")
	} else if c.Args&Underline == 0 && last.Args&Underline != 0 {
		codes = append(codes, "24")
	}

	if c.Args&Reverse != 0 && last.Args&Reverse == 0 {
		codes = append(codes, "7")
	} else if c.Args&Reverse == 0 && last.Args&Reverse != 0 {
		codes = append(codes, "27")
	}

	if c.Args&Blink != 0 && last.Args&Blink == 0 {
		codes = append(codes, "5")
	} else if c.Args&Blink == 0 && last.Args&Blink != 0 {
		codes = append(codes, "25")
	}

	if c.Fg != last.Fg {
		if c.Fg == "" {
			codes = append(codes, "39")
		} else {
			codes = append(codes, strings.Split(c.Fg, ";")...)
		}
	}

	if c.Bg != last.Bg {
		if c.Bg == "" {
			codes = append(codes, "49")
		} else {
			codes = append(codes, strings.Split(c.Bg, ";")...)
		}
	}

	if c.Args&Reset != 0 {
		codes = []string{"0"}
	}

	if len(codes) == 0 {
		return ""
	}
	return "\033[" + strings.Join(codes, ";") + "m"
}

func (c Style) Merge(new Style) Style {
	if new.Args&Reset != 0 {
		return Style{}
	}

	c.Args |= new.Args

	if new.Fg != "" {
		c.Fg = new.Fg
	}
	if new.Bg != "" {
		c.Bg = new.Bg
	}

	return c
}

func parseANSI(seq string) Style {
	if !strings.HasPrefix(seq, "\033[") || !strings.HasSuffix(seq, "m") {
		return Style{}
	}
	params := strings.Split(strings.TrimSuffix(strings.TrimPrefix(seq, "\033["), "m"), ";")
	if len(params) == 0 {
		return Style{}
	}

	var s Style

	i := 0
	for i < len(params) {
		v, _ := strconv.Atoi(params[i])
		switch v {
		case 0:
			s.Args |= Reset
			return s
		case 1:
			s.Args |= Bold
		case 3:
			s.Args |= Italic
		case 4:
			s.Args |= Underline
		case 5:
			s.Args |= Blink
		case 7:
			s.Args |= Reverse
		case 22:
			s.Args &^= Bold
		case 23:
			s.Args &^= Italic
		case 24:
			s.Args &^= Underline
		case 25:
			s.Args &^= Blink
		case 27:
			s.Args &^= Reverse
		case 30, 31, 32, 33, 34, 35, 36, 37:
			s.Fg = fmt.Sprintf("%d", v)
		case 39:
			s.Fg = ""
		case 40, 41, 42, 43, 44, 45, 46, 47:
			s.Bg = fmt.Sprintf("%d", v)
		case 49:
			s.Bg = ""
		case 38:
			if i+1 < len(params) {
				if params[i+1] == "2" && i+3 < len(params) {
					s.Fg = fmt.Sprintf("38;2;%s;%s;%s", params[i+2], params[i+3], params[i+4])
					i += 4
				} else if params[i+1] == "5" && i+2 < len(params) {
					s.Fg = fmt.Sprintf("38;5;%s", params[i+2])
					i += 2
				}
			}
		case 48:
			if i+1 < len(params) {
				if params[i+1] == "2" && i+3 < len(params) {
					s.Bg = fmt.Sprintf("48;2;%s;%s;%s", params[i+2], params[i+3], params[i+4])
					i += 4
				} else if params[i+1] == "5" && i+2 < len(params) {
					s.Bg = fmt.Sprintf("48;5;%s", params[i+2])
					i += 2
				}
			}
		default:
			// ignore
		}
		i++
	}
	return s
}

// Parse разбирает строку с ANSI-кодами и возвращает слайс ячеек.
func Parse(s string) []Cell {
	if s == "" {
		return nil
	}

	matches, _ := ansi.Find(s)
	if len(matches) == 0 {
		runes := []rune(s)
		cells := make([]Cell, len(runes))
		for i, r := range runes {
			cells[i] = Cell{Char: r}
		}
		return cells
	}

	seqMap := make(map[int]string)
	for _, m := range matches {
		seqMap[m.Index] = m.Seq
	}

	var cells []Cell
	var currentStyle Style

	runes := []rune(s)

	i := 0
	for i < len(runes) {
		if seq, ok := seqMap[i]; ok {
			newStyle := parseANSI(seq)
			currentStyle = currentStyle.Merge(newStyle)
			i += len([]rune(seq))
			continue
		}
		cells = append(cells, Cell{Char: runes[i], Style: currentStyle})
		i++
	}

	return cells
}
