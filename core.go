package tui

import (
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

	"github.com/romanSPB15/tui-compose/v4/builder"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
	termL "github.com/romanSPB15/tui-compose/v4/term"
	"golang.org/x/term"
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

type eventHandlerWithPos struct {
	EventHandler
	p Pos
}

var currentWindow *window

type window struct {
	focusableWidgets []eventHandlerWithPos
	wgt              []eventHandlerWithPos
	f                io.Writer
	focusIndex       int
	stopCh           chan struct{}
	keyHandlers      []KeyboardEventHandler
	log              *os.File
	runned           bool
	work             chan *task
	focusChange      bool

	mouseHandlers []MouseEventHandler
	content       Widget
	initCell      cell.Cell
	bufferPool    *sync.Pool
	cellBuf       []cell.Cell
	worker        atomic.Int32
	builderPool   sync.Pool
	styleFunc     func(Widget)
	widgetBuf     [][]cell.Cell
	maxWidgetSize Pos
	subBuf        [][]cell.Cell

	stdout    *os.File
	stderr    *os.File
	oldMode   *term.State
	last      cell.Style
	cursorPos Pos

	buf [][]cell.Cell
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

	if evh, ok := wgt.(EventHandler); ok {
		wnd.wgt = append(wnd.wgt, eventHandlerWithPos{
			EventHandler: evh,
			p:            offset,
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

	if evh, ok := wgt.(EventHandler); ok {
		if isFocusable(evh) {
			wnd.focusableWidgets = append(wnd.focusableWidgets, eventHandlerWithPos{
				EventHandler: evh,
				p:            offset,
			})
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

	wnd.indexClickable(wnd.content, Pos{0, 0})
	wnd.indexFocusable(wnd.content, Pos{0, 0})

	wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line = wnd.calcMaxWidgetSize(wnd.content, 0, 0)

	wnd.applyStyles(wnd.content)
}

func (wnd *window) draw(wgt Widget, rect [2]Pos, buf [][]cell.Cell) {
	if wgt == nil {
		return
	}
	if c, ok := wgt.(Container); ok {
		for i, ch := range c.Child() {
			childRect := [2]Pos{
				{Line: rect[0].Line + c.Pos(i).Line, Col: rect[0].Col + c.Pos(i).Col},
				{Line: rect[1].Line + c.Pos(i).Line, Col: rect[1].Col + c.Pos(i).Col},
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
		newBuf := make([][]cell.Cell, newH)
		for i := range newBuf {
			newBuf[i] = make([]cell.Cell, newW)
		}
		wnd.subBuf = newBuf
	} else if currentH > h*2 || currentW > w*2 {
		newBuf := make([][]cell.Cell, h)
		for i := range newBuf {
			newBuf[i] = make([]cell.Cell, w)
		}
		wnd.subBuf = newBuf
	}

	subBuf := wnd.subBuf[:h]
	for y := range subBuf {
		if len(subBuf[y]) < w {
			subBuf[y] = append(subBuf[y], make([]cell.Cell, w-len(subBuf[y]))...)
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

func (wnd *window) render() [][]cell.Cell {
	h := wnd.Height()
	w := wnd.Width()

	buf := wnd.newBuffer(h, w)

	if wnd.content == nil {
		return buf
	}

	ww, hw := wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line
	if ww == 0 && hw == 0 {
		wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line = wnd.calcMaxWidgetSize(wnd.content, 0, 0)
		ww, hw = wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line
	}

	wnd.draw(wnd.content, [2]Pos{{Line: 0, Col: 0}, {Line: h, Col: w}}, buf)

	return buf
}

func (wnd *window) newBuffer(h, w int) [][]cell.Cell {
	if wnd.bufferPool != nil {
		buf := wnd.bufferPool.Get().([][]cell.Cell)
		if len(buf) != h || len(buf[0]) != w {
			buf := make([][]cell.Cell, h)
			for i := range buf {
				buf[i] = make([]cell.Cell, w)
				for x := 0; x < w; x++ {
					buf[i][x] = wnd.initCell
				}
			}
			return buf
		}

		for y := 0; y < h; y++ {
			row := buf[y]
			for x := 0; x < w; x++ {
				row[x] = wnd.initCell
			}
		}
		return buf
	}

	buf := make([][]cell.Cell, h)
	for i := range buf {
		buf[i] = make([]cell.Cell, w)
		for x := 0; x < w; x++ {
			buf[i][x] = wnd.initCell
		}
	}
	return buf
}

func (wnd *window) newEmptyBuffer(h, w int) [][]cell.Cell {
	buf := make([][]cell.Cell, h)
	for i := range buf {
		buf[i] = make([]cell.Cell, w)
		for x := 0; x < w; x++ {
			buf[i][x] = cell.Cell{Char: ' '}
		}
	}
	return buf
}

func (wnd *window) releaseBuffer(buf [][]cell.Cell) {
	if wnd.bufferPool != nil && buf != nil {
		for y := range buf {
			for x := range buf[y] {
				buf[y][x] = wnd.initCell
			}
		}
		wnd.bufferPool.Put(buf)
	}
}

var capture bool

func reflow(old [][]cell.Cell, newW, newH int, initCell cell.Cell) [][]cell.Cell {
	if len(old) == 0 {
		return emptyBuffer(newW, newH, initCell)
	}

	oldW := len(old[0])
	if oldW == 0 {
		return emptyBuffer(newW, newH, initCell)
	}

	var lines [][]cell.Cell
	if oldW == newW {
		lines = old
	} else {
		var stream []cell.Cell
		stream = make([]cell.Cell, 0, oldW*len(old))
		for y := range old {
			stream = append(stream, old[y]...)
		}

		newLineCount := (len(stream) + newW - 1) / newW
		lines = make([][]cell.Cell, newLineCount)
		for y := range newLineCount {
			lines[y] = make([]cell.Cell, newW)
			for x := range newW {
				idx := y*newW + x
				if idx < len(stream) {
					lines[y][x] = stream[idx]
				} else {
					lines[y][x] = cell.Cell{Char: ' ', Style: initCell.Style}
				}
			}
		}
	}

	result := make([][]cell.Cell, newH)
	if len(lines) >= newH {
		start := len(lines) - newH
		for y := range newH {
			result[y] = lines[start+y]
		}
	} else {
		copy(result, lines)

		for y := len(lines); y < newH; y++ {
			result[y] = emptyRow(newW, initCell)
		}
	}

	return result
}

func emptyBuffer(w, h int, initCell cell.Cell) [][]cell.Cell {
	buf := make([][]cell.Cell, h)
	for y := range buf {
		buf[y] = emptyRow(w, initCell)
	}
	return buf
}

func emptyRow(w int, initCell cell.Cell) []cell.Cell {
	row := make([]cell.Cell, w)
	for x := range row {
		row[x] = initCell
	}
	return row
}

func (wnd *window) Redraw() {
	renderStart := time.Now()
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("Redraw called outside worker goroutine: data race")
	}
	if !wnd.runned {
		return
	}

	h := wnd.Height()
	w := wnd.Width()

	if h <= 0 || w <= 0 {
		return
	}

	newBuf := wnd.render()
	if newBuf == nil {
		return
	}

	if wnd.buf == nil || len(wnd.buf) != h || (len(wnd.buf) > 0 && len(wnd.buf[0]) != w) {
		if wnd.buf != nil {
			wnd.buf = reflow(wnd.buf, w, h, wnd.initCell)
			wnd.cursorPos = Pos{-1, -1}
		} else {
			wnd.buf = wnd.newEmptyBuffer(h, w)
		}
	}

	oldBuf := wnd.buf

	b := wnd.builderPool.Get().(*builder.Builder)
	b.Reset()

	defer func() {
		wnd.releaseBuffer(newBuf)
		wnd.builderPool.Put(b)
	}()

	if capture {
		json.NewEncoder(b).Encode(newBuf)
		b.Copy(wnd.f)
		return
	}

	renderDur := time.Since(renderStart)
	makeStringStart := time.Now()

	for y := range h {
		if len(newBuf[y]) < w {
			continue
		}
		for x := range w {
			if newBuf[y][x] != oldBuf[y][x] {
				if wnd.cursorPos.Line != y || wnd.cursorPos.Col != x {
					b.WriteString("\033[")
					b.WriteInt(y + 1)
					b.WriteByte(';')
					b.WriteInt(x + 1)
					b.WriteByte('H')
				}

				newBuf[y][x].Style.WriteANSI(wnd.last, b)
				wnd.last = newBuf[y][x].Style

				if newBuf[y][x].Char == 0 {
					wnd.LogInfo("null rune detected at [%d, %d]", x, y)
					b.WriteRune(' ')
				} else {
					b.WriteRune(newBuf[y][x].Char)
				}

				oldBuf[y][x] = newBuf[y][x]

				x2, y2 := x+1, y

				if x2 == w {
					x2 = 0
					y2++
				}

				wnd.cursorPos = Pos{
					Line: y2,
					Col:  x2,
				}
			}
		}
	}

	makeStringDur := time.Since(makeStringStart)

	writeStart := time.Now()

	b.Copy(wnd.f)

	writeDur := time.Since(writeStart)

	var fps, fpsIO int
	t := renderDur + makeStringDur
	if t == 0 {
		fps = -1
	} else {
		fps = int(time.Second / t)
	}

	t = renderDur + makeStringDur + writeDur
	if t == 0 {
		fpsIO = -1
	} else {
		fpsIO = int(time.Second / t)
	}

	wnd.LogInfo("Redraw timings: %s %s %s, FPS: %d:%d", renderDur, makeStringDur, writeDur, fps, fpsIO)
}

func (wnd *window) Run() {
	defer func() {
		if err := recover(); err != nil {
			wnd.LogFatal("tui: Произошла паника: %v", err)
		}
		if DEBUG {
			wnd.log.Close()
		}
	}()
	if !capture {
		if !term.IsTerminal(int(os.Stdout.Fd())) {
			wnd.LogFatal("tui: stdout is not terminal")
		}
		if err := termL.MakeRaw(); err != nil {
			wnd.LogInfo("tui: Cannot make raw: %s", err)
		}
	}

	wnd.stdout = os.Stdout
	wnd.stderr = os.Stderr
	os.Stdout, os.Stderr = wnd.log, wnd.log

	if !capture {
		fmt.Fprint(wnd.f, "\033[37m")
		wnd.last = cell.Style{Fg: "37"}

		fmt.Fprint(wnd.f, "\033[2J\033[0m\033[?25l\033[?1006h\033[?1000h")
	}

	go wnd.startStopSignalCatcher()
	go wnd.startScreenResizeChecker()
	go wnd.startInputCatcher()

	wnd.runned = true

	wnd.Redraw()

	wnd.runWorker()

	wnd.runned = false

	wnd.restoreOut()

	if !capture {
		termL.Restore()

		if wnd.last != (cell.Style{}) {
			fmt.Fprint(wnd.f, "\033[0m")
		}
		fmt.Fprint(wnd.f, "\033[2J\033[H\033[?25h")
		fmt.Fprint(wnd.f, "\033[?1006l\033[?1000l")
	}
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
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("IsRunned called outside worker goroutine: data race")
	}
	return wnd.runned
}

const taskBufSize = 32

func NewWindow() Window {
	wnd := &window{f: os.Stdout, stopCh: make(chan struct{}), keyHandlers: []KeyboardEventHandler{},
		work: make(chan *task, taskBufSize), focusIndex: -1, focusChange: true, cellBuf: make([]cell.Cell, 0, 256),
		initCell: cell.Cell{Char: ' '}, builderPool: sync.Pool{
			New: func() any {
				return &builder.Builder{}
			},
		},
		cursorPos: Pos{-1, -1},
	}
	if DEBUG {
		f, err := os.Create(fmt.Sprintf("debug_log_%d", time.Now().UnixMilli()))
		if err != nil {
			log.Fatal(err)
		}
		wnd.log = f
	}
	termL.EnableANSIWindows()
	if DEBUG {
		wnd.worker.Store(getGorID())
	}
	currentWindow = wnd
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
	defer recover()

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
	defer recover()

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

func getGorID() int32 {
	var buf [128]byte
	n := runtime.Stack(buf[:], false)
	line := string(buf[:n])

	parts := strings.SplitN(line, " ", 3)
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
	if capture {
		if w := os.Getenv("TUI_WIDTH"); w != "" {
			if val, err := strconv.Atoi(w); err == nil && val > 0 {
				return val
			}
		}
		return 80
	}
	w, _ := termL.SizeFd(wnd.stdout.Fd())
	return w
}

func (wnd *window) Height() int {
	if capture {
		if h := os.Getenv("TUI_HEIGHT"); h != "" {
			if val, err := strconv.Atoi(h); err == nil && val > 0 {
				return val
			}
		}
		return 24
	}
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
	if wnd.wgt != nil {
		for _, cl := range wnd.wgt {
			if ev.Pos.Y >= cl.p.Line && ev.Pos.Y < cl.p.Line+cl.Height() &&
				ev.Pos.X >= cl.p.Col && ev.Pos.X < cl.p.Col+cl.Width() {
				wnd.doWithMessage(func() {
					ev2 := &input.MouseEvent{
						Button: ev.Button,
						Pos: input.Point{
							X: ev.Pos.X - cl.p.Col,
							Y: ev.Pos.Y - cl.p.Line,
						},
					}
					cl.Send(ev2)
				}, "mouse event")
				break
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
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("RegisterClickHandler called outside worker goroutine: data race")
	}
	wnd.mouseHandlers = append(wnd.mouseHandlers, h)
}

func (wnd *window) CopyToClipboard(text string) {
	termL.CopyToClipboard(text)
}

func (wnd *window) startInputCatcher() {
	wnd.Do(func() {
		wnd.RegisterKeyHandler(func(ke *input.KeyboardEvent) {
			if wnd.focusIndex != -1 {
				wnd.focusableWidgets[wnd.focusIndex].Send(ke)
			}
			if !wnd.focusChange {
				return
			}
			switch ke.Key {
			case input.KeyTab:
				wnd.NextFocus()
			case input.KeyShiftTab:
				wnd.BeforeFocus()
			case input.KeyCtrlC:
				wnd.Quit()
			}
		})
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
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetContent called outside worker goroutine: data race")
	}
	wnd.content = w
	wnd.Index()
	wnd.ClearFocus()
}

func (wnd *window) SetTitle(title string) {
	if !capture {
		fmt.Fprintf(wnd.f, "\033]0;%s\033\\", title)
	}
}

func (wnd *window) Focus() FocusManager {
	return wnd
}

func (wnd *window) Commit(f func()) {
	wnd.Do(func() {
		f()
		wnd.Redraw()
	})
}

func CurrentWindow() Window {
	return currentWindow
}

// SetInitCell устанавливает ячейку по умолчанию для всех пустых позиций окна.
// Обычно используется для установки фона.
func (wnd *window) SetInitCell(c cell.Cell) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetInitCell called outside worker goroutine: data race")
	}
	wnd.initCell = c
	wnd.buf = nil
}

// SetInitCell устанавливает фон пустых позиций окна.
func (wnd *window) SetBackground(s Style) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetBackground called outside worker goroutine: data race")
	}
	wnd.SetInitCell(cell.Cell{Char: ' ', Style: ConvertToCellStyle(s)})
}

// SetStyleFunc устанавливает функцию для стилизации виджетов.
func (wnd *window) SetStyleFunc(fn func(Widget)) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetStyleFunc called outside worker goroutine: data race")
	}
	wnd.styleFunc = fn
}
