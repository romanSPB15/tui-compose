//go:build debug

package tui

import (
	"fmt"
	"os"
	"runtime/debug"
)

const DEBUG = true

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
func (wnd *window) LogInfo(message string, args ...any) {
	fmt.Fprintf(wnd.log, message+"\r\n", args...)
}

// LogFatal() логирует указанное сообщение вместе со стеком подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит
func (wnd *window) LogFatal(message string, args ...any) {
	msg := fmt.Sprintf(message, args...)

	fmt.Fprintf(wnd.log, msg+"\r\n", args...)

	stack := debug.Stack()
	fmt.Fprintf(wnd.log, "Stack trace:\n%s\n", stack)
	recoveryScreen(msg)
	os.Exit(1)
}
