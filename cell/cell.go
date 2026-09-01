package cell

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v3/ansi"
	"github.com/romanSPB15/tui-compose/v3/builder"
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

	resetFg
	resetBg
)

// Cell представляет одну ячейку экрана.
type Cell struct {
	Char  rune
	Style Style
}

// ANSI возвращает последовательность для перехода от предыдущего стиля к текущему.
func (c Style) ANSI(last Style) string {
	if (c == Style{}) && (last != Style{}) {
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

func parseANSI(seq string) (Style, uint16) {
	if !strings.HasPrefix(seq, "\033[") || !strings.HasSuffix(seq, "m") {
		return Style{}, 0
	}
	params := strings.Split(strings.TrimSuffix(strings.TrimPrefix(seq, "\033["), "m"), ";")
	if len(params) == 0 {
		return Style{}, 0
	}

	var s Style
	var clearMask uint16

	i := 0
	for i < len(params) {
		v, _ := strconv.Atoi(params[i])
		switch v {
		case 0:
			s.Args |= Reset
			return s, 0
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
			clearMask |= Bold
		case 23:
			clearMask |= Italic
		case 24:
			clearMask |= Underline
		case 25:
			clearMask |= Blink
		case 27:
			clearMask |= Reverse
		case 30, 31, 32, 33, 34, 35, 36, 37:
			s.Fg = fmt.Sprintf("%d", v)
		case 90, 91, 92, 93, 94, 95, 96, 97:
			s.Fg = fmt.Sprintf("%d", v)
		case 40, 41, 42, 43, 44, 45, 46, 47:
			s.Bg = fmt.Sprintf("%d", v)
		case 100, 101, 102, 103, 104, 105, 106, 107:
			s.Bg = fmt.Sprintf("%d", v)
		case 39:
			clearMask |= resetFg
		case 49:
			clearMask |= resetBg
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
	return s, clearMask
}

// Parse разбирает строку с ANSI-кодами и возвращает слайс ячеек.
// zero-allocation
func Parse(s string, buf *[]Cell) []Cell {
	cells, _ := ParseFromTo(s, buf, Style{})
	return cells
}

func ParseFromTo(s string, buf *[]Cell, currentStyle Style) ([]Cell, Style) {
	if s == "" {
		return nil, currentStyle
	}

	matches, _ := ansi.Find(s)

	cells := *buf
	cells = cells[:0]

	if len(matches) == 0 {
		for i := 0; i < len(s); {
			r, size := utf8.DecodeRuneInString(s[i:])
			cells = append(cells, Cell{Char: r, Style: currentStyle})
			i += size
		}
		return cells, currentStyle
	}

	i := 0
	mi := 0

	for i < len(s) {
		if mi < len(matches) && matches[mi].Index == i {
			seq := matches[mi].Seq
			newStyle, clearMask := parseANSI(seq)

			currentStyle.Args &^= uint32(clearMask)

			if clearMask&resetFg != 0 {
				currentStyle.Fg = ""
			}

			if clearMask&resetBg != 0 {
				currentStyle.Bg = ""
			}

			currentStyle = currentStyle.Merge(newStyle)
			i += len(seq)
			mi++
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		cells = append(cells, Cell{Char: r, Style: currentStyle})
		i += size
	}

	return cells, currentStyle
}

func ParseMultiline(s string) [][]Cell {
	if s == "" {
		return nil
	}

	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	result := make([][]Cell, len(lines))

	current := Style{}

	for i, line := range lines {
		line = strings.TrimSuffix(line, "\r")

		var buf []Cell
		cells, s := ParseFromTo(line, &buf, current)
		current = s

		result[i] = append([]Cell(nil), cells...) // копируем
	}

	maxW := 0
	for _, row := range result {
		if len(row) > maxW {
			maxW = len(row)
		}
	}
	for i, row := range result {
		if len(row) < maxW {
			pad := make([]Cell, maxW-len(row))
			for j := range pad {
				pad[j] = Cell{Char: ' ', Style: Style{}}
			}
			result[i] = append(row, pad...)
		}
	}
	return result
}

func ToString(cells [][]Cell) string {
	if len(cells) == 0 {
		return ""
	}
	var builder builder.Builder
	for rowIdx, row := range cells {
		if rowIdx > 0 {
			builder.WriteString("\n")
		}
		var last Style
		for _, cell := range row {
			if ansi := cell.Style.ANSI(last); ansi != "" {
				builder.WriteString(ansi)
			}
			builder.WriteRune(cell.Char)
			last = cell.Style
		}

		if last != (Style{}) {
			builder.WriteString("\033[0m")
		}
	}
	return builder.String()
}
