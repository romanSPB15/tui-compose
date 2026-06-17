//go:build !windows

package tui

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (wnd *window) startScreenResizeChecker() {
	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)
	for range sigwinch {
		select {
		case <-sigwinch:
			fmt.Fprint(wnd.f, "\033[2J")
			wnd.doWithMessageAndWait(wnd.Redraw, "window resize (Unix)")
		case <-wnd.stopCh:
			return
		}
	}
}
