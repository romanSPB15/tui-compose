package tui

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/charmbracelet/x/term"
	"github.com/eiannone/keyboard"
	"github.com/google/uuid"
)

// Color — это код цвета.
type Color int

const NoColor Color = 0

// Обычные цвета.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Яркие цвета(работают не во всем терминалах).
const (
	BrightBlack Color = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

// DisplayMode — это режим отображения виджета.
type DisplayMode int

const (
	DisplayInline  DisplayMode = iota // В одну строку.
	DisplayBlock                      // На отдельной строке.
	DisplayNewLine                    // Перенос строки.
)

type pos struct {
	Line int
	Col  int
}

type task struct {
	done chan struct{}
	f    func()
	msg  string
}

type app struct {
	comp        []Widget
	f           *os.File
	stopCh      chan struct{}
	keyHandlers map[keyboard.Key]func()
	currentPos  pos
	posWidgets  []pos
	window      Window
	log         io.WriteCloser
	debug       bool
	runned      bool
	work        chan *task
}

var currentApp *app

// Widgets() возвращает список компонентов, добавленных в приложение.
func (a *app) Widgets() []Widget {
	return a.comp
}

// Window() возвращает интерфейс окна приложения. Из него можно получить длину и ширину окна в символах.
func (a *app) Window() Window {
	return a.window
}

// Redraw() перерисовывает все компоненты. Он потокобезопасен.
// Важно: такая перерисовка вызывает мерцание.
func (a *app) Redraw() {
	a.doWithMessage(func() {
		fmt.Fprint(a.f, "\033[2J\033[H")
		a.posWidgets = []pos{}
		a.currentPos = pos{0, 0}
		for idx, c := range a.comp {
			if c != nil {
				if len(stripansi.Strip(c.InnerText())) > c.MaxLength() {
					a.LogFatal("Ошибка перерисовки: MaxLength() не верен.")
				}
				c.SetIndex(idx)
				switch c.DisplayMode() {
				case DisplayInline:
					if a.currentPos.Col+c.MaxLength() >= a.window.Width() {
						a.currentPos.Col = 0
						a.currentPos.Line++
					}
					a.posWidgets = append(a.posWidgets, a.currentPos)

					fmt.Fprint(a.f, c.InnerText()+strings.Repeat(" ", c.MaxLength()-len([]rune(stripansi.Strip(c.InnerText())))))
					a.currentPos.Col += c.MaxLength()
				case DisplayBlock:
					a.currentPos.Col = 0
					a.currentPos.Line++

					fmt.Fprintln(a.f)

					a.posWidgets = append(a.posWidgets, a.currentPos)

					fmt.Fprint(a.f, c.InnerText()+strings.Repeat(" ", c.MaxLength()-len([]rune(stripansi.Strip(c.InnerText())))))

					fmt.Fprintln(a.f)

					a.currentPos.Col = 0
					a.currentPos.Line++
				case DisplayNewLine:

					a.posWidgets = append(a.posWidgets, a.currentPos)

					a.currentPos.Col = 0
					a.currentPos.Line++

					fmt.Fprintln(a.f)

				}
			}
		}
	}, "redraw all")
}

// RedrawWidget() перерисовывает конкретный компонент. Потокобезопасен.
// index - это номер компонента, который нужно перерисовать.
func (a *app) RedrawWidget(index int) {
	a.doWithMessage(func() {
		a.LogInfo("RedrawWidget %v", a.posWidgets)
		pos := a.posWidgets[index]
		fmt.Fprintf(a.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
		a.LogInfo("%v %d", pos, index)
		fmt.Print(a.comp[index].InnerText() + strings.Repeat(" ", a.comp[index].MaxLength()-len(stripansi.Strip(a.comp[index].InnerText()))))
	}, "redraw widget")
}

// AddWidgets() добавляет компонент в приложение. Потокобезопасен.
func (a *app) AddWidgets(c ...Widget) {
	a.doWithMessageAndWait(func() {
		a.comp = append(a.comp, c...)
	}, "add widget")
}

// Clear() очищает список компонентов приложения без перерисовки. Потокобезопасен.
func (a *app) Clear() {
	a.doWithMessageAndWait(func() {
		a.comp = []Widget{}
		a.posWidgets = []pos{}
	}, "clear")
}

// cursor
// hide \033[?25l
// show \033[?25h

// Run() - это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
func (a *app) Run() {
	defer func() {
		if a.debug {
			a.log.Close()
		}
		if err := recover(); err != nil {
			a.LogFatal("Произошла panic: %v", err)
		}
	}()
	fmt.Fprintln(a.f, "Загрузка...")
	if !term.IsTerminal(currentApp.f.Fd()) {
		fmt.Fprintln(a.f, "Приложение запущено не в терминале. Выход...")
		time.Sleep(time.Second * 3)
		a.LogFatal("tui: stdout is not terminal")
	}
	if runtime.GOOS == "windows" {
		EnableANSI()
	}
	a.runned = true

	fmt.Fprint(a.f, "\033[?25l")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	keys, err := keyboard.GetKeys(2)
	if err != nil {
		a.LogFatal("tui: keyboard errror")
	}

	go func() {
		<-stop
		select {
		case <-a.stopCh:
		default:
			close(a.stopCh)
		}
	}()

	// Обработка нажатий клавиш
	go func() {
		for ev := range keys {
			if ev.Key == keyboard.KeyCtrlC {
				close(a.stopCh)
			}
			a.doWithMessageAndWait(func() {
				if v, ok := a.keyHandlers[ev.Key]; ok {
					a.Do(v)
				} else if ev.Err != nil {
					if ev.Err.Error() == "operation canceled" {
						close(a.stopCh)
						return
					}
					a.LogFatal("tui: keyboard error")
				}
			}, "key handler")
		}
	}()

	a.Redraw()

	<-a.stopCh
}

// Quit() — это принудительный выход из приложения.
func (a *app) Quit() {
	close(a.stopCh)
}

// Run() возвращает канал сигнализации о выходе.
func (a *app) OnQuit() <-chan struct{} {
	return a.stopCh
}

// IsRunned() возращает true, если приложение уже запущено. Иначе возвращает false.
func (a *app) IsRunned() bool {
	return a.runned
}

const taskBufSize = 16

// NewApp() создаёт объект приложения без логирования.
func NewApp() App {
	app := &app{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[keyboard.Key]func()),
		window: &window{}, debug: false, work: make(chan *task, taskBufSize),
	}
	currentApp = app
	go app.runWorker()
	return app
}

// NewDebugApp() создаёт объект приложения с логированием.
func NewDebugApp() App {
	f, err := os.Create(fmt.Sprintf("debug_log_%s", uuid.New().String()))
	if err != nil {
		log.Fatal(err)
	}
	app := &app{log: f, f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[keyboard.Key]func()),
		window: &window{}, debug: true, work: make(chan *task, taskBufSize),
	}
	currentApp = app
	go app.runWorker()
	return app
}

// AddKeyHandler() добавляет обработчик нажатия клавиши.
func (a *app) AddKeyHandler(key keyboard.Key, h func()) {
	a.keyHandlers[key] = h
}

// LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
func (a *app) LogInfo(message string, args ...any) {
	if a.debug {
		fmt.Fprintf(a.log, message+"\r\n", args...)
	}
}

// LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит
func (a *app) LogFatal(message string, args ...any) {
	recoveryScreen(fmt.Sprintf(message, args...))
	if a.debug {
		fmt.Fprintf(a.log, message+"\r\n", args...)
	}
	os.Exit(1)
}

// Do() запускает функцию f в потоке GUI, что спасает от data racing при изменении виджетов.
func (a *app) Do(f func()) {
	a.work <- &task{f: f}
}

// Do() запускает функцию f в потоке GUI и ждёт завершения.
func (a *app) DoAndWait(f func()) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
	}
	a.work <- tsk
	<-tsk.done
}

func (a *app) doWithMessage(f func(), msg string) {
	a.work <- &task{
		f:   f,
		msg: msg,
	}
}

func (a *app) doWithMessageAndWait(f func(), msg string) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
		msg:  msg,
	}
	a.work <- tsk
	<-tsk.done
}

func (a *app) runWorker() {
	a.LogInfo("Воркер запущен...")
	for {
		select {
		case <-a.stopCh:
			a.runned = false
			keyboard.Close()
			fmt.Print("\033[?25l")
			fmt.Fprint(a.f, "\033[2J\033[H\033[?25h")
			a.LogInfo("Воркер остановлен...")
			return
		case tsk := <-a.work:
			if tsk.msg != "" {
				a.LogInfo("Принята задача: '%s'", tsk.msg)
			} else {
				a.LogInfo("Принята задача")
			}
			tsk.f()
			if tsk.done != nil {
				close(tsk.done)
			}
			if tsk.msg != "" {
				a.LogInfo("Завершена задача: '%s'", tsk.msg)
			} else {
				a.LogInfo("Завершена задача")
			}
		}
	}
}
