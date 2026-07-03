// Добавлено в TUI 3.1.0.
package input

import (
	"fmt"
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

func ParseMouseEvent(input []byte) (*MouseEvent, error) {
	s := string(input)
	if !strings.HasPrefix(s, "\x1b[<") {
		return nil, fmt.Errorf("не SGR последовательность")
	}
	rest := strings.TrimPrefix(s, "\x1b[<")
	rest = strings.TrimSuffix(rest, "m")

	parts := strings.Split(rest, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("неверный формат: %v", parts)
	}
	btn, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	button := btn & 0x03

	return &MouseEvent{
		Button: button,
		Pos:    Point{y - 1, x - 1},
	}, nil
}
