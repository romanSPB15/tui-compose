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

// LogFatal() логирует указанное сообщение вместе со стеком подобно fmt.Printf() в файл, если приложение создано как Debug. Потом выходит
func (wnd *window) LogFatal(message string, args ...any) {
	msg := fmt.Sprintf(message+"\r\n", args...)
	fmt.Fprint(wnd.log, msg)

	stack := debug.Stack()
	fmt.Fprintf(wnd.log, "Stack trace:\n%s\n", stack)
	wnd.log.Close()
	wnd.recoveryScreen(msg)
	os.Exit(1)
}
