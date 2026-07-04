package input

import (
	"os"
	"time"
)

var (
	stopCh chan struct{}
)

func Start(buf int) (mouse <-chan *MouseEvent, keyboard <-chan *KeyboardEvent) {
	stopCh = make(chan struct{})
	mouseW := make(chan *MouseEvent, buf)
	keyboardW := make(chan *KeyboardEvent, buf)
	go readLoop(mouseW, keyboardW)
	return mouseW, keyboardW
}

func Stop() {
	close(stopCh)
}

func readLoop(mouse chan<- *MouseEvent, keyboard chan<- *KeyboardEvent) {
	buf := make([]byte, 512)
	for {
		os.Stdin.SetReadDeadline(time.Now().Add(60 * time.Millisecond))
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

		data := buf[:n]

		if ev := ParseKeyboardInput(data); ev != nil {
			keyboard <- ev
		}
		if ev := ParseMouseEvent(data); ev != nil {
			mouse <- ev
		}
	}
}
