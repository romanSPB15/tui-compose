//go:build windows

package term

import (
	"os"

	"golang.org/x/sys/windows"
)

// EnableANSIWindows включает ENABLE_VIRTUAL_TERMINAL_PROCESSING на Windows для поддержки ANSI.
// На других ОС заглушка.
func EnableANSIWindows() {
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(stdout, &mode); err != nil {
		return
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(stdout, mode)
}

// EnableANSIWindows включает ENABLE_VIRTUAL_TERMINAL_PROCESSING на Windows для поддержки ANSI в указанном файле.
// На других ОС заглушка.
func EnableANSIWindowsFile(f *os.File) {
	if f == nil {
		return
	}
	h := windows.Handle(f.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(h, &mode); err != nil {
		return
	}
	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(h, mode)
}
