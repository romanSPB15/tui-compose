//go:build !debug

package tui

import (
	"fmt"
	"os"
)

const DEBUG = false

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если при сборке использовался тег debug.
func (wnd *window) LogInfo(message string, args ...any) {}

// LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если при сборке использовался тег debug. Потом вых
func (wnd *window) LogFatal(message string, args ...any) {
	recoveryScreen(fmt.Sprintf(message, args...))
	os.Exit(1)
}
