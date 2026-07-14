//go:build windows

package tui

import (
	"fmt"
	"time"
)

func (wnd *window) startScreenResizeChecker() {
	prevW, prevH := wnd.Width(), wnd.Height()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			newW, newH := wnd.Width(), wnd.Height()
			if newW != prevW || newH != prevH {
				if newH < prevH {
					fmt.Fprintf(wnd.f, "\033[%d;1H\033[J\033[H", newH+1)
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
