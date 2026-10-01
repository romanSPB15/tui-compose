package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/acell/terminfo"
	"github.com/romanSPB15/tui-compose/v4/ansi"
)

var capture bool

// ColorRGB — это цвет в RGB.
type ColorRGB struct {
	R, G, B uint8
}

type (
	MouseEventHandler    func(*acell.MouseEvent)
	KeyboardEventHandler func(*acell.KeyboardEvent)
	ResizeHandler        func(width, height int)
)

// Pos - точка на экране.
type Pos struct {
	Line int
	Col  int
}

type task struct {
	done chan struct{}
	f    func()
	msg  string
}

type eventHandlerWithPos struct {
	EventHandler
	p Pos
}

type windowStats struct {
	frames      atomic.Int64
	redrawCount atomic.Int64
}

type window struct {
	focusableWidgets []eventHandlerWithPos
	wgt              []eventHandlerWithPos
	focusIndex       int
	stopCh           chan struct{}

	keyHandlers    []KeyboardEventHandler
	mouseHandlers  []MouseEventHandler
	resizeHandlers []ResizeHandler

	log         *os.File
	runned      bool
	work        chan *task
	focusChange bool
	worker      atomic.Int32

	styleFunc func(Widget)

	hovered map[EventHandler]struct{}

	maxWidgetSize Pos

	content  Widget
	initCell acell.Cell
	subBuf   [][]acell.Cell
	quitOnce sync.Once

	stats      windowStats
	lastReport time.Time

	t *acell.Terminal

	width  int
	height int

	capture bool
}

func (wnd *window) indexWgt(wgt Widget, offset Pos) {
	if c, ok := wgt.(Container); ok {
		for i, child := range c.Child() {
			childOffset := Pos{
				Line: offset.Line + c.Pos(i).Line,
				Col:  offset.Col + c.Pos(i).Col,
			}
			wnd.indexWgt(child, childOffset)
		}
	}

	if evh, ok := wgt.(EventHandler); ok {
		evh.Send(&WindowEvent{Window: wnd})
		wnd.wgt = append(wnd.wgt, eventHandlerWithPos{EventHandler: evh, p: offset})
		if isFocusable(evh) {
			wnd.focusableWidgets = append(wnd.focusableWidgets, eventHandlerWithPos{EventHandler: evh, p: offset})
		}
	}
}

func (wnd *window) applyStyles(wgt Widget) {
	if wnd.styleFunc == nil {
		return
	}
	wnd.styleFunc(wgt)
	if c, ok := wgt.(Container); ok {
		for _, child := range c.Child() {
			wnd.applyStyles(child)
		}
	}
}

func (wnd *window) Index() {
	if wnd.content == nil {
		return
	}
	wnd.focusableWidgets = nil
	wnd.wgt = nil

	wnd.indexWgt(wnd.content, Pos{0, 0})

	wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line = wnd.calcMaxWidgetSize(wnd.content, 0, 0)

	wnd.applyStyles(wnd.content)
}

func (wnd *window) draw(wgt Widget, rect [2]Pos, buf [][]acell.Cell) {
	if wgt == nil {
		return
	}
	if c, ok := wgt.(Container); ok {
		contBottom := rect[0].Line + wgt.Height()
		contRight := rect[0].Col + wgt.Width()
		if contBottom > rect[1].Line {
			contBottom = rect[1].Line
		}
		if contRight > rect[1].Col {
			contRight = rect[1].Col
		}

		for i, ch := range c.Child() {
			offLine := c.Pos(i).Line
			offCol := c.Pos(i).Col
			childRect := [2]Pos{
				{Line: rect[0].Line + offLine, Col: rect[0].Col + offCol},
				{Line: contBottom, Col: contRight},
			}
			wnd.draw(ch, childRect, buf)
		}
		return
	}

	w, h := wgt.Width(), wgt.Height()
	if w <= 0 || h <= 0 {
		return
	}

	currentH := len(wnd.subBuf)
	var currentW int
	if currentH > 0 {
		currentW = len(wnd.subBuf[0])
	}

	if currentH < h || currentW < w {
		newH := max(currentH, h)
		newW := max(currentW, w)
		newBuf := make([][]acell.Cell, newH)
		for i := range newBuf {
			newBuf[i] = make([]acell.Cell, newW)
		}
		wnd.subBuf = newBuf
	} else if currentH > h*2 || currentW > w*2 {
		newBuf := make([][]acell.Cell, h)
		for i := range newBuf {
			newBuf[i] = make([]acell.Cell, w)
		}
		wnd.subBuf = newBuf
	}

	subBuf := wnd.subBuf[:h]
	for y := range subBuf {
		if len(subBuf[y]) < w {
			subBuf[y] = append(subBuf[y], make([]acell.Cell, w-len(subBuf[y]))...)
		}
		for x := range w {
			subBuf[y][x] = wnd.initCell
		}
	}

	wgt.Render(subBuf)

	for y := 0; y < h; y++ {
		destY := y + rect[0].Line
		if destY < 0 || destY >= len(buf) {
			continue
		}
		if destY >= rect[1].Line {
			break
		}
		srcRow := subBuf[y]
		dstRow := buf[destY]

		if rect[0].Col < 0 {
			continue
		}

		copyLen := w
		if rect[0].Col+copyLen > rect[1].Col {
			copyLen = rect[1].Col - rect[0].Col
		}
		if rect[0].Col+copyLen > len(dstRow) {
			copyLen = len(dstRow) - rect[0].Col
		}
		if copyLen > 0 {
			copy(dstRow[rect[0].Col:rect[0].Col+copyLen], srcRow[:copyLen])
		}
	}
}

func (wnd *window) calcMaxWidgetSize(wgt Widget, w, h int) (int, int) {
	if wgt == nil {
		return w, h
	}
	if c, ok := wgt.(Container); ok {
		for _, v := range c.Child() {
			w2, h2 := wnd.calcMaxWidgetSize(v, w, h)
			w = max(w, w2)
			h = max(h, h2)
		}
		return w, h
	}
	return max(w, wgt.Width()), max(h, wgt.Height())
}

func (wnd *window) render() {
	h := wnd.Height()
	w := wnd.Width()
	if h <= 0 || w <= 0 {
		return
	}

	wnd.t.Fill(wnd.initCell)

	if wnd.content == nil {
		return
	}

	ww, hw := wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line
	if ww == 0 && hw == 0 {
		wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line = wnd.calcMaxWidgetSize(wnd.content, 0, 0)
	}

	wnd.draw(wnd.content, [2]Pos{{Line: 0, Col: 0}, {Line: h, Col: w}}, wnd.t.Buf)
}

func (wnd *window) Redraw() {
	if !wnd.runned {
		return
	}

	if DEBUG {
		wnd.stats.redrawCount.Add(1)
		if wnd.stats.redrawCount.Load()%64 == 1 && !wnd.isWorker() {
			wnd.LogFatal("Redraw called outside worker goroutine: data race")
		}
	}

	wnd.render()

	if wnd.capture {
		json.NewEncoder(wnd.t).Encode(wnd.t.Buf)
		return
	}

	wnd.t.Flush()
	wnd.stats.frames.Add(1)
	wnd.maybeReport()
}

func (wnd *window) maybeReport() {
	if !DEBUG {
		return
	}
	now := time.Now()
	elapsed := now.Sub(wnd.lastReport)
	if elapsed < time.Second {
		return
	}
	wnd.lastReport = now

	frames := wnd.stats.frames.Swap(0)
	if frames == 0 {
		return
	}
	wnd.LogInfo("FPS: %.0f", float64(frames)/elapsed.Seconds())
}

func getEnvInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func (wnd *window) Run() {
	defer func() {
		if err := recover(); err != nil {
			wnd.LogFatal("tui: Произошла паника: %v", err)
		}
		if wnd.t != nil {
			wnd.t.Close()
		}
		if DEBUG && wnd.log != nil {
			wnd.log.Close()
		}
	}()

	if wnd.capture {
		w := wnd.width
		h := wnd.height
		if w == 0 {
			w = getEnvInt("TUI_WIDTH", 80)
		}
		if h == 0 {
			h = getEnvInt("TUI_HEIGHT", 24)
		}
		wnd.width, wnd.height = w, h
		wnd.t = acell.NewWithTerm(&fakeRawTerminal{width: w, height: h})
		wnd.t.Buf = acell.NewBuf(w, h)
		wnd.runned = true
		wnd.Index()
		wnd.Redraw()
		return
	}

	wnd.t = acell.New(acell.TTY())

	w, h := wnd.t.Size()
	if wnd.width > 0 {
		w = wnd.width
	}
	if wnd.height > 0 {
		h = wnd.height
	}
	wnd.t.Buf = acell.NewBuf(w, h)

	go wnd.startStopSignalCatcher()
	go wnd.startInputCatcher()

	wnd.runned = true
	wnd.Redraw()
	wnd.runWorker()
	wnd.runned = false
}

func (wnd *window) Quit() {
	wnd.quitOnce.Do(func() { close(wnd.stopCh) })
}

func (wnd *window) OnQuit() <-chan struct{} { return wnd.stopCh }

func (wnd *window) IsRunned() bool {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("IsRunned called outside worker goroutine: data race")
	}
	return wnd.runned
}

// Options — настройки окна.
type Options struct {
	Terminal *acell.Terminal
}

type Option func(*Options)

func WithTerminal(t *acell.Terminal) Option {
	return func(o *Options) { o.Terminal = t }
}

const taskBufSize = 32

func NewWindow(opts ...Option) Window {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	if o.Terminal == nil {
		o.Terminal = acell.New(acell.TTY())
	}

	wnd := &window{
		t:           o.Terminal,
		stopCh:      make(chan struct{}),
		keyHandlers: []KeyboardEventHandler{},
		work:        make(chan *task, taskBufSize),
		focusIndex:  -1,
		focusChange: true,
		initCell:    acell.Cell{Char: ' '},
		lastReport:  time.Now(),
		hovered:     make(map[EventHandler]struct{}),
	}

	if DEBUG {
		f, err := os.Create(fmt.Sprintf("debug_log_%d", time.Now().UnixMilli()))
		if err != nil {
			log.Fatal(err)
		}
		wnd.log = f
	}
	return wnd
}

func (wnd *window) isWorker() bool {
	return wnd.worker.Load() == getGorID()
}

func (wnd *window) RegisterKeyHandler(keh KeyboardEventHandler) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("RegisterKeyHandler called outside worker goroutine: data race")
	}
	wnd.keyHandlers = append(wnd.keyHandlers, keh)
}

func (wnd *window) Do(f func()) {
	defer recover()
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- &task{f: f}:
	}
}

func (wnd *window) DoAndWait(f func()) {
	defer recover()
	if DEBUG && wnd.isWorker() {
		f()
		return
	}
	tsk := &task{f: f, done: make(chan struct{})}
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- tsk:
		<-tsk.done
	}
}

func (wnd *window) doWithMessage(f func(), msg string) {
	defer recover()
	select {
	case <-wnd.stopCh:
		return
	case wnd.work <- &task{f: f, msg: msg}:
	}
}

func getGorID() int32 {
	var buf [128]byte
	n := runtime.Stack(buf[:], false)
	parts := strings.SplitN(string(buf[:n]), " ", 3)
	if len(parts) < 2 {
		return -1
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1
	}
	return int32(id)
}

func (wnd *window) runWorker() {
	wnd.worker.Store(getGorID())
	wnd.LogInfo("Воркер запущен...")
	for {
		select {
		case <-wnd.stopCh:
			wnd.LogInfo("Воркер остановлен")
			return
		case tsk := <-wnd.work:
			if tsk.msg != "" {
				wnd.LogInfo("Принята задача: '%s'", tsk.msg)
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						wnd.LogInfo("Задача вызвала панику: %v", r)
					}
				}()
				tsk.f()
			}()
			if tsk.done != nil {
				close(tsk.done)
			}
		}
	}
}

func (wnd *window) Width() int {
	if wnd.width > 0 {
		return wnd.width
	}
	if wnd.t != nil {
		w, _ := wnd.t.Size()
		return w
	}
	return 80
}

func (wnd *window) Height() int {
	if wnd.height > 0 {
		return wnd.height
	}
	if wnd.t != nil {
		_, h := wnd.t.Size()
		return h
	}
	return 24
}

func (wnd *window) SetSize(w, h int) {
	wnd.Do(func() {
		wnd.width, wnd.height = w, h
		if wnd.t != nil {
			wnd.t.Buf = acell.NewBuf(w, h)
			wnd.t.Invalidate()
		}
		wnd.Redraw()
	})
}

func (wnd *window) startStopSignalCatcher() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	select {
	case <-wnd.stopCh:
		return
	default:
		wnd.Quit()
	}
}

func (wnd *window) handleMouseEvent(ev *acell.MouseEvent) {
	if wnd.wgt == nil {
		return
	}

	var under []eventHandlerWithPos
	for _, cl := range wnd.wgt {
		if ev.Pos.Y >= cl.p.Line && ev.Pos.Y < cl.p.Line+cl.Height() &&
			ev.Pos.X >= cl.p.Col && ev.Pos.X < cl.p.Col+cl.Width() {
			under = append(under, cl)
		}
	}

	if ev.Action == acell.MouseMove {
		wnd.updateHover(under, ev.Pos)
	}

	for _, cl := range under {
		cl := cl
		wnd.doWithMessage(func() {
			cl.Send(&acell.MouseEvent{
				Action: ev.Action,
				Button: ev.Button,
				Pos: acell.Point{
					X: ev.Pos.X - cl.p.Col,
					Y: ev.Pos.Y - cl.p.Line,
				},
				Shift: ev.Shift,
				Alt:   ev.Alt,
				Ctrl:  ev.Ctrl,
			})
		}, "mouse event")
	}

	for _, h := range wnd.mouseHandlers {
		h := h
		wnd.doWithMessage(func() { h(ev) }, "mouse handler")
	}
}

func (wnd *window) updateHover(under []eventHandlerWithPos, pos acell.Point) {
	now := make(map[EventHandler]struct{}, len(under))
	for _, cl := range under {
		now[cl.EventHandler] = struct{}{}
	}

	for h := range wnd.hovered {
		if _, still := now[h]; still {
			continue
		}
		p := wnd.posOf(h)
		ev := &MouseHoverEvent{
			Entered: false,
			Pos:     acell.Point{X: pos.X - p.Col, Y: pos.Y - p.Line},
		}
		h := h
		wnd.doWithMessage(func() { h.Send(ev) }, "mouse leave")
	}

	for _, cl := range under {
		h := cl.EventHandler
		if _, was := wnd.hovered[h]; was {
			continue
		}
		ev := &MouseHoverEvent{
			Entered: true,
			Pos:     acell.Point{X: pos.X - cl.p.Col, Y: pos.Y - cl.p.Line},
		}
		wnd.doWithMessage(func() { h.Send(ev) }, "mouse enter")
	}

	wnd.hovered = now
}

func (wnd *window) posOf(h EventHandler) Pos {
	for _, cl := range wnd.wgt {
		if cl.EventHandler == h {
			return cl.p
		}
	}
	return Pos{}
}

func (wnd *window) RegisterClickHandler(h MouseEventHandler) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("RegisterClickHandler called outside worker goroutine: data race")
	}
	wnd.mouseHandlers = append(wnd.mouseHandlers, h)
}

func (wnd *window) RegisterResizeHandler(h ResizeHandler) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("RegisterResizeHandler called outside worker goroutine: data race")
	}
	wnd.resizeHandlers = append(wnd.resizeHandlers, h)
}

func (wnd *window) startInputCatcher() {
	wnd.Do(func() {
		wnd.RegisterKeyHandler(func(ke *acell.KeyboardEvent) {
			if wnd.focusIndex != -1 {
				wnd.focusableWidgets[wnd.focusIndex].Send(ke)
			}
			if !wnd.focusChange {
				return
			}
			switch ke.Key {
			case acell.KeyTab:
				wnd.NextFocus()
			case acell.KeyShiftTab:
				wnd.BeforeFocus()
			case acell.KeyCtrlC:
				wnd.Quit()
			}
		})
	})

	for {
		select {
		case <-wnd.stopCh:
			return
		case ev, ok := <-wnd.t.Events():
			if !ok {
				return
			}
			switch e := ev.(type) {
			case *acell.KeyboardEvent:
				key := e
				wnd.doWithMessage(func() {
					for _, h := range wnd.keyHandlers {
						h(key)
					}
				}, "key handler")
			case *acell.MouseEvent:
				mouse := e
				wnd.doWithMessage(func() { wnd.handleMouseEvent(mouse) }, "mouse dispatch")
			case *acell.ResizeEvent:
				wnd.Do(func() {
					wnd.t.Buf = acell.NewBuf(e.Width, e.Height)
					wnd.t.Invalidate()
					for _, h := range wnd.resizeHandlers {
						h(e.Width, e.Height)
					}
					wnd.Redraw()
				})
			}
		}
	}
}

func (wnd *window) SetContent(w Widget) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetContent called outside worker goroutine: data race")
	}
	wnd.content = w
	wnd.Index()
	wnd.ClearFocus()
}

func (wnd *window) SetTitle(title string) {
	if wnd.capture {
		return
	}
	title = ansi.Strip(title)
	wnd.Do(func() {
		if wnd.t != nil {
			wnd.t.SetTitle(title)
		}
	})
}

func (wnd *window) Focus() FocusManager { return wnd }

func (wnd *window) Commit(f func()) {
	wnd.Do(func() {
		f()
		wnd.Redraw()
	})
}

// SetInitCell устанавливает ячейку по умолчанию для всех пустых позиций окна.
func (wnd *window) SetInitCell(c acell.Cell) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetInitCell called outside worker goroutine: data race")
	}
	wnd.initCell = c
	if wnd.t != nil {
		wnd.t.Invalidate()
	}
}

// SetBackground устанавливает фон пустых позиций окна.
func (wnd *window) SetBackground(s Style) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetBackground called outside worker goroutine: data race")
	}
	wnd.SetInitCell(acell.Cell{Char: ' ', Style: ConvertToCellStyle(s)})
}

// SetStyleFunc устанавливает функцию для стилизации виджетов.
func (wnd *window) SetStyleFunc(fn func(Widget)) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetStyleFunc called outside worker goroutine: data race")
	}
	wnd.styleFunc = fn
}

type fakeRawTerminal struct {
	width, height int
	events        chan any
	out           bytes.Buffer
	mu            sync.Mutex
}

func (f *fakeRawTerminal) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.out.Write(p)
}
func (f *fakeRawTerminal) Read(p []byte) (int, error) { return 0, io.EOF }
func (f *fakeRawTerminal) MakeRaw() error             { return nil }
func (f *fakeRawTerminal) Restore() error             { return nil }
func (f *fakeRawTerminal) EnableANSI() error          { return nil }
func (f *fakeRawTerminal) Events() <-chan any         { return f.events }
func (f *fakeRawTerminal) Close() error               { return nil }
func (f *fakeRawTerminal) Size() (int, int)           { return f.width, f.height }
func (f *fakeRawTerminal) StartInput()                {}
func (f *fakeRawTerminal) Info() terminfo.Info        { return terminfo.Info{} }

// Reset очищает буфер между итерациями теста.
func (f *fakeRawTerminal) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.out.Reset()
}
