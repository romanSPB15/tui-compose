package cell

import (
	"strings"

	"github.com/romanSPB15/tui-compose/v3/ansi" // ваш пакет ansi
)

type Cell struct {
	Char rune
	ANSI []string // коды, например ["31", "1"]
}

// Parse преобразует строку с ANSI-кодами в слайс ячеек.
func Parse(s string) []Cell {
	if s == "" {
		return nil
	}

	matches, _ := ansi.Find(s)
	if len(matches) == 0 {
		runes := []rune(s)
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

	runes := []rune(s)
	var cells []Cell
	currentStyles := []string{}

	for i := 0; i < len(runes); i++ {
		if seq, ok := styleMap[i]; ok {
			params := extractParams(seq)
			if len(params) > 0 && params[0] == "0" {
				currentStyles = []string{}
			} else {
				currentStyles = params
			}
			seqLen := len([]rune(seq))
			i += seqLen - 1
			continue
		}

		// Обычный символ
		cell := Cell{
			Char: runes[i],
			ANSI: copyStyles(currentStyles),
		}
		cells = append(cells, cell)
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

func copyStyles(styles []string) []string {
	if len(styles) == 0 {
		return nil
	}
	cpy := make([]string, len(styles))
	copy(cpy, styles)
	return cpy
}
