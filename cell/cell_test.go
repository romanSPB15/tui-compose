package cell_test

import (
	"slices"
	"testing"

	"github.com/romanSPB15/tui-compose/v3/cell"
)

func cells(chars string, style cell.Style) []cell.Cell {
	res := make([]cell.Cell, len(chars))
	for i, ch := range chars {
		res[i] = cell.Cell{Char: ch, Style: style}
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
	}
	for i, test := range tt {
		got := cell.Parse(test.Input)
		if !slices.Equal(got, test.Expected) {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}
