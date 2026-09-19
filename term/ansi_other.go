//go:build !windows

package term

import "os"

// EnableANSIWindows включает ENABLE_VIRTUAL_TERMINAL_PROCESSING на Windows для поддержки ANSI.
// На других ОС заглушка.
func EnableANSIWindows() {
}

func EnableANSIWindowsFile(f *os.File) {
}
