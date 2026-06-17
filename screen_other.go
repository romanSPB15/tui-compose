//go:build !windows

package tui

import (
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
			wnd.doWithMessageAndWait(wnd.Redraw, "window resize (Unix)")
		case <-wnd.stopCh:
			return
		}
	}
}
