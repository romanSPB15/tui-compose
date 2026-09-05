package tui

func isFocusable(w EventHandler) bool {
	ev := &CheckFocusableEvent{}
	w.Send(ev)
	return ev.Result
}

func (wnd *window) FocusedWidget() EventHandler {
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
		wnd.setFocusTo(idx)
		return
	}
}

func (wnd *window) setFocusTo(idx int) {
	if wnd.focusIndex != -1 {
		wnd.focusableWidgets[wnd.focusIndex].Send(FocusEvent{Focused: false})
	}
	wnd.focusIndex = idx
	wnd.focusableWidgets[idx].Send(FocusEvent{Focused: true})
	wnd.Do(wnd.Redraw)
}

func (wnd *window) SetFocus(f EventHandler) bool {
	if !wnd.focusChange || !isFocusable(f) {
		return false
	}
	for i, w := range wnd.focusableWidgets {
		if w == f {
			wnd.setFocusTo(i)
			return true
		}
	}
	return false
}

func (wnd *window) ClearFocus() {
	if wnd.focusIndex != -1 {
		wnd.focusableWidgets[wnd.focusIndex].Send(FocusEvent{Focused: false})
		wnd.focusIndex = -1
		wnd.Do(wnd.Redraw)
	}
}

func (wnd *window) Disable() {
	wnd.focusChange = false
	wnd.ClearFocus()
}

func (wnd *window) Enable() {
	wnd.focusChange = true
}

func (wnd *window) FocusedIndex() int {
	return wnd.focusIndex
}

func (wnd *window) SetIndex(idx int) {
	if idx < -1 || idx >= len(wnd.focusableWidgets) {
		return
	}

	if wnd.focusIndex == idx {
		return
	}

	if wnd.focusIndex != -1 && wnd.focusIndex < len(wnd.focusableWidgets) {
		wnd.focusableWidgets[wnd.focusIndex].Send(FocusEvent{Focused: false})
	}

	if idx == -1 {
		wnd.focusIndex = -1
		wnd.Redraw()
		return
	}

	wnd.focusIndex = idx
	wnd.focusableWidgets[idx].Send(FocusEvent{Focused: true})
	wnd.Redraw()
}
