package tui_test

import (
	"testing"

	"github.com/romanSPB15/tui-compose/v3"
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
