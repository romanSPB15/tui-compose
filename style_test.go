package tui_test

import (
	"testing"

	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
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
		Expected acell.Style
	}{
		{Style: tui.Style(0),
			Expected: acell.Style{},
		},
		{Style: tui.Italic | tui.Bold,
			Expected: acell.Style{Args: acell.Italic | acell.Bold},
		},
		{Style: tui.Underline | tui.Reverse,
			Expected: acell.Style{Args: acell.Underline | acell.Reverse},
		},
		{Style: tui.Blink,
			Expected: acell.Style{Args: acell.Blink},
		},
		{Style: tui.BgRed | tui.FrBlack,
			Expected: acell.Style{Bg: "41", Fg: "30"},
		},
		{
			Style:    tui.BgBrightRed | tui.FrBrightBlack,
			Expected: acell.Style{Fg: "90", Bg: "101"},
		},
		{
			Style:    tui.BgBrightRed | tui.FrBrightCyan | tui.Blink,
			Expected: acell.Style{Fg: "96", Bg: "101", Args: acell.Blink},
		},
	}
	for i, test := range tt {
		if got := tui.ConvertToCellStyle(test.Style); got != test.Expected {
			t.Errorf("#%d: expected %v, but got %v", i, test.Expected, got)
		}
	}
}
