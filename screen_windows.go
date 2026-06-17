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
				prevW, prevH = newW, newH
				if newH < prevH {
					fmt.Fprint(wnd.f, "\033[2J")
				}

				wnd.doWithMessageAndWait(wnd.Redraw, "window resize (Windows)")
			}
		case <-wnd.stopCh:
			return
		}
	}
}
