//go:build windows

package tui

import "time"

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
				wnd.doWithMessageAndWait(func() {
					wnd.currentPos = pos{0, 0}
					wnd.index()
					wnd.Redraw()
				}, "window resize (Windows)")
			}
		case <-wnd.stopCh:
			return
		}
	}
}
