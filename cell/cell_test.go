package cell_test

import (
	"slices"
	"testing"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

func cells(chars string, styles ...cell.Style) []cell.Cell {
	runes := []rune(chars)

	res := make([]cell.Cell, len(runes))

	currentStyle := styles[0]

	sIdx := 1

	for i, ch := range runes {
		res[i] = cell.Cell{Char: ch, Style: currentStyle}
		if sIdx < len(styles) {
			currentStyle = styles[sIdx]
			sIdx++
		}
	}
	return res
}

func TestParse(t *testing.T) {
	tt := []struct {
		Input    string
		Expected []cell.Cell
	}{
		{Input: "\033[31mHello\033[0m World",
			Expected: append(
				cells("Hello", cell.Style{Fg: "31"}),
				cells(" World", cell.Style{})...,
			),
		},
		{
			Input: "\033[1;32mBold green\033[0mNormal",
			Expected: append(
				cells("Bold green", cell.Style{Fg: "32", Args: cell.Bold}),
				cells("Normal", cell.Style{})...,
			),
		},
		{
			Input: "\033[40;30m▀\033[0m▀",
			Expected: append(
				cells("▀", cell.Style{Fg: "30", Bg: "40"}),
				cells("▀", cell.Style{})...,
			),
		},
		{
			Input: "\033[40;30m▀\033[0m▀\033[43m",
			Expected: append(
				cells("▀", cell.Style{Fg: "30", Bg: "40"}),
				cells("▀", cell.Style{})...,
			),
		},
		{
			Input:    "\033[30m▀\033[40m▀\033[35m",
			Expected: cells("▀▀", cell.Style{Fg: "30"}, cell.Style{Fg: "30", Bg: "40"}),
		},
	}
	for i, test := range tt {
		got := cell.Parse(test.Input)
		if !slices.Equal(got, test.Expected) {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}

func TestStyleANSI(t *testing.T) {
	tt := []struct {
		name     string
		last     cell.Style
		new      cell.Style
		expected string
	}{
		{
			name:     "both empty",
			last:     cell.Style{},
			new:      cell.Style{},
			expected: "",
		},
		{
			name:     "from empty to fg color",
			last:     cell.Style{},
			new:      cell.Style{Fg: "31"},
			expected: "\x1b[31m",
		},
		{
			name:     "from empty to bg color",
			last:     cell.Style{},
			new:      cell.Style{Bg: "44"},
			expected: "\x1b[44m",
		},
		{
			name:     "from empty to bold",
			last:     cell.Style{},
			new:      cell.Style{Args: cell.Bold},
			expected: "\x1b[1m",
		},
		{
			name:     "from fg to empty -> reset",
			last:     cell.Style{Fg: "31"},
			new:      cell.Style{},
			expected: "\x1b[0m",
		},
		{
			name:     "from bg to empty -> reset",
			last:     cell.Style{Bg: "44"},
			new:      cell.Style{},
			expected: "\x1b[0m",
		},
		{
			name:     "from bold to empty -> reset",
			last:     cell.Style{Args: cell.Bold},
			new:      cell.Style{},
			expected: "\x1b[0m",
		},
		{
			name:     "change fg color",
			last:     cell.Style{Fg: "31"},
			new:      cell.Style{Fg: "32"},
			expected: "\x1b[32m",
		},
		{
			name:     "change bg color",
			last:     cell.Style{Bg: "44"},
			new:      cell.Style{Bg: "45"},
			expected: "\x1b[45m",
		},
		{
			name:     "turn off bold, turn on italic",
			last:     cell.Style{Args: cell.Bold},
			new:      cell.Style{Args: cell.Italic},
			expected: "\x1b[22;3m",
		},
		{
			name:     "turn off italic only",
			last:     cell.Style{Args: cell.Bold | cell.Italic},
			new:      cell.Style{Args: cell.Bold},
			expected: "\x1b[23m",
		},
		{
			name:     "explicit reset",
			last:     cell.Style{Fg: "31", Args: cell.Bold},
			new:      cell.Style{Args: cell.Reset},
			expected: "\x1b[0m",
		},
		{
			name:     "from empty to reset -> reset",
			last:     cell.Style{},
			new:      cell.Style{Args: cell.Reset},
			expected: "\x1b[0m",
		},
		{
			name:     "multiple changes: fg and bold to empty",
			last:     cell.Style{Fg: "31", Args: cell.Bold | cell.Underline},
			new:      cell.Style{},
			expected: "\x1b[0m",
		},
		{
			name:     "change fg and bold -> only fg change (bold stays)",
			last:     cell.Style{Fg: "31", Args: cell.Bold},
			new:      cell.Style{Fg: "32", Args: cell.Bold},
			expected: "\x1b[32m",
		},
		{
			name:     "turn off underline, keep fg",
			last:     cell.Style{Fg: "31", Args: cell.Underline},
			new:      cell.Style{Fg: "31"},
			expected: "\x1b[24m",
		},
		{
			name:     "complex: last with italic and bg, new with fg and bold",
			last:     cell.Style{Bg: "44", Args: cell.Italic},
			new:      cell.Style{Fg: "31", Args: cell.Bold},
			expected: "\x1b[1;23;31;49m",
		},
		{
			name:     "from non-empty to empty",
			last:     cell.Style{Fg: "31", Args: cell.Bold | cell.Underline},
			new:      cell.Style{},
			expected: "\x1b[0m",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.new.ANSI(tc.last)
			if got != tc.expected {
				t.Errorf("ANSI(%+v) = %q, want %q", tc.last, got, tc.expected)
			}
		})
	}
}
