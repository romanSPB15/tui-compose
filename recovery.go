package tui

import (
	"fmt"

	term "github.com/romanSPB15/tui-compose/v4/term"
)

func recoveryScreen(message string) {
	fmt.Fprint(currentWindow.f, "\033[0m")
	fmt.Fprint(currentWindow.f, "\033[2J\033[H\033[?25h")
	fmt.Fprint(currentWindow.f, "\033[?1006l\033[?1000l")

	currentWindow.restoreOut()
	term.Restore()

	if currentWindow.stopCh != nil {
		close(currentWindow.stopCh)
	}

	fmt.Fprintln(currentWindow.f, message)
}
