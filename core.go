package tui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

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

type Pos struct {
	Line int
	Col  int
}

type task struct {
	done chan struct{}
	f    func()
	msg  string
}

type clickableWidgetWithPos struct {
	Clickable
	p Pos
}

type clickableAtWidgetWithPos struct {
	ClickableAt
	p Pos
}

var currentWindow *window

type window struct {
	cl            []clickableWidgetWithPos
	clAt          []clickableAtWidgetWithPos
	f             *os.File
	focusIndex    int
	stopCh        chan struct{}
	keyHandlers   map[Key]func()
	log           *os.File
	runned        bool
	work          chan *task
	focusChange   bool
	stdout        *os.File
	stderr        *os.File
	oldMode       *term.State
	mouseHandlers []MouseEventHandler
	content       Widget
}

func (wnd *window) drawContainer(buf io.Writer, p Pos, c Container) {
	for i, w := range c.Child() {
		pos := c.Pos(i)
		if c, ok := w.(Container); ok {
			wnd.drawContainer(buf, pos, c)
		} else {
			fmt.Fprintf(buf, "\033[%d;%dH", pos.Line+1+p.Line, pos.Col+1+p.Col)
			fmt.Fprint(buf, w.InnerText())
		}
	}
}

func (wnd *window) indexClickable() {
	wnd.cl = nil
	if wnd.content == nil {
		return
	}
	wnd.indexRec(wnd.content, Pos{0, 0})

	fmt.Println(wnd.cl)
}

func (wnd *window) indexRec(w Widget, offset Pos) {
	if w == nil {
		return
	}

	// Проверяем Clickable
	if cl, ok := w.(Clickable); ok {
		wnd.cl = append(wnd.cl, clickableWidgetWithPos{
			Clickable: cl,
			p:         offset,
		})
	}

	// Проверяем ClickableAt (независимо от Clickable)
	if clAt, ok := w.(ClickableAt); ok {
		wnd.clAt = append(wnd.clAt, clickableAtWidgetWithPos{
			ClickableAt: clAt,
			p:           offset,
		})
	}

	// Если это контейнер, обходим детей
	if c, ok := w.(Container); ok {
		for i, child := range c.Child() {
			childPos := c.Pos(i)
			newOffset := Pos{
				Line: offset.Line + childPos.Line,
				Col:  offset.Col + childPos.Col,
			}
			wnd.indexRec(child, newOffset)
		}
	}
}

// Redraw() перерисовывает все компоненты. Он потокобезопасен.
func (wnd *window) Redraw() {
	if wnd.content == nil {
		return
	}
	wnd.doWithMessage(func() {
		buf := &bytes.Buffer{}
		fmt.Fprint(buf, "\033[2J\033[H")

		if c, ok := wnd.content.(Container); ok {
			wnd.drawContainer(buf, Pos{0, 0}, c)
		} else {
			fmt.Fprint(buf, wnd.content.InnerText())
		}
		io.Copy(wnd.f, buf)
	}, "redraw all")
}

// RedrawWidget() перерисовывает конкретный компонент. Потокобезопасен.
// index - это номер компонента, который нужно перерисовать.
func (wnd *window) RedrawWidget(index int) {
	wnd.doWithMessage(func() {
		// TODO: сделать RedrawWidget()
		/*
			wnd.LogInfo("RedrawWidget %v", wnd.posWidgets)
			pos := wnd.posWidgets[index]
			fmt.Fprintf(wnd.f, "\033[%d;%dH", pos.Line+1, pos.Col+1)
			wnd.LogInfo("%v %d", pos, index)
			fmt.Fprint(wnd.f, wnd.comp[index].InnerText()+strings.Repeat(" ", wnd.comp[index].MaxWidth()-len(stripansi.Strip(wnd.comp[index].InnerText()))))
		*/
	}, "redraw widget")
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

	wnd.indexClickable()

	wnd.enableRawMode()
	defer wnd.restoreTerminalMode()

	fmt.Fprint(wnd.f, "\033[?25l")

	wnd.Redraw()

	go wnd.startStopSignalCatcher()
	go wnd.startScreenResizeChecker()
	go wnd.startInputCatcher()

	wnd.runned = true
	<-wnd.stopCh
	wnd.restoreOut()
	fmt.Fprint(wnd.f, "\033[2J\033[H")
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
	wnd := &window{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: make(map[Key]func()),
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
func (wnd *window) RegisterKeyHandler(key Key, h func()) {
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
			fmt.Fprint(wnd.f, "\033[?25l")
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
	fmt.Fprint(wnd.f, "\033[?1000h\033[?1006h")
}

func (wnd *window) restoreTerminalMode() {
	if wnd.oldMode != nil {
		fmt.Fprint(wnd.f, "\033[?1006l\033[?1000l")
		term.Restore(int(os.Stdin.Fd()), wnd.oldMode)
	}
}

type MouseEvent struct {
	Button  int // 0=левый, 1=средний, 2=правый, 128=отпущена
	Pos     Pos
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
		Pos:     Pos{y - 1, x - 1},
		IsPress: !isRelease && !isDrag,
		IsDrag:  isDrag,
	}, nil
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
	if wnd.cl != nil {
		for _, cl := range wnd.cl {
			if ev.Pos.Line >= cl.p.Line && ev.Pos.Line < cl.p.Line+cl.MaxHeight() && ev.Pos.Col >= cl.p.Col && ev.Pos.Col < cl.p.Col+cl.MaxWidth() {

				if ev.IsPress {
					// Пользователь нажал на этот виджет
					wnd.doWithMessage(cl.OnClick, "click handler")
				}
			}
		}
	}
	if wnd.clAt != nil {
		for _, clAt := range wnd.clAt {
			if ev.Pos.Line >= clAt.p.Line && ev.Pos.Line < clAt.p.Line+clAt.MaxHeight() &&
				ev.Pos.Col >= clAt.p.Col && ev.Pos.Col < clAt.p.Col+clAt.MaxWidth() {
				relX := ev.Pos.Col - clAt.p.Col
				relY := ev.Pos.Line - clAt.p.Line
				wnd.doWithMessage(func() {
					clAt.OnClickAt(relX, relY)
				}, "clickAt handler")
				return
			}
		}
	}
	for _, h := range wnd.mouseHandlers {
		wnd.doWithMessage(func() {
			h(ev)
		}, "mouse handler")
	}
}

func (wnd *window) RegisterClickHandler(h func(ev *MouseEvent)) {
	wnd.Do(func() {
		wnd.mouseHandlers = append(wnd.mouseHandlers, h)
	})
}

func (wnd *window) CopyToClipboard(text string) {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Fprintf(wnd.f, "\033]52;c;%s\007", encoded)
}

func (wnd *window) startInputCatcher() {
	// 	if wnd.focusChange && len(wnd.compF) != 0 {
	// 		wnd.RegisterKeyHandler(KeyArrowLeft, func() {
	// 			if wnd.focusIndex <= 0 {
	// 				return
	// 			}
	// 			wnd.compF[wnd.focusIndex].OnBlur()
	// 			wnd.focusIndex--
	// 			wnd.compF[wnd.focusIndex].OnFocus()
	// 		})

	// 		wnd.RegisterKeyHandler(KeyArrowRight, func() {
	// 			if wnd.focusIndex > len(wnd.compF)-2 {
	// 				return
	// 			}
	// 			if wnd.focusIndex == -1 {
	// 				wnd.compF[0].OnFocus()
	// 				wnd.focusIndex = 0
	// 				return
	// 			}
	// 			wnd.compF[wnd.focusIndex].OnBlur()
	// 			wnd.focusIndex++
	// 			wnd.compF[wnd.focusIndex].OnFocus()
	// 		})
	// 	}

	// 	wnd.RegisterKeyHandler(KeyEnter, func() {
	// 		if wnd.focusIndex != -1 {
	// 			if cl, ok := wnd.compF[wnd.focusIndex].(Clickable); ok {
	// 				wnd.Do(cl.OnClick)
	// 			}
	// 		}
	// 	})

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
			data := buf[:n]
			if ev, err := parseMouseEvent(string(data)); err == nil {
				wnd.handleMouseEvent(ev)
				continue
			}

			wnd.handleKeyboardInput(data)
		}
	}
}

func parseEscapeSequence(data []byte) (Key, int) {
	if len(data) < 2 || data[0] != 0x1B {
		return 0, 0
	}

	if data[1] == '[' {
		if len(data) < 3 {
			return 0, 0
		}
		switch data[2] {
		case 'A':
			return KeyArrowUp, 3
		case 'B':
			return KeyArrowDown, 3
		case 'C':
			return KeyArrowRight, 3
		case 'D':
			return KeyArrowLeft, 3
		case 'H':
			return KeyHome, 3
		case 'F':
			return KeyEnd, 3
		case '5', '6': // PgUp/PgDn (CSI 5 ~, CSI 6 ~)
			if len(data) >= 4 && data[3] == '~' {
				if data[2] == '5' {
					return KeyPgup, 4
				}
				if data[2] == '6' {
					return KeyPgdn, 4
				}
			}
			return 0, 0
		case '1', '2', '3', '4': // Home, End, Insert, Delete (CSI 1 ~, CSI 2 ~, CSI 3 ~, CSI 4 ~)
			if len(data) >= 4 && data[3] == '~' {
				switch data[2] {
				case '1':
					return KeyHome, 4
				case '2':
					return KeyInsert, 4
				case '3':
					return KeyDelete, 4
				case '4':
					return KeyEnd, 4
				}
			}
			return 0, 0
		}
	}

	// F1-F4 с ESC O P/Q/R/S (старый стиль)
	if data[1] == 'O' && len(data) >= 3 {
		switch data[2] {
		case 'P':
			return KeyF1, 3
		case 'Q':
			return KeyF2, 3
		case 'R':
			return KeyF3, 3
		case 'S':
			return KeyF4, 3
		}
	}

	// F5-F12: ESC [ 1 5 ~, ESC [ 1 7 ~ и т.д.
	if data[1] == '[' && len(data) >= 5 && data[3] == '~' {
		switch data[2] {
		case '1':
			switch data[4] {
			case '5': // ESC [ 1 5 ~ -> F5
				return KeyF5, 5
			case '7': // ESC [ 1 7 ~ -> F6
				return KeyF6, 5
			case '9': // ESC [ 1 9 ~ -> F7
				return KeyF7, 5
			}
		case '2':
			switch data[4] {
			case '0': // ESC [ 2 0 ~ -> F8
				return KeyF8, 5
			case '1': // ESC [ 2 1 ~ -> F9
				return KeyF9, 5
			case '3': // ESC [ 2 3 ~ -> F10
				return KeyF10, 5
			case '4': // ESC [ 2 4 ~ -> F11
				return KeyF11, 5
			case '5': // ESC [ 2 5 ~ -> F12
				return KeyF12, 5
			}
		}
	}

	if len(data) == 1 {
		return KeyEsc, 1
	}

	if len(data) >= 2 && data[1] >= 0x20 && data[1] <= 0x7E {
		return KeyEsc, 2
	}

	return 0, 0
}

func (wnd *window) handleKeyboardInput(data []byte) {
	if len(data) == 0 {
		return
	}

	if key, n := parseEscapeSequence(data); n > 0 {
		wnd.doWithMessageAndWait(func() {
			if handler, ok := wnd.keyHandlers[key]; ok {
				handler()
			}
		}, "key handler")
		return
	}

	if len(data) == 1 {
		b := data[0]
		var key Key
		switch b {
		case 0x03:
			close(wnd.stopCh)
			return
		case 0x0D:
			key = KeyEnter
		case 0x09:
			key = KeyTab
		case 0x7F, 0x08:
			key = KeyBackspace
		case 0x1B:
			key = KeyEsc
		default:
			if b >= 0x20 && b <= 0x7E {
				key = Key(b)
			} else if b >= 0x01 && b <= 0x1A {
				key = KeyCtrlA + Key(b) - 1
			} else {
				return
			}
		}

		wnd.doWithMessageAndWait(func() {
			if handler, ok := wnd.keyHandlers[key]; ok {
				handler()
			}
		}, "key handler")
		return
	}
}

func (wnd *window) SetContent(w Widget) {
	wnd.content = w
}

func (wnd *window) SetTitle(title string) {
	fmt.Printf("\033]0;%s\033\\", title)
}

func CurrentWindow() Window {
	return currentWindow
}
