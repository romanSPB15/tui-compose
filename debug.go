//go:build debug

package tui

import (
	"fmt"
	"os"
)

const DEBUG = true

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
func (w *window) LogInfo(message string, args ...any) {
	fmt.Fprintf(w.log, message+"\r\n", args...)
}

// LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит
func (w *window) LogFatal(message string, args ...any) {
	recoveryScreen(fmt.Sprintf(message, args...))
	fmt.Fprintf(w.log, message+"\r\n", args...)
	os.Exit(1)
}
