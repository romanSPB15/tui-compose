package term

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"golang.org/x/term"
)

var (
	chans  []chan []byte
	h      []func([]byte)
	mx     sync.Mutex
	stopCh chan struct{}

	started bool
)

// OnInput подписывается на ввод.
// Если fn != nil, добавляет функцию-обработчик.
// Если fn == nil, создаёт и возвращает новый канал.
func OnInput(fn func([]byte)) <-chan []byte {
	if fn == nil {
		ch := make(chan []byte, 16)
		mx.Lock()
		chans = append(chans, ch)
		mx.Unlock()
		return ch
	}

	mx.Lock()
	h = append(h, fn)
	mx.Unlock()

	return nil
}

// Start запускает чтение из reader.
func Start(input io.Reader) {
	if input == nil {
		input = os.Stdin
	}
	mx.Lock()
	defer mx.Unlock()
	if started {
		return
	}
	started = true
	stopCh = make(chan struct{})
	go readLoop(input)
}

// IsStarted возращает, запущено ли чтение.
func IsStarted() bool {
	mx.Lock()
	defer mx.Unlock()
	v := started
	return v
}

// Start останавливает чтение из stdin.
func Stop() {
	mx.Lock()
	defer mx.Unlock()

	if !started {
		return
	}
	started = false

	close(stopCh)

	for _, ch := range chans {
		func() {
			defer func() { recover() }()
			close(ch)
		}()
	}
	chans = nil
	h = nil
}

func readLoop(r io.Reader) {
	buf := make([]byte, 1024)

	// Локальный случай: *os.File (tty, pipe) поддерживает SetReadDeadline.
	// Это позволяет периодически проверять stopCh без блокировки.
	//
	// SSH-случай: ssh.Channel — это io.Reader без deadline.
	// Читаем блокирующе; остановка происходит через закрытие канала
	// снаружи (в Stop() нужно вызвать r.Close() если r — io.Closer).
	type deadlineReader interface {
		SetReadDeadline(time.Time) error
	}
	dr, hasDeadline := r.(deadlineReader)

	for {
		if hasDeadline {
			dr.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		}

		n, err := r.Read(buf) // ← БЫЛО os.Stdin.Read, должно быть r.Read
		if err != nil {
			if os.IsTimeout(err) {
				select {
				case <-stopCh:
					return
				default:
					continue
				}
			}
			return
		}
		if n == 0 {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		mx.Lock()

		select {
		case <-stopCh:
			mx.Unlock()
			return
		default:
		}

		for _, ch := range chans {
			select {
			case <-stopCh:
				mx.Unlock()
				return
			case ch <- data:
			default:
			}
		}

		for _, h := range h {
			func() {
				defer func() {
					if err := recover(); err != nil {
						fmt.Printf("term: panic: %v\r\n", err)
					}
				}()
				h(data)
			}()
		}

		mx.Unlock()
	}
}

var (
	ErrorNotRaw = errors.New("term: terminal not in raw mode")
)

var (
	oldState *term.State
	oldFile  *os.File
)

// MakeRawFile вводит переданный файл в raw режим.
func MakeRawFile(f *os.File) error {
	s, err := term.MakeRaw(int(f.Fd()))
	if err == nil {
		oldState = s
		oldFile = f
	}
	return err
}

// MakeRaw вводит os.Stdin в raw режим (для совместимости).
func MakeRaw() error {
	return MakeRawFile(os.Stdin)
}

// Restore выводит из raw тот файл, который был введён MakeRawFile.
func Restore() error {
	if oldState == nil || oldFile == nil {
		return ErrorNotRaw
	}
	return term.Restore(int(oldFile.Fd()), oldState)
}

// Width возвращает ширину терминала в символах.
// В случае ошибки возвращает 0.
func Width() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0
	}
	return width
}

// Height возвращает высоту терминала в строках.
// В случае ошибки возвращает 0.
func Height() int {
	_, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return height
}

// Size возвращает ширину и высоту терминала.
// В случае ошибки возвращает (0, 0).
func Size() (int, int) {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return 0, 0
	}
	return w, h
}

// SizeFd возвращает размеры терминала по заданному дескриптору.
// В случае ошибки возвращает (0, 0).
func SizeFd(fd uintptr) (int, int) {
	w, h, err := term.GetSize(int(fd))
	if err != nil {
		return 0, 0
	}
	return w, h
}

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
