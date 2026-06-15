//go:build debug

package tui

import (
	"fmt"
	"os"
)

const DEBUG = true

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
func (wnd *window) LogInfo(message string, args ...any) {
	fmt.Fprintf(wnd.log, message+"\r\n", args...)
}

// LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит
func (wnd *window) LogFatal(message string, args ...any) {
	recoveryScreen(fmt.Sprintf(message, args...))
	fmt.Fprintf(wnd.log, message+"\r\n", args...)
	os.Exit(1)
}
