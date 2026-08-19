package term

import (
	"os/exec"
	"runtime"
)

// OpenURL открывает переданный URL в браузере пользователя.
// Поддерживает Windows, macOS и Linux.
func OpenURL(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
