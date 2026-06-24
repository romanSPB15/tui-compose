package tui

import (
	"bytes"
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

type (
	MouseEventHandler    func(*MouseEvent)
	KeyboardEventHandler func(*KeyboardEvent)
)

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
	cl               []clickableWidgetWithPos
	clAt             []clickableAtWidgetWithPos
	f                *os.File
	focusIndex       int
	stopCh           chan struct{}
	keyHandlers      []KeyboardEventHandler
	log              *os.File
	runned           bool
	work             chan *task
	focusChange      bool
	stdout           *os.File
	stderr           *os.File
	oldMode          *term.State
	mouseHandlers    []MouseEventHandler
	content          Widget
	buf              []string
	focusableWidgets []Focusable
}

func (wnd *window) drawContainer(buf io.Writer, p Pos, c Container) {
	for i, w := range c.Child() {
		pos := c.Pos(i)
		if c, ok := w.(Container); ok {
			wnd.drawContainer(buf, pos, c)
		} else {
			fmt.Fprintf(buf, "\033[%d;%dH\033[0m", pos.Line+1+p.Line, pos.Col+1+p.Col)
			it := w.InnerText()
			r := []rune(it)
			if len(r) > w.MaxWidth() {
				fmt.Fprint(buf, string(r[:w.MaxWidth()]))
			} else {
				fmt.Fprint(buf, it+strings.Repeat(" ", w.MaxWidth()-len(r)))
			}

		}
	}
}

func (wnd *window) indexClickableAndFocusable(wgt Widget, offset Pos) {
	if c, ok := wgt.(Container); ok {
		for i, child := range c.Child() {
			childOffset := Pos{
				Line: offset.Line + c.Pos(i).Line,
				Col:  offset.Col + c.Pos(i).Col,
			}
			wnd.indexClickableAndFocusable(child, childOffset)
		}
		return
	}

	if cl, ok := wgt.(Clickable); ok {
		if _, ok := wgt.(ClickableAt); !ok {
			wnd.cl = append(wnd.cl, clickableWidgetWithPos{
				Clickable: cl,
				p:         offset,
			})
		}
	}

	if clAt, ok := wgt.(ClickableAt); ok {
		wnd.clAt = append(wnd.clAt, clickableAtWidgetWithPos{
			ClickableAt: clAt,
			p:           offset,
		})
	}

	if foc, ok := wgt.(Focusable); ok {
		wnd.focusableWidgets = append(wnd.focusableWidgets, foc)
	}
}

func (wnd *window) index() {
	if wnd.content == nil {
		return
	}
	wnd.cl = nil
	wnd.clAt = nil

	wnd.indexClickableAndFocusable(wnd.content, Pos{0, 0})

	wnd.indexRec(wnd.content, Pos{0, 0})

}

func (wnd *window) indexRec(w Widget, offset Pos) {
	if w == nil {
		return
	}

	if cl, ok := w.(Clickable); ok {
		wnd.cl = append(wnd.cl, clickableWidgetWithPos{
			Clickable: cl,
			p:         offset,
		})
	}

	if clAt, ok := w.(ClickableAt); ok {
		wnd.clAt = append(wnd.clAt, clickableAtWidgetWithPos{
			ClickableAt: clAt,
			p:           offset,
		})
	}

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

func (wnd *window) Redraw() {
	wnd.doWithMessage(func() {
		if wnd.content == nil || wnd.runned == false {
			return
		}
		buf := &bytes.Buffer{}
		fmt.Fprint(buf, "\033[H")

		if c, ok := wnd.content.(Container); ok {
			wnd.drawContainer(buf, Pos{0, 0}, c)
		} else {
			fmt.Fprint(buf, wnd.content.InnerText())
		}
		new := strings.Split(buf.String(), "\r\n")

		changed := []int{}
		for i := range new {
			if i >= len(wnd.buf) || wnd.buf[i] != new[i] {
				changed = append(changed, i)
			}
		}
		h := wnd.Height()

		switch {
		case len(changed) == len(new):
			if len(new) > h {
				for i := range h {
					fmt.Fprint(wnd.f, new[i], "\033[K")
				}
			} else {
				io.Copy(wnd.f, buf)
			}
			wnd.buf = new
		case len(changed) == 0:
			return
		default:
			for _, idx := range changed {
				if idx > h {
					return
				}
				fmt.Fprintf(wnd.f, "\033[%d;1H%s\033[K", idx+1, new[idx])
			}
			wnd.buf = new
		}
	}, "redraw all")
}

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

	fmt.Fprint(wnd.f, "\033[2J")

	fmt.Fprint(wnd.f, "\033[?25l")

	wnd.Redraw()

	go wnd.startStopSignalCatcher()
	go wnd.startScreenResizeChecker()
	go wnd.startInputCatcher()

	wnd.runned = true
	<-wnd.stopCh
	wnd.restoreOut()
	fmt.Fprint(wnd.f, "\033[2J\033[H\033[?25h")
}

func (wnd *window) restoreOut() {
	os.Stdout = wnd.stdout
	os.Stderr = wnd.stderr
}

func (wnd *window) Quit() {
	close(wnd.stopCh)
}

func (wnd *window) OnQuit() <-chan struct{} {
	return wnd.stopCh
}

func (wnd *window) IsRunned() bool {
	return wnd.runned
}

const taskBufSize = 32

func NewWindow() Window {
	wnd := &window{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: []KeyboardEventHandler{},
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

func (wnd *window) RegisterKeyHandler(keh KeyboardEventHandler) {
	wnd.keyHandlers = append(wnd.keyHandlers, keh)
}

func (wnd *window) Do(f func()) {
	wnd.work <- &task{f: f}
}

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
		return
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
	copyToClipboard(text)
}

func (wnd *window) startInputCatcher() {
	if wnd.focusChange && len(wnd.focusableWidgets) != 0 {
		wnd.RegisterKeyHandler(func(ke *KeyboardEvent) {
			switch ke.Key {
			case KeyTab:
				if wnd.focusIndex > len(wnd.focusableWidgets)-2 {
					return
				}
				if wnd.focusIndex == -1 {
					wnd.focusableWidgets[0].OnFocus()
					wnd.focusIndex = 0
					return
				}
				wnd.focusableWidgets[wnd.focusIndex].OnBlur()
				wnd.focusIndex++
				wnd.focusableWidgets[wnd.focusIndex].OnFocus()
			case KeyShiftTab:
				if wnd.focusIndex <= 0 {
					return
				}
				wnd.focusableWidgets[wnd.focusIndex].OnBlur()
				wnd.focusIndex--
				wnd.focusableWidgets[wnd.focusIndex].OnFocus()
			case KeyEnter:
				if wnd.focusIndex != -1 {
					if cl, ok := wnd.focusableWidgets[wnd.focusIndex].(Clickable); ok {
						wnd.Do(cl.OnClick)
					}
				}
			}
		})
		wnd.RegisterKeyHandler(func(ke *KeyboardEvent) {
			if wnd.focusIndex != -1 {
				if te, ok := wnd.focusableWidgets[wnd.focusIndex].(TextInput); ok {
					te.OnKeyPress(ke)
				}
			}
		})
	}

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

func (wnd *window) SetContent(w Widget) {
	wnd.content = w
	wnd.index() // перестраиваем список кликабельных
	wnd.focusIndex = -1
}

func (wnd *window) SetTitle(title string) {
	fmt.Printf("\033]0;%s\033\\", title)
}

func CurrentWindow() Window {
	return currentWindow
}
