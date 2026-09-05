package input

import (
	"sync"

	"github.com/romanSPB15/tui-compose/v4/term"
)

var (
	mouseCh    chan *MouseEvent
	keyboardCh chan *KeyboardEvent
	mu         sync.Mutex
	started    bool
	stopCh     chan struct{}
)

// Start запускает чтение событий ввода.
// Возвращает каналы для мыши и клавиатуры.
func Start(buf int) (<-chan *MouseEvent, <-chan *KeyboardEvent) {
	mu.Lock()
	defer mu.Unlock()

	if started {
		return mouseCh, keyboardCh
	}

	term.Start()

	mouseCh = make(chan *MouseEvent, buf)
	keyboardCh = make(chan *KeyboardEvent, buf)
	stopCh = make(chan struct{})

	term.OnInput(func(data []byte) {
		if ev := parseMouseEvent(data); ev != nil {
			select {
			case mouseCh <- ev:
			default:
			}
			return
		}

		if ev := parseKeyboardInput(data); ev != nil {
			select {
			case keyboardCh <- ev:
			default:
			}
		}
	})

	started = true
	return mouseCh, keyboardCh
}

// Stop останавливает чтение и закрывает каналы.
func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if !started {
		return
	}

	if stopCh != nil {
		close(stopCh)
		stopCh = nil
	}

	term.Stop()

	if mouseCh != nil {
		close(mouseCh)
		mouseCh = nil
	}
	if keyboardCh != nil {
		close(keyboardCh)
		keyboardCh = nil
	}

	started = false
}
