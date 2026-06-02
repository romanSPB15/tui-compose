//go:build debug

package tui

import (
	"fmt"
	"os"
)

const DEBUG = true

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
func (a *app) LogInfo(message string, args ...any) {
	fmt.Fprintf(a.log, message+"\r\n", args...)
}

// LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит
func (a *app) LogFatal(message string, args ...any) {
	recoveryScreen(fmt.Sprintf(message, args...))
	fmt.Fprintf(a.log, message+"\r\n", args...)
	os.Exit(1)
}
