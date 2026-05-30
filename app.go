package tui

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/x/term"
	"github.com/eiannone/keyboard"
	"github.com/google/uuid"
)

type pos struct {
	Line int
	Col  int
}

type app struct {
	comp          []Component
	f             *os.File
	stopCh        chan struct{}
	keyHandlers   map[keyboard.Key]func()
	currentPos    pos
	posComponents []pos
	window        Window
	log           io.WriteCloser
	debug         bool
	access        *accessManager
}

var currentApp *app

func (a *app) Components() []Component {
	return a.comp
}

func (a *app) Window() Window {
	return a.window
}

// Метод для перерисовки всех компонентов
func (a *app) Redraw() {
	a.access.components(func() {
		fmt.Fprint(a.f, "\033[2J\033[H")
		a.posComponents = []pos{}
		a.currentPos = pos{0, 0}
		for idx, c := range a.comp {
			if c != nil {
				c.setIndex(idx)
				switch c.DisplayMode() {
				case DisplayInline:
					if a.currentPos.Col+c.MaxWidth() >= a.window.Width() {
						a.currentPos.Col = 0
						a.currentPos.Line++
					}
					a.posComponents = append(a.posComponents, a.currentPos)

					fmt.Fprint(a.f, c.innerText()+strings.Repeat(" ", c.MaxWidth()-len([]rune(stripansi.Strip(c.innerText())))), "  ")
					a.currentPos.Col += c.MaxWidth() + 2
				case DisplayBlock:
					a.currentPos.Col = 0
					a.currentPos.Line++

					fmt.Fprintln(a.f)

					a.posComponents = append(a.posComponents, a.currentPos)

					fmt.Fprint(a.f, c.innerText()+strings.Repeat(" ", c.MaxWidth()-len([]rune(stripansi.Strip(c.innerText())))), "  ")

					fmt.Fprintln(a.f)

					a.currentPos.Col = 0
					a.currentPos.Line++
				case DisplayNewLine:

					a.posComponents = append(a.posComponents, a.currentPos)

					a.currentPos.Col = 0
					a.currentPos.Line++

					fmt.Fprintln(a.f)

				}
			}
		}
	})

}

// Перерисовать конкретный компонент
func (a *app) RedrawComponent(index int) {
	a.access.components(func() {
		a.LogInfo("RedrawComponent %v", a.posComponents)
		pos := a.posComponents[index]
		fmt.Fprintf(a.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
		a.LogInfo("%v %d", pos, index)
		fmt.Print(a.comp[index].innerText() + strings.Repeat(" ", a.comp[index].MaxWidth()-len(stripansi.Strip(a.comp[index].innerText()))))
	})
}

// Метод добавления компонент в массив.
func (a *app) AddComponents(c ...Component) {
	a.comp = append(a.comp, c...)
}

// Метод для очистки списка компонентов приложения
func (a *app) Clear() {
	a.comp = []Component{}
	a.posComponents = []pos{}
}

// cursor
// hide \033[?25l
// show \033[?25h

type accessManager struct {
	mtxComponents *sync.Mutex
	mtxEvents     *sync.Mutex
}

func newAccessManager() *accessManager {
	return &accessManager{
		mtxComponents: &sync.Mutex{},
		mtxEvents:     &sync.Mutex{},
	}
}

func (am *accessManager) components(f func()) {
	am.mtxComponents.Lock()
	f()
	am.mtxComponents.Unlock()
}

func (am *accessManager) events(f func()) {
	am.mtxEvents.Lock()
	f()
	am.mtxEvents.Unlock()
}

// Запуск консольного приложения. Если приложение закроется, то будет произведён graceful shutdown и выход из метода.
func (a *app) Run() {
	if !term.IsTerminal(currentApp.f.Fd()) {
		a.LogFatal("tui: stdout is not terminal")
	}
	if runtime.GOOS == "windows" { // Windows не поддерживает ANSI escape sequnces по умолчанию.
		EnableANSI()
	}
	currentApp = a
	a.Redraw()

	defer func() {
		if a.debug {
			a.log.Close()
		}
	}()

	fmt.Fprint(a.f, "\033[?25l")

	// Реализация graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	keys, err := keyboard.GetKeys(2)
	if err != nil {
		a.LogFatal("tui: keyboard errror")
	}

	go func() {
		<-stop

		close(a.stopCh)
	}()

	// Обработка нажатий клавиш
	go func() {
		for ev := range keys {
			if ev.Key == keyboard.KeyCtrlC {
				close(a.stopCh)
			}
			a.access.events(func() {
				if v, ok := a.keyHandlers[ev.Key]; ok {
					v()
				}
				if ev.Err != nil {
					if err.Error() == "operation canceled" {
						close(a.stopCh)
						return
					}
					a.LogFatal("tui: keyboard error")
				}
			})
		}
	}()

	<-a.stopCh
	fmt.Print("\033[?25l")
	fmt.Fprint(a.f, "\033[2J\033[H\033[?25h")
}

// Выход из приложения.
func (a *app) Quit() {
	close(a.stopCh)
}

// Возвращает канал сигнализации выхода - приложение закрыто или вызван метод App.Quit().
func (a *app) OnQuit() <-chan struct{} {
	return a.stopCh
}

// Метод создаёт объект приложения без логирования.
func NewApp() App {
	app := &app{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[keyboard.Key]func()),
		window: &window{}, debug: false, access: newAccessManager(),
	}
	currentApp = app
	return app
}

// Метод создаёт объект приложения без логирования.
func NewDebugApp() App {
	f, err := os.Create(fmt.Sprintf("debug_log_%s", uuid.New().String()))
	if err != nil {
		log.Fatal(err)
	}
	app := &app{log: f, f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[keyboard.Key]func()),
		window: &window{}, debug: true, access: newAccessManager(),
	}
	currentApp = app
	return app
}

// Метод добавляет обработчик нажатия клавиши.
func (a *app) AddKeyHandler(key keyboard.Key, h func()) {
	a.keyHandlers[key] = h
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Система логирования приложения

// Критическая ошибка с выходом.
// Ничего не выводит если приложение не создано как Debug
func (a *app) LogFatal(message string, args ...any) {
	if a.debug {
		fmt.Fprintf(a.log, message+"\r\n", args...)
	}
	os.Exit(1)
}

// Логирования сообщения как fmt.Printf(), но с переносом строки в конце.
// Ничего не выводит если приложение не создано как Debug
func (a *app) LogInfo(message string, args ...any) {
	if a.debug {
		fmt.Fprintf(a.log, message+"\r\n", args...)
	}
}
