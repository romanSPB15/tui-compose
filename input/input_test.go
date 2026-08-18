package input

import (
	"reflect"
	"testing"
)

func TestParseAnsiKeyboardInput(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantRune rune
		wantKey  Key
	}{
		{"enter", []byte{13}, 0, KeyEnter},
		{"tab", []byte{9}, 0, KeyTab},
		{"space", []byte{32}, ' ', KeySpace},
		{"slash", []byte{47}, '/', KeySlash},
		{"backslash", []byte{92}, '\\', KeyReverseSlash},
		{"backspace", []byte{127}, 0, KeyBackspace},
		{"normal 'a'", []byte{'a'}, 'a', KeyUnknown},
		{"normal 'A'", []byte{'A'}, 'A', KeyUnknown},
		{"ctrl+a", []byte{1}, 0, KeyCtrlA},
		{"ctrl+z", []byte{26}, 0, KeyCtrlZ},
		{"unknown", []byte{0}, 0, KeyUnknown},

		{"pgup", []byte{27, 91, 53, 126}, 0, KeyPgUp},
		{"pgdown", []byte{27, 91, 54, 126}, 0, KeyPgDown},
		{"shifttab", []byte{27, 91, 90}, 0, KeyShiftTab},
		{"delete", []byte{27, 91, 51, 126}, 0, KeyDelete},
		{"end", []byte{27, 91, 70}, 0, KeyEnd},
		{"home", []byte{27, 91, 72}, 0, KeyHome},
		{"insert", []byte{27, 91, 50, 126}, 0, KeyInsert},
		{"f1", []byte{27, 79, 80}, 0, KeyF1},
		{"f2", []byte{27, 79, 81}, 0, KeyF2},
		{"f3", []byte{27, 79, 82}, 0, KeyF3},
		{"f4", []byte{27, 79, 83}, 0, KeyF4},
		{"f5", []byte{27, 91, 49, 53, 126}, 0, KeyF5},
		{"f6", []byte{27, 91, 49, 55, 126}, 0, KeyF6},
		{"f7", []byte{27, 91, 49, 56, 126}, 0, KeyF7},
		{"f8", []byte{27, 91, 49, 57, 126}, 0, KeyF8},
		{"f9", []byte{27, 91, 50, 48, 126}, 0, KeyF9},
		{"f10", []byte{27, 91, 50, 49, 126}, 0, KeyF10},
		{"f11", []byte{27, 91, 50, 51, 126}, 0, KeyF11},
		{"f12", []byte{27, 91, 50, 52, 126}, 0, KeyF12},
		{"up", []byte{27, 91, 65}, 0, KeyArrowUp},
		{"right", []byte{27, 91, 67}, 0, KeyArrowRight},
		{"down", []byte{27, 91, 66}, 0, KeyArrowDown},
		{"left", []byte{27, 91, 68}, 0, KeyArrowLeft},
		{"unknown seq", []byte{27, 91, 1}, 0, KeyUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, k := parseAnsiKeyboardInput(tt.data)
			if r != tt.wantRune || k != tt.wantKey {
				t.Errorf("got (%q, %v), want (%q, %v)", r, k, tt.wantRune, tt.wantKey)
			}
		})
	}
}

func TestParseKeyboardInput(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want *KeyboardEvent
	}{
		{"empty", []byte{}, nil},
		{"normal 'a'", []byte{'a'}, &KeyboardEvent{Key: KeyUnknown, Rune: 'a', Alt: false}},
		{"enter", []byte{13}, &KeyboardEvent{Key: KeyEnter, Rune: 0, Alt: false}},
		{"ctrl+a", []byte{1}, &KeyboardEvent{Key: KeyCtrlA, Rune: 0, Alt: false}},
		{"alt+a", []byte{27, 'a'}, &KeyboardEvent{Key: KeyUnknown, Rune: 'a', Alt: true}},
		{"alt+enter", []byte{27, 13}, &KeyboardEvent{Key: KeyEnter, Rune: 0, Alt: true}},
		{"alt+up", []byte{27, 27, 91, 65}, &KeyboardEvent{Key: KeyArrowUp, Rune: 0, Alt: true}},
		{"unknown", []byte{0}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseKeyboardInput(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				if got == nil {
					t.Errorf("got <nil>, want %v", *tt.want)
				} else {
					t.Errorf("got %v, want %v", got, *tt.want)
				}
			}
		})
	}
}

func TestParseMouseEvent(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want *MouseEvent
	}{
		{"left click", []byte("\x1b[<0;10;20m"), &MouseEvent{Button: 0, Pos: Point{X: 9, Y: 19}}},
		{"middle click", []byte("\x1b[<1;5;5m"), &MouseEvent{Button: 1, Pos: Point{X: 4, Y: 4}}},
		{"right click", []byte("\x1b[<2;1;1m"), &MouseEvent{Button: 2, Pos: Point{X: 0, Y: 0}}},
		{"invalid prefix", []byte("abc"), nil},
		{"invalid parts", []byte("\x1b[<0;10m"), nil},
		{"invalid number 1", []byte("\x1b[<a;10;20m"), nil},
		{"invalid number 2", []byte("\x1b[<1;a;20m"), nil},
		{"invalid number 3", []byte("\x1b[<1;10;am"), nil},
		{"empty", []byte{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMouseEvent(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
