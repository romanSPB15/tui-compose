// Добавлено в TUI 3.0.0.
package tui

type KeyboardEvent struct {
	Key  Key
	Rune rune
	Alt  bool
}

type Key uint16

const (
	KeyUnknown Key = iota

	KeyF1
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

	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlI
	KeyCtrlJ
	KeyCtrlK
	KeyCtrlL
	KeyCtrlM
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ

	KeyEnter
	KeySpaсe
	KeyPgUp
	KeyPgDown
	KeySlash
	KeyReverseSlash
	KeyTab
	KeyShiftTab

	KeyBackspase
	KeyDelete
	KeyInsert
	KeyHome
	KeyEnd

	KeyArrowUp
	KeyArrowRight
	KeyArrowDown
	KeyArrowLeft
)

func parseAnsiKeyboardInput(data []byte) (rune, Key) {
	if len(data) == 1 {
		v := data[0]
		switch {
		case v < 27:
			if v == 13 {
				return 0, KeyEnter
			}
			return 0, KeyCtrlA + Key(v-1)
		case v > 64 && v < 91:
			return 'A' + rune(v-65), KeyUnknown
		case v > 96 && v < 123:
			return 'a' + rune(v-97), KeyUnknown
		case v > 47 && v < 58:
			return '0' + rune(v-48), KeyUnknown
		case v == 32:
			return ' ', KeySpaсe
		case v == 47:
			return 0, KeySlash
		case v == 92:
			return 0, KeyReverseSlash
		case v == 9:
			return 0, KeyTab
		case v == 35:
			return 0, KeyBackspase
		default:
			return 0, KeyUnknown
		}
	}
	m := map[string]Key{
		string([]byte{27, 91, 53, 126}):     KeyPgUp,
		string([]byte{27, 91, 54, 126}):     KeyPgDown,
		string([]byte{27, 91, 90}):          KeyShiftTab,
		string([]byte{27, 91, 51, 126}):     KeyDelete,
		string([]byte{27, 91, 70}):          KeyEnd,
		string([]byte{27, 91, 72}):          KeyHome,
		string([]byte{27, 91, 50, 126}):     KeyInsert,
		string([]byte{27, 79, 80}):          KeyF1,
		string([]byte{27, 79, 81}):          KeyF2,
		string([]byte{27, 79, 82}):          KeyF3,
		string([]byte{27, 79, 83}):          KeyF4,
		string([]byte{27, 91, 49, 53, 126}): KeyF5,
		string([]byte{27, 91, 49, 55, 126}): KeyF6,
		string([]byte{27, 91, 49, 56, 126}): KeyF7,
		string([]byte{27, 91, 49, 57, 126}): KeyF8,
		string([]byte{27, 91, 50, 48, 126}): KeyF9,
		string([]byte{27, 91, 50, 49, 126}): KeyF10,
		string([]byte{27, 91, 50, 51, 126}): KeyF11,
		string([]byte{27, 91, 50, 52, 126}): KeyF12,
		string([]byte{27, 91, 65}):          KeyArrowUp,
		string([]byte{27, 91, 67}):          KeyArrowRight,
		string([]byte{27, 91, 66}):          KeyArrowDown,
		string([]byte{27, 91, 68}):          KeyArrowLeft,
	}
	key := string(data)
	if v, ok := m[key]; ok {
		return 0, v
	}

	return 0, KeyUnknown
}

func (wnd *window) handleKeyboardInput(data []byte) {
	if len(data) == 0 {
		return
	}
	if data[0] == 3 { // Ctrl+C
		wnd.Quit()
		return
	}

	var alt = false
	var r rune
	var k Key

	if len(data) > 2 && data[0] == 27 {
		alt = true
		r, k = parseAnsiKeyboardInput(data[1:])
	} else {
		r, k = parseAnsiKeyboardInput(data)
	}

	if r == 0 && k == KeyUnknown {
		wnd.LogInfo("Нераспознано: %v", data)
		return
	}

	wnd.doWithMessageAndWait(func() {
		for _, h := range wnd.keyHandlers {
			wnd.doWithMessage(func() {
				h(&KeyboardEvent{
					Key:  k,
					Rune: r,
					Alt:  alt,
				})
			}, "keyboard handler")
		}
	}, "key handler")
}
