package input

import (
	"io"
	"os"
	"sync"

	"github.com/romanSPB15/tui-compose/v4/term"
)

var (
	mouseCh    chan *MouseEvent
	keyboardCh chan *KeyboardEvent
	mu         sync.Mutex
	started    bool
)

// Start запускает чтение событий ввода и возвращает каналы для событий мыши
// и клавиатуры.
//
// Если чтение уже запущено, повторный вызов возвращает уже существующие
// каналы, а buf игнорируется.
func Start(input io.Reader, buf int) (<-chan *MouseEvent, <-chan *KeyboardEvent) {
	if input == nil {
		input = os.Stdin
	}
	mu.Lock()
	defer mu.Unlock()

	if started {
		return mouseCh, keyboardCh
	}

	term.Start(input)

	mouseCh = make(chan *MouseEvent, buf)
	keyboardCh = make(chan *KeyboardEvent, buf)

	mouse, keyboard := mouseCh, keyboardCh

	term.OnInput(func(data []byte) {
		if ev := parseMouseEvent(data); ev != nil {
			select {
			case mouse <- ev:
			default:
			}
			return
		}

		if ev := parseKeyboardInput(data); ev != nil {
			select {
			case keyboard <- ev:
			default:
			}
		}
	})

	started = true
	return mouseCh, keyboardCh
}

// Stop останавливает чтение событий и закрывает каналы.
// Идемпотентен: повторный вызов без предшествующего Start ничего не делает.
func Stop() {
	mu.Lock()
	if !started {
		mu.Unlock()
		return
	}
	mouse, keyboard := mouseCh, keyboardCh
	mouseCh, keyboardCh = nil, nil
	started = false
	mu.Unlock()

	// Останавливаем ридер до закрытия каналов: после возврата term.Stop
	// callback гарантированно не будет вызван.
	term.Stop()

	close(mouse)
	close(keyboard)
}
