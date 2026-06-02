//go:build !no_recovery && debug

package tui

import (
	"fmt"
	"strings"
	"time"
)

func recoveryScreen(message string) {
	if currentWindow.IsRunned() {
		currentWindow.Quit()
	}
	fmt.Fprint(currentWindow.f, "\033[2J\033[H")
	fmt.Fprint(currentWindow.f, "\033[44m")
	time.Sleep(time.Millisecond * 300)
	w := currentWindow.Window().Width()
	format := fmt.Sprintf("%%-%ds", w) + "\r\n"
	fmt.Fprintf(currentWindow.f, format, "go-tui")
	fmt.Fprintf(currentWindow.f, format, message)
	fmt.Fprintf(currentWindow.f, format, "Нажмите ENTER для выхода...")
	for range currentWindow.Window().Height() - 4 {
		fmt.Println(strings.Repeat(" ", w))
	}
	fmt.Fprint(currentWindow.f, "\033[0m")
	fmt.Scanln()
}
