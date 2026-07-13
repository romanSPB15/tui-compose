package term

import (
	"os"

	"golang.org/x/term"
)

// Width возвращает ширину терминала в символах.
// В случае ошибки возвращает 0.
func Width() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0
	}
	return width
}

// Height возвращает высоту терминала в строках.
// В случае ошибки возвращает 0.
func Height() int {
	_, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0
	}
	return height
}

// Size возвращает ширину и высоту терминала.
// В случае ошибки возвращает (0, 0).
func Size() (int, int) {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0
	}
	return w, h
}
