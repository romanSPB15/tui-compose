package tui

import (
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

// EnableANSI() включает поддержку ANSI в терминале в случае если нет(Windows)
func EnableANSI() {
	if runtime.GOOS != "windows" {
		return
	}
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err != nil {
		return
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(stdout, mode)
}
