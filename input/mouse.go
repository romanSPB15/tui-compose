// Добавлено в TUI 3.1.0.
package input

import (
	"strconv"
	"strings"
)

// Point представляет собой 2D-вектор.
type Point struct {
	X, Y int
}

// MouseAction описывает тип события мыши.
type MouseAction uint8

const (
	MousePress MouseAction = iota
	MouseRelease
	MouseMove
	MouseWheelUp
	MouseWheelDown
	MouseWheelLeft
	MouseWheelRight
)

// NoButton означает отсутствие кнопки: скролл или чистое движение без нажатия.
const NoButton = -1

// MouseEvent представляет собой событие мыши в SGR-режиме (xterm 1006).
type MouseEvent struct {
	Action MouseAction
	Button int // 0=левый, 1=средний, 2=правый, NoButton для скролла/motion без кнопки
	Pos    Point
	Shift  bool
	Alt    bool
	Ctrl   bool
}

func parseMouseEvent(input []byte) *MouseEvent {
	s := string(input)
	if !strings.HasPrefix(s, "\x1b[<") {
		return nil
	}

	body := s[3:]
	if len(body) < 2 {
		return nil
	}

	// SGR: `M` — press / motion / wheel, `m` — release.
	suffix := body[len(body)-1]
	if suffix != 'M' && suffix != 'm' {
		return nil
	}
	body = body[:len(body)-1]

	parts := strings.Split(body, ";")
	if len(parts) != 3 {
		return nil
	}

	cb, err := strconv.Atoi(parts[0])
	if err != nil || cb < 0 {
		return nil
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	ev := &MouseEvent{
		Button: NoButton,
		Pos:    Point{X: x - 1, Y: y - 1},
		Shift:  cb&4 != 0,
		Alt:    cb&8 != 0,
		Ctrl:   cb&16 != 0,
	}

	const (
		motionBit = 32
		wheelBit  = 64
	)

	switch {
	case cb&wheelBit != 0:
		switch cb & 3 {
		case 0:
			ev.Action = MouseWheelUp
		case 1:
			ev.Action = MouseWheelDown
		case 2:
			ev.Action = MouseWheelLeft
		case 3:
			ev.Action = MouseWheelRight
		}

	case cb&motionBit != 0:
		ev.Action = MouseMove
		if b := cb & 3; b != 3 {
			ev.Button = b
		}

	case suffix == 'm':
		ev.Action = MouseRelease
		if b := cb & 3; b != 3 {
			ev.Button = b
		}

	default:
		ev.Action = MousePress
		if b := cb & 3; b != 3 {
			ev.Button = b
		}
	}

	return ev
}
