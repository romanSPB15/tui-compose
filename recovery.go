package tui

import (
	"fmt"

	term "github.com/romanSPB15/tui-compose/v4/term"
)

func (wnd *window) recoveryScreen(message string) {
	fmt.Fprint(wnd.f, "\033[0m")
	fmt.Fprint(wnd.f, "\033[2J\033[H\033[?25h")
	fmt.Fprint(wnd.f, "\033[?1006l\033[?1000l\033[?1049l")

	wnd.restoreOut()
	term.Restore()

	if wnd.stopCh != nil {
		close(wnd.stopCh)
	}

	fmt.Fprintln(wnd.f, message)
}
