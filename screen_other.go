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

	for {
		select {
		case <-sigwinch:
			newW, newH := wnd.Width(), wnd.Height()
			if newW != prevW || newH != prevH {
				if newH < prevH {
					fmt.Fprintf(wnd.f, "\033[%d;1H\033[J", newH+1)
				}
				prevW, prevH = newW, newH
				wnd.doWithMessage(func() {
					wnd.buf = nil
					wnd.Redraw()
				}, "buf reset")
			}
		case <-wnd.stopCh:
			return
		}
	}
}
