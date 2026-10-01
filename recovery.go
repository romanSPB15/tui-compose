package tui

import (
	"fmt"
)

func (wnd *window) recoveryScreen(message string) {
	wnd.t.Close()

	if wnd.stopCh != nil {
		close(wnd.stopCh)
	}

	fmt.Fprintln(wnd.t, message)
}
