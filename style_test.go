package tui_test

import (
	"testing"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
)

func TestStyleString(t *testing.T) {
	tt := []struct {
		Style    tui.Style
		Expected string
	}{
		{Style: tui.Style(0),
			Expected: "",
		},
		{Style: tui.Italic | tui.Bold,
			Expected: "\x1b[1;3m",
		},
		{Style: tui.Underline | tui.Reverse,
			Expected: "\x1b[4;7m",
		},
		{Style: tui.Reset,
			Expected: "\x1b[0m",
		},
		{Style: tui.Reset | tui.Blink,
			Expected: "\x1b[0m",
		},
		{Style: tui.Blink,
			Expected: "\x1b[5m",
		},
		{Style: tui.BgRed | tui.FrBlack,
			Expected: "\x1b[30;41m",
		},
		{Style: tui.BgBrightRed | tui.FrBrightBlack,
			Expected: "\x1b[90;101m",
		},
		{Style: tui.BgBrightRed | tui.FrBrightCyan | tui.Blink,
			Expected: "\x1b[96;101;5m",
		},
	}
	for i, test := range tt {
		if got := test.Style.String(); got != test.Expected {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}

func TestConvertToCellStyle(t *testing.T) {
	tt := []struct {
		Style    tui.Style
		Expected cell.Style
	}{
		{Style: tui.Style(0),
			Expected: cell.Style{},
		},
		{Style: tui.Italic | tui.Bold,
			Expected: cell.Style{Args: cell.Italic | cell.Bold},
		},
		{Style: tui.Underline | tui.Reverse,
			Expected: cell.Style{Args: cell.Underline | cell.Reverse},
		},
		{Style: tui.Reset,
			Expected: cell.Style{Args: cell.Reset},
		},
		{Style: tui.Reset | tui.Blink,
			Expected: cell.Style{Args: cell.Reset | cell.Blink},
		},
		{Style: tui.Blink,
			Expected: cell.Style{Args: cell.Blink},
		},
		{Style: tui.BgRed | tui.FrBlack,
			Expected: cell.Style{Bg: "41", Fg: "30"},
		},
		{
			Style:    tui.BgBrightRed | tui.FrBrightBlack,
			Expected: cell.Style{Fg: "90", Bg: "101"},
		},
		{
			Style:    tui.BgBrightRed | tui.FrBrightCyan | tui.Blink,
			Expected: cell.Style{Fg: "96", Bg: "101", Args: cell.Blink},
		},
	}
	for i, test := range tt {
		if got := tui.ConvertToCellStyle(test.Style); got != test.Expected {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}
