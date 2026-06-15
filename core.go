package tui

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/eiannone/keyboard"
	"golang.org/x/term"
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

type MouseEventHandler func(*MouseEvent)

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
	comp          []Widget
	compF         []Focusable
	f             *os.File
	focusIndex    int
	stopCh        chan struct{}
	keyHandlers   map[keyboard.Key]func()
	currentPos    pos
	posWidgets    []pos
	posWidgetsF   []pos
	log           *os.File
	runned        bool
	work          chan *task
	focusChange   bool
	stdout        *os.File
	stderr        *os.File
	oldMode       *term.State
	mouseHandlers []MouseEventHandler
}

// Widgets() возвращает список компонентов, добавленных в приложение.
func (wnd *window) Widgets() []Widget {
	return wnd.comp
}

// Redraw() перерисовывает все компоненты. Он потокобезопасен.
// Важно: такая перерисовка вызывает мерцание.
func (wnd *window) Redraw() {
	wnd.doWithMessage(func() {
		fmt.Fprint(wnd.f, "\033[2J\033[H")

		for idx, c := range wnd.comp {
			if c != nil {
				if idx >= len(wnd.posWidgets) {
					wnd.LogFatal("позиция для виджета %d не найдена", idx)
				}
				pos := wnd.posWidgets[idx]
				fmt.Fprintf(wnd.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
				fmt.Fprint(wnd.f, c.InnerText())
			}
		}
	}, "redraw all")
}

func (wnd *window) index() {
	wnd.compF = []Focusable{}
	wnd.posWidgets = []pos{}
	wnd.posWidgetsF = []pos{}
	wnd.currentPos = pos{0, 0}
	for idx, c := range wnd.comp {
		if c != nil {
			if len(stripansi.Strip(c.InnerText())) > c.MaxLength() {
				wnd.LogFatal("Ошибка индексации: MaxLength() не верен.")
			}
			focusable := false
			if f, ok := c.(Focusable); ok {
				wnd.compF = append(wnd.compF, f)
				focusable = true
			}
			c.SetIndex(idx)
			switch c.DisplayMode() {
			case DisplayInline:
				if wnd.currentPos.Col+c.MaxLength() >= wnd.Width() {
					wnd.currentPos.Col = 0
					wnd.currentPos.Line++
				}
				wnd.posWidgets = append(wnd.posWidgets, wnd.currentPos)
				if focusable {
					wnd.posWidgetsF = append(wnd.posWidgetsF, wnd.currentPos)
				}

				wnd.currentPos.Col += c.MaxLength()
			case DisplayBlock:
				wnd.currentPos.Col = 0
				wnd.currentPos.Line++

				wnd.posWidgets = append(wnd.posWidgets, wnd.currentPos)
				if focusable {
					wnd.posWidgetsF = append(wnd.posWidgetsF, wnd.currentPos)
				}

				wnd.currentPos.Col = 0
				wnd.currentPos.Line++
			case DisplayNewLine:
				wnd.posWidgets = append(wnd.posWidgets, wnd.currentPos)
				if focusable {
					wnd.posWidgetsF = append(wnd.posWidgetsF, wnd.currentPos)
				}

				wnd.currentPos.Col = 0
				wnd.currentPos.Line++
			}
		}
	}
}

// RedrawWidget() перерисовывает конкретный компонент. Потокобезопасен.
// index - это номер компонента, который нужно перерисовать.
func (wnd *window) RedrawWidget(index int) {
	wnd.doWithMessage(func() {
		wnd.LogInfo("RedrawWidget %v", wnd.posWidgets)
		pos := wnd.posWidgets[index]
		fmt.Fprintf(wnd.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
		wnd.LogInfo("%v %d", pos, index)
		fmt.Print(wnd.comp[index].InnerText() + strings.Repeat(" ", wnd.comp[index].MaxLength()-len(stripansi.Strip(wnd.comp[index].InnerText()))))
	}, "redraw widget")
}

// AddWidgets() добавляет компонент в приложение. Потокобезопасен.
func (wnd *window) AddWidgets(c ...Widget) {
	wnd.doWithMessageAndWait(func() {
		wnd.comp = append(wnd.comp, c...)
	}, "add widget")
}

// Clear() очищает список компонентов приложения без перерисовки. Потокобезопасен.
func (wnd *window) Clear() {
	wnd.doWithMessageAndWait(func() {
		wnd.comp = []Widget{}
		wnd.compF = []Focusable{}
		wnd.posWidgets = []pos{}
	}, "clear")
}

// Run() - это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
func (wnd *window) Run() {
	defer func() {
		if DEBUG {
			wnd.log.Close()
		}
		if err := recover(); err != nil {
			wnd.LogFatal("Произошла panic: %v", err)
		}
	}()
	if !term.IsTerminal(int(wnd.f.Fd())) {
		fmt.Fprintln(wnd.f, "Приложение запущено не в терминале. Выход...")
		time.Sleep(time.Second * 3)
		wnd.LogFatal("tui: stdout is not terminal")
	}
	wnd.stdout = os.Stdout
	wnd.stderr = os.Stderr
	os.Stdout, os.Stderr = wnd.log, wnd.log

	wnd.enableRawMode()
	defer wnd.restoreTerminalMode()

	wnd.index()

	fmt.Fprint(wnd.f, "\033[?25l")

	wnd.Redraw()

	go wnd.startStopSignalCatcher()
	go wnd.startScreenResizeChecker()
	go wnd.startMouseCatcher()
	go wnd.startKeyCatcher()

	if len(wnd.compF) != 0 {
		wnd.Do(func() {
			wnd.compF[0].OnFocus()
			wnd.focusIndex = 0
		})
	}

	wnd.runned = true
	<-wnd.stopCh
	wnd.restoreOut()
}

func (wnd *window) restoreOut() {
	os.Stdout = wnd.stdout
	os.Stderr = wnd.stderr
}

// Quit() — это принудительный выход из приложения.
func (wnd *window) Quit() {
	close(wnd.stopCh)
}

// Run() возвращает канал сигнализации о выходе.
func (wnd *window) OnQuit() <-chan struct{} {
	return wnd.stopCh
}

// IsRunned() возращает true, если приложение уже запущено. Иначе возвращает false.
func (wnd *window) IsRunned() bool {
	return wnd.runned
}

const taskBufSize = 32

// NewWindow() создаёт объект приложения.
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
	enableANSI()
	currentWindow = wnd
	go wnd.runWorker()
	return wnd
}

// RegisterKeyHandler() добавляет обработчик нажатия клавиши.
func (wnd *window) RegisterKeyHandler(key keyboard.Key, h func()) {
	wnd.keyHandlers[key] = h
}

// Do() запускает функцию f в потоке GUI, что спасает от data racing при изменении виджетов.
func (wnd *window) Do(f func()) {
	wnd.work <- &task{f: f}
}

// Do() запускает функцию f в потоке GUI и ждёт завершения.
func (wnd *window) DoAndWait(f func()) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
	}
	wnd.work <- tsk
	<-tsk.done
}

func (wnd *window) doWithMessage(f func(), msg string) {
	wnd.work <- &task{
		f:   f,
		msg: msg,
	}
}

func (wnd *window) doWithMessageAndWait(f func(), msg string) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
		msg:  msg,
	}
	wnd.work <- tsk
	<-tsk.done
}

func (wnd *window) runWorker() {
	wnd.LogInfo("Воркер запущен...")
	for {
		select {
		case <-wnd.stopCh:
			wnd.runned = false
			keyboard.Close()
			fmt.Print("\033[?25l")
			fmt.Fprint(wnd.f, "\033[2J\033[H\033[?25h")
			wnd.LogInfo("Воркер остановлен...")
			return
		case tsk := <-wnd.work:
			if tsk.msg != "" {
				wnd.LogInfo("Принята задача: '%s'", tsk.msg)
			} else {
				wnd.LogInfo("Принята задача")
			}
			tsk.f()
			if tsk.done != nil {
				close(tsk.done)
			}
			if tsk.msg != "" {
				wnd.LogInfo("Завершена задача: '%s'", tsk.msg)
			} else {
				wnd.LogInfo("Завершена задача")
			}
		}
	}
}

func (wnd *window) Width() int {
	width, _, err := term.GetSize(int(wnd.f.Fd()))
	if err != nil {
		wnd.LogFatal("tui: get window size error")
	}
	return width
}

func (wnd *window) Height() int {
	_, height, err := term.GetSize(int(wnd.f.Fd()))
	if err != nil {
		wnd.LogFatal("tui: get window size error")
	}
	return height
}

func (wnd *window) DisableFocusChange() {
	wnd.Do(func() {
		wnd.focusChange = false
	})
}

func (wnd *window) enableRawMode() {
	old, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		wnd.LogFatal("Ошибка перехода в RAW режим:")
	}
	wnd.oldMode = old
	fmt.Print("\033[?1000h\033[?1006h")
}

func (wnd *window) restoreTerminalMode() {
	if wnd.oldMode != nil {
		fmt.Print("\033[?1006l\033[?1000l")
		term.Restore(int(os.Stdin.Fd()), wnd.oldMode)
	}
}

type MouseEvent struct {
	Button  int // 0=левый, 1=средний, 2=правый, 128=отпущена
	X, Y    int // в символах
	IsPress bool
	IsDrag  bool
}

func parseMouseEvent(input string) (*MouseEvent, error) {
	if !strings.HasPrefix(input, "\x1b[<") {
		return nil, fmt.Errorf("не SGR последовательность")
	}
	rest := strings.TrimPrefix(input, "\x1b[<")
	rest = strings.TrimSuffix(rest, "m")

	parts := strings.Split(rest, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("неверный формат: %v", parts)
	}
	btn, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	isRelease := (btn & 0x80) != 0
	isDrag := (btn & 0x40) != 0
	button := btn & 0x03

	return &MouseEvent{
		Button:  button,
		X:       x - 1,
		Y:       y - 1,
		IsPress: !isRelease && !isDrag,
		IsDrag:  isDrag,
	}, nil
}

func (wnd *window) startMouseCatcher() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-wnd.stopCh:
			return
		default:
			n, err := os.Stdin.Read(buf)
			if err != nil {
				return
			}
			input := string(buf[:n])

			if n == 1 && buf[0] == 3 {
				close(wnd.stopCh)
				return
			}

			ev, err := parseMouseEvent(input)
			if err != nil {
				continue
			}

			wnd.handleMouseEvent(ev)
		}
	}
}

func (wnd *window) startStopSignalCatcher() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	select {
	case <-wnd.stopCh:
	default:
		close(wnd.stopCh)
	}
}

func (wnd *window) handleMouseEvent(ev *MouseEvent) {
	for i, pos := range wnd.posWidgetsF {
		if ev.Y == pos.Line && ev.X > pos.Col && ev.Y <= pos.Col+wnd.compF[i].MaxLength() {
			// Пользователь нажал на виджет
			if cl, ok := wnd.compF[i].(Clickable); ok {
				wnd.Do(cl.OnClick)
			}
		}
	}
	for _, h := range wnd.mouseHandlers {
		h(ev)
	}
}

func (wnd *window) RegisterClickHandler(h func(ev *MouseEvent)) {
	wnd.Do(func() {
		wnd.mouseHandlers = append(wnd.mouseHandlers, h)
	})
}

func (wnd *window) startKeyCatcher() {
	keys, err := keyboard.GetKeys(2)
	if err != nil {
		wnd.LogFatal("tui: keyboard errror")
	}
	if wnd.focusChange && len(wnd.compF) != 0 {
		wnd.RegisterKeyHandler(keyboard.KeyArrowLeft, func() {
			if wnd.focusIndex <= 0 {
				return
			}
			wnd.compF[wnd.focusIndex].OnBlur()
			wnd.focusIndex--
			wnd.compF[wnd.focusIndex].OnFocus()
		})

		wnd.RegisterKeyHandler(keyboard.KeyArrowRight, func() {
			if wnd.focusIndex > len(wnd.compF)-2 {
				return
			}
			if wnd.focusIndex == -1 {
				wnd.compF[0].OnFocus()
				wnd.focusIndex = 0
				return
			}
			wnd.compF[wnd.focusIndex].OnBlur()
			wnd.focusIndex++
			wnd.compF[wnd.focusIndex].OnFocus()
		})
	}

	wnd.RegisterKeyHandler(keyboard.KeyEnter, func() {
		if wnd.focusIndex != -1 {
			if cl, ok := wnd.compF[wnd.focusIndex].(Clickable); ok {
				wnd.Do(cl.OnClick)
			}
		}
	})
	for {
		select {
		case ev := <-keys:
			if ev.Key == keyboard.KeyCtrlC {
				close(wnd.stopCh)
			}
			wnd.doWithMessageAndWait(func() {
				if v, ok := wnd.keyHandlers[ev.Key]; ok {
					wnd.Do(v)
				} else if ev.Err != nil {
					if ev.Err.Error() == "operation canceled" {
						close(wnd.stopCh)
						return
					}
					wnd.LogFatal("tui: keyboard error")
				}
			}, "key handler")
		case <-wnd.stopCh:
			keyboard.Close()
		}
	}
}

func CurrentWindow() Window {
	return currentWindow
}
