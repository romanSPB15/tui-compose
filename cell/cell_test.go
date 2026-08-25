package cell_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

func cells(chars string, styles ...cell.Style) []cell.Cell {
	runes := []rune(chars)

	res := make([]cell.Cell, len(runes))

	var currentStyle cell.Style
	if len(styles) > 0 {
		currentStyle = styles[0]
	}

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
		{
			Input:    "\033[30m⣾\033[0m⢿",
			Expected: cells("⣾⢿", cell.Style{Fg: "30"}, cell.Style{}),
		},
		{
			Input:    "\033[30m⣾\033[0m⢿",
			Expected: cells("⣾⢿", cell.Style{Fg: "30"}, cell.Style{}),
		},
		{
			Input: "\033[2mH\033[38;2;200;100;50m\033[48;2;120;255;80mi",
			Expected: cells("Hi", cell.Style{}, cell.Style{
				Fg: "38;2;200;100;50",
				Bg: "48;2;120;255;80",
			}),
		},
		{
			Input: "\033[101;5;4;7m\033[95mH\033[39mi",
			Expected: cells("Hi", cell.Style{
				Fg:   "95",
				Bg:   "101",
				Args: cell.Blink | cell.Underline | cell.Reverse,
			}, cell.Style{
				Bg:   "101",
				Args: cell.Blink | cell.Underline | cell.Reverse,
			}),
		},
		{
			Input: "\x1b[1;3;4;5;7m1\x1b[22;23;24;25;27m2",
			Expected: cells("12", cell.Style{
				Args: cell.Blink | cell.Underline | cell.Reverse | cell.Bold | cell.Italic,
			}, cell.Style{}),
		},
		{
			Input:    "\033[m123",
			Expected: cells("123"),
		},
		{
			Input:    "\033[41m1\033[49m23",
			Expected: cells("123", cell.Style{Bg: "41"}, cell.Style{}),
		},
	}
	buf := make([]cell.Cell, 0, 256)
	for i, test := range tt {
		got := cell.Parse(test.Input, &buf)
		if !slices.Equal(got, test.Expected) {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}

func TestParseFromTo(t *testing.T) {
	tt := []struct {
		Input    string
		Expected []cell.Cell
		From, To cell.Style
	}{
		{Input: "\033[31mHello\033[0m World",
			Expected: append(
				cells("Hello", cell.Style{Fg: "31", Args: cell.Bold}),
				cells(" World", cell.Style{})...,
			),
			From: cell.Style{Args: cell.Bold},
			To:   cell.Style{},
		},
		{Input: "1\033[33m2\033[3m3\033[5m4\033[0m",
			Expected: cells("1234", cell.Style{Args: cell.Bold, Fg: "90"}, cell.Style{Fg: "33", Args: cell.Bold}, cell.Style{Fg: "33", Args: cell.Bold | cell.Italic}, cell.Style{Fg: "33", Args: cell.Bold | cell.Italic | cell.Blink}),
			From:     cell.Style{Args: cell.Bold, Fg: "90"},
			To:       cell.Style{},
		},
		{Input: "",
			Expected: nil,
			From:     cell.Style{Args: cell.Underline, Fg: "90"},
			To:       cell.Style{Args: cell.Underline, Fg: "90"},
		},
		{Input: "Hello",
			Expected: cells("Hello", cell.Style{Args: cell.Underline, Fg: "90"}),
			From:     cell.Style{Args: cell.Underline, Fg: "90"},
			To:       cell.Style{Args: cell.Underline, Fg: "90"},
		},
	}
	buf := make([]cell.Cell, 0, 128)
	for i, test := range tt {
		got, to := cell.ParseFromTo(test.Input, &buf, test.From)
		if !slices.Equal(got, test.Expected) {
			t.Errorf("#%d: expected cells %v, but got %v", i, test.Expected, got)
		}
		if to != test.To {
			t.Errorf("#%d: expected to %v, but got %v", i, test.To, to)
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
		{
			name:     "bright black",
			last:     cell.Style{Fg: "90"},
			new:      cell.Style{Fg: "90"},
			expected: "",
		},

		{
			name:     "reverse - set",
			last:     cell.Style{Fg: "31", Args: cell.Reverse},
			new:      cell.Style{Fg: "31"},
			expected: "\033[27m",
		},

		{
			name:     "reverse - reset",
			last:     cell.Style{},
			new:      cell.Style{Args: cell.Reverse},
			expected: "\033[7m",
		},

		{
			name:     "reverse - set",
			last:     cell.Style{Fg: "31", Args: cell.Blink},
			new:      cell.Style{Fg: "31"},
			expected: "\033[25m",
		},

		{
			name:     "reverse - reset",
			last:     cell.Style{},
			new:      cell.Style{Args: cell.Blink},
			expected: "\033[5m",
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

func TestToString(t *testing.T) {
	tt := []struct {
		input    [][]cell.Cell
		expected string
	}{
		{
			input: [][]cell.Cell{
				cells("123"),
			},
			expected: "123",
		},
		{
			input: [][]cell.Cell{
				cells("456"),
				cells("123"),
			},
			expected: "456\n123",
		},
		{
			input: [][]cell.Cell{
				cells("456", cell.Style{Fg: "30", Args: cell.Italic}),
				cells("123"),
			},
			expected: "\033[3;30m456\033[0m\n123",
		},
		{
			input: [][]cell.Cell{
				cells("text", cell.Style{Args: cell.Bold | cell.Italic | cell.Underline}),
			},
			expected: "\x1b[1;3;4mtext\x1b[0m",
		},
		{
			input:    [][]cell.Cell{},
			expected: "",
		},
	}

	for i, tc := range tt {
		got := cell.ToString(tc.input)
		if got != tc.expected {
			t.Errorf("#%d: expected %v, but got %v", i, []byte(tc.expected), []byte(got))
		}
	}
}

func TestParseMultiline(t *testing.T) {
	tt := []struct {
		Input    string
		Expected [][]cell.Cell
	}{
		{Input: "",
			Expected: nil,
		},
		{Input: "Hello\nH",
			Expected: [][]cell.Cell{
				cells("Hello"),
				cells("H    "),
			},
		},
	}
	for i, test := range tt {
		got := cell.ParseMultiline(test.Input)
		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("#%d: expected cells %v, but got %v", i, test.Expected, got)
		}
	}
}
