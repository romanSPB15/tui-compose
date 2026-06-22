// Добавлено в TUI 3.0.0.
package tui

import "unicode/utf8"

type Key uint16

type KeyboardEvent struct {
	Key  Key
	Rune rune
}

const (
	KeyF1 Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
)

const (
	KeyCtrlTilde      Key = 0x00
	KeyCtrl2          Key = 0x00
	KeyCtrlSpace      Key = 0x00
	KeyCtrlA          Key = 0x01
	KeyCtrlB          Key = 0x02
	KeyCtrlC          Key = 0x03
	KeyCtrlD          Key = 0x04
	KeyCtrlE          Key = 0x05
	KeyCtrlF          Key = 0x06
	KeyCtrlG          Key = 0x07
	KeyBackspace      Key = 0x08
	KeyCtrlH          Key = 0x08
	KeyTab            Key = 0x09
	KeyCtrlI          Key = 0x09
	KeyCtrlJ          Key = 0x0A
	KeyCtrlK          Key = 0x0B
	KeyCtrlL          Key = 0x0C
	KeyEnter          Key = 0x0D
	KeyCtrlM          Key = 0x0D
	KeyCtrlN          Key = 0x0E
	KeyCtrlO          Key = 0x0F
	KeyCtrlP          Key = 0x10
	KeyCtrlQ          Key = 0x11
	KeyCtrlR          Key = 0x12
	KeyCtrlS          Key = 0x13
	KeyCtrlT          Key = 0x14
	KeyCtrlU          Key = 0x15
	KeyCtrlV          Key = 0x16
	KeyCtrlW          Key = 0x17
	KeyCtrlX          Key = 0x18
	KeyCtrlY          Key = 0x19
	KeyCtrlZ          Key = 0x1A
	KeyEsc            Key = 0x1B
	KeyCtrlLsqBracket Key = 0x1B
	KeyCtrl3          Key = 0x1B
	KeyCtrl4          Key = 0x1C
	KeyCtrlBackslash  Key = 0x1C
	KeyCtrl5          Key = 0x1D
	KeyCtrlRsqBracket Key = 0x1D
	KeyCtrl6          Key = 0x1E
	KeyCtrl7          Key = 0x1F
	KeyCtrlSlash      Key = 0x1F
	KeyCtrlUnderscore Key = 0x1F
	KeySpace          Key = 0x20
	KeyBackspace2     Key = 0x7F
	KeyCtrl8          Key = 0x7F
)

func parseKeyboardKey(data []byte) (Key, int) {
	if len(data) < 2 || data[0] != 0x1B {
		return 0, 0
	}

	if data[1] == '[' {
		if len(data) < 3 {
			return 0, 0
		}
		switch data[2] {
		case 'A':
			return KeyArrowUp, 3
		case 'B':
			return KeyArrowDown, 3
		case 'C':
			return KeyArrowRight, 3
		case 'D':
			return KeyArrowLeft, 3
		case 'H':
			return KeyHome, 3
		case 'F':
			return KeyEnd, 3
		case '5', '6': // PgUp/PgDn (CSI 5 ~, CSI 6 ~)
			if len(data) >= 4 && data[3] == '~' {
				if data[2] == '5' {
					return KeyPgup, 4
				}
				if data[2] == '6' {
					return KeyPgdn, 4
				}
			}
			return 0, 0
		case '1', '2', '3', '4': // Home, End, Insert, Delete (CSI 1 ~, CSI 2 ~, CSI 3 ~, CSI 4 ~)
			if len(data) >= 4 && data[3] == '~' {
				switch data[2] {
				case '1':
					return KeyHome, 4
				case '2':
					return KeyInsert, 4
				case '3':
					return KeyDelete, 4
				case '4':
					return KeyEnd, 4
				}
			}
			return 0, 0
		}
	}

	// F1-F4 с ESC O P/Q/R/S (старый стиль)
	if data[1] == 'O' && len(data) >= 3 {
		switch data[2] {
		case 'P':
			return KeyF1, 3
		case 'Q':
			return KeyF2, 3
		case 'R':
			return KeyF3, 3
		case 'S':
			return KeyF4, 3
		}
	}

	// F5-F12: ESC [ 1 5 ~, ESC [ 1 7 ~ и т.д.
	if data[1] == '[' && len(data) >= 5 && data[3] == '~' {
		switch data[2] {
		case '1':
			switch data[4] {
			case '5': // ESC [ 1 5 ~ -> F5
				return KeyF5, 5
			case '7': // ESC [ 1 7 ~ -> F6
				return KeyF6, 5
			case '9': // ESC [ 1 9 ~ -> F7
				return KeyF7, 5
			}
		case '2':
			switch data[4] {
			case '0': // ESC [ 2 0 ~ -> F8
				return KeyF8, 5
			case '1': // ESC [ 2 1 ~ -> F9
				return KeyF9, 5
			case '3': // ESC [ 2 3 ~ -> F10
				return KeyF10, 5
			case '4': // ESC [ 2 4 ~ -> F11
				return KeyF11, 5
			case '5': // ESC [ 2 5 ~ -> F12
				return KeyF12, 5
			}
		}
	}

	if len(data) == 1 {
		return KeyEsc, 1
	}

	if len(data) >= 2 && data[1] >= 0x20 && data[1] <= 0x7E {
		return KeyEsc, 2
	}

	return 0, 0
}

func parseKeyboardRune(data []byte) (rune, bool) {
	if len(data) == 0 {
		return 0, false
	}

	b := data[0]

	if b < 0x20 || b == 0x7F {
		return 0, false
	}

	if b < 0x80 {
		return rune(b), true
	}

	r, size := utf8.DecodeRune(data)
	if r == utf8.RuneError && size == 1 {
		return 0, false
	}
	return r, true
}

func (wnd *window) handleKeyboardInput(data []byte) {
	if len(data) == 0 {
		return
	}

	if key, n := parseKeyboardKey(data); n > 0 {
		wnd.doWithMessageAndWait(func() {
			for _, h := range wnd.keyHandlers {
				wnd.doWithMessage(func() {
					h(&KeyboardEvent{
						Key: key,
					})
				}, "keyboard handler")
			}
		}, "key handler")
		return
	}

	if r, ok := parseKeyboardRune(data); ok {
		for _, h := range wnd.keyHandlers {
			wnd.doWithMessage(func() {
				h(&KeyboardEvent{
					Rune: r,
				})
			}, "keyboard handler")
		}
	}

	if data[0] == 3 { // Ctrl+C
		wnd.Quit()
		return
	}

	wnd.LogInfo("Нераспознано: %v", data)
}
