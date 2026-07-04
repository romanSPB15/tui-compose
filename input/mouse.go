// Добавлено в TUI 3.1.0.
package input

import (
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

type MouseEvent struct {
	Button int // 0=левый, 1=средний, 2=правый
	Pos    Point
}

func ParseMouseEvent(input []byte) *MouseEvent {
	s := string(input)
	if !strings.HasPrefix(s, "\x1b[<") {
		return nil
	}
	rest := strings.TrimPrefix(s, "\x1b[<")
	rest = strings.TrimSuffix(rest, "m")

	parts := strings.Split(rest, ";")
	if len(parts) != 3 {
		return nil
	}
	btn, err := strconv.Atoi(parts[0])
	if err != nil {
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

	button := btn & 0x03

	return &MouseEvent{
		Button: button,
		Pos:    Point{x - 1, y - 1},
	}
}
