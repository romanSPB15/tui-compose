//go:build !no_recovery

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
	fmt.Fprint(currentWindow.f, "\033[3J")
	fmt.Fprint(currentWindow.f, "\033[44m")
	time.Sleep(time.Millisecond * 300)
	w := currentWindow.Width()
	format := fmt.Sprintf("%%-%ds", w) + "\r\n"
	fmt.Fprintf(currentWindow.f, format, "TUI Compose Framework")
	fmt.Fprintf(currentWindow.f, format, message)
	fmt.Fprintf(currentWindow.f, format, "Нажмите ENTER для выхода...")
	for range currentWindow.Height() - 4 {
		fmt.Fprintln(currentWindow.f, strings.Repeat(" ", w))
	}
	fmt.Fprint(currentWindow.f, strings.Repeat(" ", w))
	fmt.Fprint(currentWindow.f, "\033[0m")
	fmt.Scanln()
}
