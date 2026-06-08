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

// ColorRGB — это цвет в RGB.
type ColorRGB struct {
	R, G, B uint8
}

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

var currentWindow *window

type window struct {
	comp        []Widget
	compF       []Focusable
	f           *os.File
	focusIndex  int
	stopCh      chan struct{}
	keyHandlers map[keyboard.Key]func()
	currentPos  pos
	posWidgets  []pos
	log         io.WriteCloser
	runned      bool
	work        chan *task
	focusChange bool
}

// Widgets() возвращает список компонентов, добавленных в приложение.
func (w *window) Widgets() []Widget {
	return w.comp
}

// Redraw() перерисовывает все компоненты. Он потокобезопасен.
// Важно: такая перерисовка вызывает мерцание.
func (w *window) Redraw() {
	w.doWithMessage(func() {
		fmt.Fprint(w.f, "\033[2J\033[H")

		for idx, c := range w.comp {
			if c != nil {
				if idx >= len(w.posWidgets) {
					w.LogFatal("позиция для виджета %d не найдена", idx)
				}
				pos := w.posWidgets[idx]
				fmt.Fprintf(w.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
				fmt.Fprint(w.f, c.InnerText())
			}
		}
	}, "redraw all")
}

func (w *window) index() {
	w.compF = []Focusable{}
	w.posWidgets = []pos{}
	w.currentPos = pos{0, 0}
	for idx, c := range w.comp {
		if c != nil {
			if len(stripansi.Strip(c.InnerText())) > c.MaxLength() {
				w.LogFatal("Ошибка индексации: MaxLength() не верен.")
			}
			if f, ok := c.(Focusable); ok {
				w.compF = append(w.compF, f)
			}
			c.SetIndex(idx)
			switch c.DisplayMode() {
			case DisplayInline:
				if w.currentPos.Col+c.MaxLength() >= w.Width() {
					w.currentPos.Col = 0
					w.currentPos.Line++
				}
				w.posWidgets = append(w.posWidgets, w.currentPos)

				w.currentPos.Col += c.MaxLength()
			case DisplayBlock:
				w.currentPos.Col = 0
				w.currentPos.Line++

				w.posWidgets = append(w.posWidgets, w.currentPos)

				w.currentPos.Col = 0
				w.currentPos.Line++
			case DisplayNewLine:
				w.posWidgets = append(w.posWidgets, w.currentPos)

				w.currentPos.Col = 0
				w.currentPos.Line++
			}
		}
	}
}

// RedrawWidget() перерисовывает конкретный компонент. Потокобезопасен.
// index - это номер компонента, который нужно перерисовать.
func (w *window) RedrawWidget(index int) {
	w.doWithMessage(func() {
		w.LogInfo("RedrawWidget %v", w.posWidgets)
		pos := w.posWidgets[index]
		fmt.Fprintf(w.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
		w.LogInfo("%v %d", pos, index)
		fmt.Print(w.comp[index].InnerText() + strings.Repeat(" ", w.comp[index].MaxLength()-len(stripansi.Strip(w.comp[index].InnerText()))))
	}, "redraw widget")
}

// AddWidgets() добавляет компонент в приложение. Потокобезопасен.
func (w *window) AddWidgets(c ...Widget) {
	w.doWithMessageAndWait(func() {
		w.comp = append(w.comp, c...)
	}, "add widget")
}

// Clear() очищает список компонентов приложения без перерисовки. Потокобезопасен.
func (w *window) Clear() {
	w.doWithMessageAndWait(func() {
		w.comp = []Widget{}
		w.compF = []Focusable{}
		w.posWidgets = []pos{}
	}, "clear")
}

// cursor
// hide \033[?25l
// show \033[?25h

// Run() - это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
func (w *window) Run() {
	defer func() {
		if DEBUG {
			w.log.Close()
		}
		if err := recover(); err != nil {
			w.LogFatal("Произошла panic: %v", err)
		}
	}()
	if !term.IsTerminal(w.f.Fd()) {
		fmt.Fprintln(w.f, "Приложение запущено не в терминале. Выход...")
		time.Sleep(time.Second * 3)
		w.LogFatal("tui: stdout is not terminal")
	}

	w.index()
	w.runned = true

	fmt.Fprint(w.f, "\033[?25l")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	keys, err := keyboard.GetKeys(2)
	if err != nil {
		w.LogFatal("tui: keyboard errror")
	}

	go func() {
		<-stop
		select {
		case <-w.stopCh:
		default:
			close(w.stopCh)
		}
	}()

	// Обработка нажатий клавиш
	go func() {
		for ev := range keys {
			if ev.Key == keyboard.KeyCtrlC {
				close(w.stopCh)
			}
			w.doWithMessageAndWait(func() {
				if v, ok := w.keyHandlers[ev.Key]; ok {
					w.Do(v)
				} else if ev.Err != nil {
					if ev.Err.Error() == "operation canceled" {
						close(w.stopCh)
						return
					}
					w.LogFatal("tui: keyboard error")
				}
			}, "key handler")
		}
	}()

	if w.focusChange && len(w.compF) != 0 {
		w.RegisterKeyHandler(keyboard.KeyArrowLeft, func() {
			if w.focusIndex <= 0 {
				return
			}
			w.compF[w.focusIndex].OnBlur()
			w.focusIndex--
			w.compF[w.focusIndex].OnFocus()
		})

		w.RegisterKeyHandler(keyboard.KeyArrowRight, func() {
			if w.focusIndex > len(w.compF)-2 {
				return
			}
			if w.focusIndex == -1 {
				w.compF[0].OnFocus()
				w.focusIndex = 0
				return
			}
			w.compF[w.focusIndex].OnBlur()
			w.focusIndex++
			w.compF[w.focusIndex].OnFocus()
		})
	}

	w.RegisterKeyHandler(keyboard.KeyEnter, func() {
		if w.focusIndex != -1 {
			if cl, ok := w.compF[w.focusIndex].(Clickable); ok {
				cl.OnClick()
			}
		}
	})

	w.Redraw()

	w.startScreenResizeChecker()

	if len(w.compF) != 0 {
		w.Do(func() {
			w.compF[0].OnFocus()
			w.focusIndex = 0
		})
	}
	<-w.stopCh
}

// Quit() — это принудительный выход из приложения.
func (w *window) Quit() {
	close(w.stopCh)
}

// Run() возвращает канал сигнализации о выходе.
func (w *window) OnQuit() <-chan struct{} {
	return w.stopCh
}

// IsRunned() возращает true, если приложение уже запущено. Иначе возвращает false.
func (w *window) IsRunned() bool {
	return w.runned
}

const taskBufSize = 32

// NewWindow() создаёт объект приложения без логирования.
func NewWindow() Window {
	wnd := &window{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[keyboard.Key]func()),
		work: make(chan *task, taskBufSize), focusIndex: -1, focusChange: true,
	}
	if DEBUG {
		f, err := os.Create(fmt.Sprintf("debug_log_%d", time.Now().UnixMilli()))
		if err != nil {
			log.Fatal(err)
		}
		wnd.log = f
	}
	if runtime.GOOS == "windows" {
		enableANSI()
	}
	currentWindow = wnd
	go wnd.runWorker()
	return wnd
}

// RegisterKeyHandler() добавляет обработчик нажатия клавиши.
func (w *window) RegisterKeyHandler(key keyboard.Key, h func()) {
	w.keyHandlers[key] = h
}

// Do() запускает функцию f в потоке GUI, что спасает от data racing при изменении виджетов.
func (w *window) Do(f func()) {
	w.work <- &task{f: f}
}

// Do() запускает функцию f в потоке GUI и ждёт завершения.
func (w *window) DoAndWait(f func()) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
	}
	w.work <- tsk
	<-tsk.done
}

func (w *window) doWithMessage(f func(), msg string) {
	w.work <- &task{
		f:   f,
		msg: msg,
	}
}

func (w *window) doWithMessageAndWait(f func(), msg string) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
		msg:  msg,
	}
	w.work <- tsk
	<-tsk.done
}

func (w *window) runWorker() {
	w.LogInfo("Воркер запущен...")
	for {
		select {
		case <-w.stopCh:
			w.runned = false
			keyboard.Close()
			fmt.Print("\033[?25l")
			fmt.Fprint(w.f, "\033[2J\033[H\033[?25h")
			w.LogInfo("Воркер остановлен...")
			return
		case tsk := <-w.work:
			if tsk.msg != "" {
				w.LogInfo("Принята задача: '%s'", tsk.msg)
			} else {
				w.LogInfo("Принята задача")
			}
			tsk.f()
			if tsk.done != nil {
				close(tsk.done)
			}
			if tsk.msg != "" {
				w.LogInfo("Завершена задача: '%s'", tsk.msg)
			} else {
				w.LogInfo("Завершена задача")
			}
		}
	}
}

func (wnd *window) Width() int {
	w, _, err := term.GetSize(wnd.f.Fd())
	if err != nil {
		wnd.LogFatal("tui: get window size error")
	}
	return w
}

func (wnd *window) Height() int {
	_, h, err := term.GetSize(wnd.f.Fd())
	if err != nil {
		wnd.LogFatal("tui: get window size error")
	}
	return h
}

func (wnd *window) DisableFocusChange() {
	wnd.focusChange = false
}

func CurrentWindow() Window {
	return currentWindow
}
