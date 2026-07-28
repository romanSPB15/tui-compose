//go:build !no_recovery

package tui

import (
	"fmt"
	"strings"
	"time"

	term "github.com/romanSPB15/tui-compose/v3/term"
)

func recoveryScreen(message string) {
	fmt.Fprint(currentWindow.f, "\033[0m")
	fmt.Fprint(currentWindow.f, "\033[2J\033[H\033[?25h")
	fmt.Fprint(currentWindow.f, "\033[?1006l\033[?1000l")
	currentWindow.restoreOut()
	term.Restore()

	fmt.Fprint(currentWindow.f, "\033[3J")
	fmt.Fprint(currentWindow.f, "\033[44m")
	time.Sleep(time.Millisecond * 300)
	wnd := currentWindow.Width()
	format := fmt.Sprintf("%%-%ds", wnd) + "\r\n"
	fmt.Fprintf(currentWindow.f, format, "TUI Compose Framework")
	fmt.Fprintf(currentWindow.f, format, message)
	fmt.Fprintf(currentWindow.f, format, "Нажмите ENTER для выхода...")
	for range currentWindow.Height() - 4 {
		fmt.Fprintln(currentWindow.f, strings.Repeat(" ", wnd))
	}
	fmt.Fprint(currentWindow.f, strings.Repeat(" ", wnd))
	fmt.Fprint(currentWindow.f, "\033[0m")
	fmt.Scanln()
}
