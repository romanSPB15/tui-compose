package term

import (
	"errors"
	"fmt"
	"os"
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

func Start() {
	mx.Lock()
	defer mx.Unlock()
	if started {
		return
	}
	started = true
	stopCh = make(chan struct{})
	go readLoop()
}

func IsStarted() bool {
	mx.Lock()
	defer mx.Unlock()
	v := started
	return v
}

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

func readLoop() {
	buf := make([]byte, 1024)
	for {
		os.Stdin.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
		n, err := os.Stdin.Read(buf)
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

var old *term.State

// MakeRaw вводит терминал в raw режим.
func MakeRaw() error {
	s, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		old = s
	}
	return err
}

// MakeRaw выводит терминал из raw режима.
func Restore() error {
	if old == nil {
		return ErrorNotRaw
	}
	return term.Restore(int(os.Stdin.Fd()), old)
}
