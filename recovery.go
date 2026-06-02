package tui

import (
	"fmt"
	"strings"
)

func recoveryScreen(message string) {
	fmt.Fprint(currentApp.f, "\033[2J\033[H")
	fmt.Fprint(currentApp.f, "\033[44m")
	w := currentApp.Window().Width()
	fmt.Fprintf(currentApp.f, fmt.Sprintf("%%-%ds", w)+"\r\n", "Ошибка")
	fmt.Fprintf(currentApp.f, fmt.Sprintf("%%-%ds", w)+"\r\n", message)
	fmt.Fprintf(currentApp.f, fmt.Sprintf("%%-%ds", w)+"\r\n", "Нажмите ENTER для выхода...")
	for range currentApp.Window().Height() - 3 {
		fmt.Println(strings.Repeat(" ", w))
	}
	fmt.Fprint(currentApp.f, "\033[0m")
	fmt.Scanln()
}
