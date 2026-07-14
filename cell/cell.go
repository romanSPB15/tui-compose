package cell

import (
	"strconv"
	"strings"

	"github.com/romanSPB15/tui-compose/v3/ansi"
)

// Cell представляет одну ячейку экрана.
type Cell struct {
	Char rune
	ANSI []string
}

type styleState struct {
	fg        string
	bg        string
	bold      bool
	italic    bool
	underline bool
	reverse   bool
	strike    bool
}

func (s *styleState) reset() {
	s.fg = ""
	s.bg = ""
	s.bold = false
	s.italic = false
	s.underline = false
	s.reverse = false
	s.strike = false
}

func (s *styleState) applyParams(params []string) {
	if len(params) == 0 {
		return
	}
	if params[0] == "0" {
		s.reset()
		return
	}

	for i := 0; i < len(params); i++ {
		code, _ := strconv.Atoi(params[i])
		switch {
		case code == 0:
			s.reset()
		case code == 1:
			s.bold = true
		case code == 3:
			s.italic = true
		case code == 4:
			s.underline = true
		case code == 7:
			s.reverse = true
		case code == 9:
			s.strike = true
		case code == 21:
			s.underline = true
		case code == 22:
			s.bold = false
		case code == 23:
			s.italic = false
		case code == 24:
			s.underline = false
		case code == 27:
			s.reverse = false
		case code == 29:
			s.strike = false
		case code >= 30 && code <= 37:
			s.fg = strconv.Itoa(code)
		case code >= 40 && code <= 47:
			s.bg = strconv.Itoa(code)
		case code >= 90 && code <= 97:
			s.fg = strconv.Itoa(code)
		case code >= 100 && code <= 107:
			s.bg = strconv.Itoa(code)
		case code == 38:
			if i+1 < len(params) {
				if params[i+1] == "2" && i+4 < len(params) {
					r, _ := strconv.Atoi(params[i+2])
					g, _ := strconv.Atoi(params[i+3])
					b, _ := strconv.Atoi(params[i+4])
					s.fg = "38;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
					i += 4
				} else if params[i+1] == "5" && i+2 < len(params) {
					idx, _ := strconv.Atoi(params[i+2])
					s.fg = "38;5;" + strconv.Itoa(idx)
					i += 2
				}
			}
		case code == 39:
			s.fg = ""
		case code == 48:
			if i+1 < len(params) {
				if params[i+1] == "2" && i+4 < len(params) {
					r, _ := strconv.Atoi(params[i+2])
					g, _ := strconv.Atoi(params[i+3])
					b, _ := strconv.Atoi(params[i+4])
					s.bg = "48;2;" + strconv.Itoa(r) + ";" + strconv.Itoa(g) + ";" + strconv.Itoa(b)
					i += 4
				} else if params[i+1] == "5" && i+2 < len(params) {
					idx, _ := strconv.Atoi(params[i+2])
					s.bg = "48;5;" + strconv.Itoa(idx)
					i += 2
				}
			}
		case code == 49:
			s.bg = ""
		}
	}
}

func (s *styleState) toSlice() []string {
	var parts []string
	if s.fg != "" {
		parts = append(parts, s.fg)
	}
	if s.bg != "" {
		parts = append(parts, s.bg)
	}
	if s.bold {
		parts = append(parts, "1")
	}
	if s.italic {
		parts = append(parts, "3")
	}
	if s.underline {
		parts = append(parts, "4")
	}
	if s.reverse {
		parts = append(parts, "7")
	}
	if s.strike {
		parts = append(parts, "9")
	}
	return parts
}

// Parse разбирает строку с ANSI-кодами и возвращает слайс ячеек.
func Parse(s string) []Cell {
	if s == "" {
		return nil
	}

	matches, clean := ansi.Find(s)
	if len(matches) == 0 {
		runes := []rune(clean)
		cells := make([]Cell, len(runes))
		for i, r := range runes {
			cells[i] = Cell{Char: r, ANSI: nil}
		}
		return cells
	}

	styleMap := make(map[int]string)
	for _, m := range matches {
		styleMap[m.Index] = m.Seq
	}

	runes := []rune(clean)
	var cells []Cell
	state := &styleState{}
	state.reset()

	for i := 0; i < len(runes); i++ {
		if seq, ok := styleMap[i]; ok {
			params := extractParams(seq)
			if len(params) > 0 {
				state.applyParams(params)
			}
			continue
		}

		c := Cell{
			Char: runes[i],
			ANSI: state.toSlice(),
		}
		cells = append(cells, c)
	}

	return cells
}

func extractParams(seq string) []string {
	if !strings.HasPrefix(seq, "\x1b[") {
		return nil
	}
	lastIdx := strings.LastIndexFunc(seq, func(r rune) bool {
		return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
	})
	if lastIdx == -1 {
		return nil
	}
	if seq[lastIdx] != 'm' {
		return nil
	}
	params := seq[len("\x1b["):lastIdx]
	if params == "" {
		return nil
	}
	return strings.Split(params, ";")
}
