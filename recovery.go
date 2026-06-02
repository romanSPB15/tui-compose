package tui

import (
	"fmt"
	"strings"
	"time"
)

func recoveryScreen(message string) {
	currentApp.Quit()
	fmt.Fprint(currentApp.f, "\033[2J\033[H")
	fmt.Fprint(currentApp.f, "\033[44m")
	time.Sleep(time.Millisecond)
	w := currentApp.Window().Width()
	format := fmt.Sprintf("%%-%ds", w) + "\r\n"
	fmt.Fprintf(currentApp.f, format, "Ошибка")
	fmt.Fprintf(currentApp.f, format, message)
	fmt.Fprintf(currentApp.f, format, "Нажмите ENTER для выхода...")
	for range currentApp.Window().Height() - 4 {
		fmt.Println(strings.Repeat(" ", w))
	}
	fmt.Fprint(currentApp.f, "\033[0m")
	fmt.Scanln()
}
