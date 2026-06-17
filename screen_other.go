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

	prevW, prevH := wnd.Width(), wnd.Height()

	for range sigwinch {
		select {
		case <-sigwinch:
			newW, newH := wnd.Width(), wnd.Height()
			if newW != prevW || newH != prevH {
				if newH < prevH {
					fmt.Fprint(wnd.f, "\033[2J")
				}
				prevW, prevH = newW, newH
				wnd.doWithMessageAndWait(wnd.Redraw, "window resize (Unix)")
			}
		case <-wnd.stopCh:
			return
		}
	}
}
