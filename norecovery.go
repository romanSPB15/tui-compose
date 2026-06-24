//go:build no_recovery

package tui

import "fmt"

func recoveryScreen(message string) {
	fmt.Fprint(currentWindow.f, "\033[2J\033[H")
}
