package tui

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/romanSPB15/tui-compose/v3/cell"
	"github.com/romanSPB15/tui-compose/v3/input"
	termL "github.com/romanSPB15/tui-compose/v3/term"
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
	MouseEventHandler    func(*input.MouseEvent)
	KeyboardEventHandler func(*input.KeyboardEvent)
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
	buf              [][]cell.Cell
	focusableWidgets []Focusable
	overlay          Widget
	displayOverlay   bool
}

func (wnd *window) indexClickable(wgt Widget, offset Pos) {
	if c, ok := wgt.(Container); ok {
		for i, child := range c.Child() {
			childOffset := Pos{
				Line: offset.Line + c.Pos(i).Line,
				Col:  offset.Col + c.Pos(i).Col,
			}
			wnd.indexClickable(child, childOffset)
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

}

func (wnd *window) indexFocusable(wgt Widget, offset Pos) {
	if c, ok := wgt.(Container); ok {
		for i, child := range c.Child() {
			childOffset := Pos{
				Line: offset.Line + c.Pos(i).Line,
				Col:  offset.Col + c.Pos(i).Col,
			}
			wnd.indexFocusable(child, childOffset)
		}
		return
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

	wnd.indexClickable(wnd.overlay, Pos{0, 0})
	wnd.indexFocusable(wnd.overlay, Pos{0, 0})

	wnd.indexClickable(wnd.content, Pos{0, 0})
	wnd.indexFocusable(wnd.content, Pos{0, 0})
}

func (wnd *window) draw(wgt Widget, pos Pos, buf [][]cell.Cell) {
	if wgt == nil {
		return
	}
	if c, ok := wgt.(Container); ok {
		for i, ch := range c.Child() {
			wnd.draw(ch, Pos{Line: pos.Line + c.Pos(i).Line, Col: pos.Col + c.Pos(i).Col}, buf)
		}

	} else {

		txt := wgt.InnerText()

		if txt == "" {
			return
		}

		txt = strings.ReplaceAll(txt, "\r\n", "\n")
		widgetLines := strings.Split(txt, "\n")

		for i, line := range widgetLines {
			if i >= wgt.MaxHeight() {
				return
			}
			if pos.Line+i >= wnd.Height() {
				return
			}
			w := wnd.Width() - pos.Col

			if w < 0 {
				continue
			}

			c := cell.Parse(line)

			copy(buf[pos.Line+i][pos.Col:], c)
		}
	}
}

func (wnd *window) render() [][]cell.Cell {
	h := wnd.Height()
	w := wnd.Width()
	buf := make([][]cell.Cell, h)

	for i := range buf {
		buf[i] = make([]cell.Cell, w)
		for j := range buf[i] {
			buf[i][j] = cell.Cell{Char: ' ', ANSI: nil}
		}
	}

	if wnd.content == nil {
		return buf
	}

	wnd.draw(wnd.content, Pos{0, 0}, buf)
	if wnd.displayOverlay {
		wnd.draw(wnd.overlay, Pos{0, 0}, buf)
	}
	return buf
}

func (wnd *window) Redraw() {
	if wnd.content == nil || !wnd.runned {
		return
	}
	newBuf := wnd.render()
	oldBuf := wnd.buf

	h := wnd.Height()
	w := wnd.Width()

	if len(oldBuf) < h {
		newOld := make([][]cell.Cell, h)
		copy(newOld, oldBuf)
		for i := len(oldBuf); i < h; i++ {
			newOld[i] = make([]cell.Cell, w)
			for j := range newOld[i] {
				newOld[i][j] = cell.Cell{Char: ' ', ANSI: nil}
			}
		}
		oldBuf = newOld
	}

	//fmt.Fprintf(wnd.f, "\033[%d;1H%s", row+1, builder.String())

	var res strings.Builder

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if !cellsEqual(newBuf[y][x], oldBuf[y][x]) {
				fmt.Println(x, y, newBuf[y][x])
				fmt.Fprintf(&res, "\033[%d;%dH", y+1, x+1)

				res.WriteString("\033[0m")
				if len(newBuf[y][x].ANSI) > 0 {
					for _, a := range newBuf[y][x].ANSI {
						res.WriteString("\033[")
						res.WriteString(a)
						res.WriteRune('m')
					}
				}
				res.WriteString(string(newBuf[y][x].Char))

				res.WriteString("\033[0m")
			}
		}
	}

	fmt.Fprint(wnd.f, res.String())

	wnd.buf = newBuf
}

func cellsEqual(a, b cell.Cell) bool {
	if a.Char != b.Char || !slices.Equal(a.ANSI, b.ANSI) {
		return false
	}
	return true
}

func (wnd *window) SetOverlay(wgt Widget) {
	wnd.Do(func() {
		wnd.overlay = wgt
		if wnd.displayOverlay {
			wnd.Redraw()
		}
		wnd.index()
	})
}

func (wnd *window) ShowOverlay() {
	wnd.Do(func() {
		if !wnd.displayOverlay {
			wnd.displayOverlay = true
			wnd.Redraw()
		}
	})
}

func (wnd *window) HideOverlay() {
	wnd.Do(func() {
		if wnd.displayOverlay {
			wnd.displayOverlay = false
			wnd.Redraw()
		}
	})
}

func (wnd *window) Run() {
	defer func() {
		if DEBUG {
			wnd.log.Close()
		}
		if err := recover(); err != nil {
			wnd.LogFatal("Произошла паника: %v", err)
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

	termL.MakeRaw()
	defer termL.Restore()

	fmt.Fprint(wnd.f, "\033[2J")

	fmt.Fprint(wnd.f, "\033[?25l")

	go wnd.startStopSignalCatcher()
	go wnd.startScreenResizeChecker()
	go wnd.startInputCatcher()

	wnd.runned = true

	wnd.Redraw()

	<-wnd.stopCh
	wnd.runned = false

	wnd.restoreOut()
	fmt.Fprint(wnd.f, "\033[0m\033[2J\033[H\033[?25h")
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
	termL.EnableANSIWindows()
	currentWindow = wnd
	go wnd.runWorker()
	return wnd
}

func (wnd *window) RegisterKeyHandler(keh KeyboardEventHandler) {
	wnd.keyHandlers = append(wnd.keyHandlers, keh)
}

func (wnd *window) Do(f func()) {
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- &task{f: f}:
	}

}

func (wnd *window) DoAndWait(f func()) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
	}
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- tsk:
		<-tsk.done
	}
}

func (wnd *window) doWithMessage(f func(), msg string) {
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- &task{
		f:   f,
		msg: msg,
	}:
	}
}

func (wnd *window) doWithMessageAndWait(f func(), msg string) {
	tsk := &task{
		f:    f,
		done: make(chan struct{}),
		msg:  msg,
	}
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- tsk:
		<-tsk.done
	}
}

func (wnd *window) runWorker() {
	wnd.LogInfo("Воркер запущен...")
	for {
		select {
		case <-wnd.stopCh:
			close(wnd.work)
			wnd.LogInfo("Воркер остановлен")
			return
		case tsk := <-wnd.work:
			if tsk.msg != "" {
				wnd.LogInfo("Принята задача: '%s'", tsk.msg)
			} else {
				wnd.LogInfo("Принята задача")
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						if tsk.msg != "" {
							wnd.LogInfo("Задача '%s' вызвала панику: %v", tsk.msg, r)
						} else {
							wnd.LogInfo("Задача вызвала панику: %v", r)
						}
					}
				}()
				tsk.f()
			}()
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
	w, _ := termL.SizeFd(wnd.stdout.Fd())
	return w
}

func (wnd *window) Height() int {
	_, h := termL.SizeFd(wnd.stdout.Fd())
	return h
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

func (wnd *window) handleMouseEvent(ev *input.MouseEvent) {
	if wnd.cl != nil {
		for _, cl := range wnd.cl {
			if ev.Pos.Y >= cl.p.Line && ev.Pos.Y < cl.p.Line+cl.MaxHeight() && ev.Pos.X >= cl.p.Col && ev.Pos.X < cl.p.Col+cl.MaxWidth() {
				// Пользователь нажал на этот виджет
				wnd.doWithMessage(cl.OnClick, "click handler")

			}
		}
	}
	if wnd.clAt != nil {
		for _, clAt := range wnd.clAt {
			if ev.Pos.Y >= clAt.p.Line && ev.Pos.Y < clAt.p.Line+clAt.MaxHeight() &&
				ev.Pos.X >= clAt.p.Col && ev.Pos.X < clAt.p.Col+clAt.MaxWidth() {
				relX := ev.Pos.X - clAt.p.Col
				relY := ev.Pos.Y - clAt.p.Line
				wnd.doWithMessage(func() {
					clAt.OnClickAt(relX, relY)
				}, "clickAt handler")
			}
		}
	}
	for _, h := range wnd.mouseHandlers {
		wnd.doWithMessage(func() {
			h(ev)
		}, "mouse handler")
	}
}

func (wnd *window) RegisterClickHandler(h func(ev *input.MouseEvent)) {
	wnd.Do(func() {
		wnd.mouseHandlers = append(wnd.mouseHandlers, h)
	})
}

func (wnd *window) CopyToClipboard(text string) {
	termL.CopyToClipboard(text)
}

func (wnd *window) startInputCatcher() {
	wnd.RegisterKeyHandler(func(ke *input.KeyboardEvent) {
		if wnd.focusIndex != -1 {
			if te, ok := wnd.focusableWidgets[wnd.focusIndex].(KeyReceiver); ok {
				te.OnKeyPress(ke)
			}
		}
		if !wnd.focusChange {
			return
		}
		switch ke.Key {
		case input.KeyTab:
			wnd.NextFocus()
		case input.KeyShiftTab:
			wnd.BeforeFocus()
		}
	})

	mouse, keyboard := input.Start(1)
	for {
		select {
		case <-wnd.stopCh:
			input.Stop()
			return
		case ev := <-keyboard:
			wnd.doWithMessage(func() {
				for _, h := range wnd.keyHandlers {
					wnd.doWithMessage(func() {
						h(ev)
					}, "keyboard handler")
				}
			}, "key handler")
		case ev := <-mouse:
			wnd.handleMouseEvent(ev)
		}
	}
}

func (wnd *window) SetContent(w Widget) {
	wnd.content = w
	wnd.index() // перестраиваем список кликабельных
	wnd.focusIndex = -1
}

func (wnd *window) SetTitle(title string) {
	fmt.Fprintf(wnd.f, "\033]0;%s\033\\", title)
}

func (wnd *window) Focus() FocusManager {
	return wnd
}

func CurrentWindow() Window {
	return currentWindow
}
