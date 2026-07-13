package tui

func (wnd *window) FocusedWidget() Focusable {
	if wnd.focusIndex >= 0 && wnd.focusIndex < len(wnd.focusableWidgets) {
		return wnd.focusableWidgets[wnd.focusIndex]
	}
	return nil
}

func (wnd *window) NextFocus() {
	if !wnd.focusChange || len(wnd.focusableWidgets) == 0 {
		return
	}
	start := wnd.focusIndex + 1
	if start >= len(wnd.focusableWidgets) {
		start = 0
	}
	for i := 0; i < len(wnd.focusableWidgets); i++ {
		idx := (start + i) % len(wnd.focusableWidgets)
		w := wnd.focusableWidgets[idx]
		if d, ok := w.(Disablable); ok && d.IsDisabled() {
			continue
		}

		wnd.setFocusTo(idx)
		return
	}
}

func (wnd *window) BeforeFocus() {
	if !wnd.focusChange || len(wnd.focusableWidgets) == 0 {
		return
	}
	start := wnd.focusIndex - 1
	if start < 0 {
		start = len(wnd.focusableWidgets) - 1
	}
	for i := 0; i < len(wnd.focusableWidgets); i++ {
		idx := (start - i + len(wnd.focusableWidgets)) % len(wnd.focusableWidgets)
		w := wnd.focusableWidgets[idx]
		if d, ok := w.(Disablable); ok && d.IsDisabled() {
			continue
		}
		wnd.setFocusTo(idx)
		return
	}
}

func (wnd *window) setFocusTo(idx int) {
	if wnd.focusIndex != -1 {
		wnd.focusableWidgets[wnd.focusIndex].OnBlur()
	}
	wnd.focusIndex = idx
	wnd.focusableWidgets[idx].OnFocus()
	wnd.Redraw()
}

func (wnd *window) SetFocus(f Focusable) bool {
	if !wnd.focusChange {
		return false
	}
	for i, w := range wnd.focusableWidgets {
		if w == f {
			if d, ok := w.(Disablable); ok && d.IsDisabled() {
				return false
			}
			wnd.setFocusTo(i)
			return true
		}
	}
	return false
}

func (wnd *window) ClearFocus() {
	if wnd.focusIndex != -1 {
		wnd.focusableWidgets[wnd.focusIndex].OnBlur()
		wnd.focusIndex = -1
		wnd.Redraw()
	}
}

func (wnd *window) Disable() {
	wnd.focusChange = false
	wnd.ClearFocus()
}
