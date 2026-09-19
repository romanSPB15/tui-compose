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

	"github.com/romanSPB15/tui-compose/v4/ansi"
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
	frames       atomic.Int64 // количество успешных Redraw
	renderNanos  atomic.Int64 // накопленное время render()
	makeStrNanos atomic.Int64 // накопленное время diff + ANSI
	writeNanos   atomic.Int64 // накопленное время write в терминал
	tasksDone    atomic.Int64 // количество выполненных задач
	redrawCount  atomic.Int64 // счётчик для периодической проверки isWorker
}

type window struct {
	focusableWidgets []eventHandlerWithPos
	wgt              []eventHandlerWithPos
	f                io.Writer
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

	cellBuf       []cell.Cell
	maxWidgetSize Pos

	content     Widget
	initCell    cell.Cell
	bufferPool  *sync.Pool
	builderPool sync.Pool
	subBuf      [][]cell.Cell
	buf         [][]cell.Cell
	quitOnce    sync.Once

	stdout *os.File
	stderr *os.File
	stdin  *os.File

	oldMode   *term.State
	last      cell.Style
	cursorPos Pos
	altScreen bool

	stats      windowStats
	lastReport time.Time

	width   int
	height  int
	in      io.Reader
	skipRaw bool
}

func (wnd *window) SetAltScreenEnable(v bool) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetAltScreenEnable called outside worker goroutine: data race")
	}
	wnd.altScreen = v
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
		wnd.wgt = append(wnd.wgt, eventHandlerWithPos{
			EventHandler: evh,
			p:            offset,
		})
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

	wnd.indexWgt(wnd.content, Pos{0, 0})

	wnd.maxWidgetSize.Col, wnd.maxWidgetSize.Line = wnd.calcMaxWidgetSize(wnd.content, 0, 0)

	wnd.applyStyles(wnd.content)
}

func (wnd *window) draw(wgt Widget, rect [2]Pos, buf [][]cell.Cell) {
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
	if h <= 0 || w <= 0 {
		return nil
	}

	buf, _ := wnd.bufferPool.Get().([][]cell.Cell)

	if len(buf) == h && len(buf) > 0 && len(buf[0]) == w {
		for y := 0; y < h; y++ {
			row := buf[y]
			for x := 0; x < w; x++ {
				row[x] = wnd.initCell
			}
		}
		return buf
	}

	if buf != nil {
		wnd.bufferPool.Put(buf)
	}

	buf = make([][]cell.Cell, h)
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

func resizeBuf(old [][]cell.Cell, newW, newH int, initCell cell.Cell, anchorTop bool) [][]cell.Cell {
	if newW <= 0 || newH <= 0 {
		return nil
	}
	if len(old) == 0 || len(old[0]) == 0 {
		return emptyBuffer(newW, newH, initCell)
	}

	oldH := len(old)
	oldW := len(old[0])

	buf := make([][]cell.Cell, newH)
	for y := range buf {
		buf[y] = emptyRow(newW, initCell)
	}

	copyW := oldW
	if newW < copyW {
		copyW = newW
	}

	var srcY, copyH int
	if anchorTop {
		srcY = 0
		if newH < oldH {
			copyH = newH
		} else {
			copyH = oldH
		}
	} else {
		if newH <= oldH {
			srcY = oldH - newH
			copyH = newH
		} else {
			srcY = 0
			copyH = oldH
		}
	}

	for y := 0; y < copyH; y++ {
		copy(buf[y][:copyW], old[srcY+y][:copyW])
	}

	return buf
}

func (wnd *window) Redraw() {
	renderStart := time.Now()

	if DEBUG {
		wnd.stats.redrawCount.Add(1)
		if wnd.stats.redrawCount.Load()%64 == 1 && !wnd.isWorker() {
			wnd.LogFatal("Redraw called outside worker goroutine: data race")
		}
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
			wnd.buf = resizeBuf(wnd.buf, w, h, wnd.initCell, wnd.altScreen)
			wnd.cursorPos = Pos{-1, -1}
		} else {
			wnd.buf = wnd.newEmptyBuffer(h, w)
		}
	}

	oldBuf := wnd.buf

	b := wnd.builderPool.Get().(*builder.Builder)
	b.Reset()

	if !capture {
		b.WriteString("\033[?2026h")
	}

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

	b.WriteString("\033[?2026l")

	makeStringDur := time.Since(makeStringStart)

	writeStart := time.Now()

	b.Copy(wnd.f)

	writeDur := time.Since(writeStart)

	wnd.stats.frames.Add(1)
	wnd.stats.renderNanos.Add(int64(renderDur))
	wnd.stats.makeStrNanos.Add(int64(makeStringDur))
	wnd.stats.writeNanos.Add(int64(writeDur))

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
	r := wnd.stats.renderNanos.Swap(0)
	m := wnd.stats.makeStrNanos.Swap(0)
	wrt := wnd.stats.writeNanos.Swap(0)
	tasks := wnd.stats.tasksDone.Swap(0)

	avgRender := time.Duration(r / frames)
	avgMake := time.Duration(m / frames)
	avgWrite := time.Duration(wrt / frames)

	fpsCPU := float64(frames) / elapsed.Seconds()

	var fpsIO float64
	totalNanos := r + m + wrt
	if totalNanos > 0 {
		fpsIO = float64(frames) / (float64(totalNanos) / float64(time.Second))
	}

	wnd.LogInfo(
		"FPS: %.0f (cpu) / %.0f (io), frames=%d tasks=%d, avg render=%s makeString=%s write=%s",
		fpsCPU, fpsIO, frames, tasks, avgRender, avgMake, avgWrite,
	)
}

func (wnd *window) Run() {
	defer func() {
		if err := recover(); err != nil {
			wnd.LogFatal("tui: Произошла паника: %v", err)
		}
		if DEBUG && wnd.log != nil {
			wnd.log.Close()
		}
	}()

	if !capture && !wnd.skipRaw {
		if wnd.stdin != nil && term.IsTerminal(int(wnd.stdin.Fd())) {
			if err := termL.MakeRawFile(wnd.stdin); err != nil {
				wnd.LogInfo("tui: Cannot make raw: %s", err)
			}
		}
	}

	if !capture {
		fmt.Fprint(wnd.f, "\033[37m")
		wnd.last = cell.Style{Fg: "37"}

		b := wnd.builderPool.Get().(*builder.Builder)
		b.Reset()
		b.WriteString("\033[0m\033[?25l\033[?1003h\033[?1006h")
		if wnd.altScreen {
			b.WriteString("\033[?1049h")
		}
		b.WriteString("\033[2J")
		b.Copy(wnd.f)
		wnd.builderPool.Put(b)
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
		if wnd.skipRaw == false && wnd.stdout != nil {
			termL.Restore()
		}

		if wnd.last != (cell.Style{}) {
			fmt.Fprint(wnd.f, "\033[0m")
		}

		b := wnd.builderPool.Get().(*builder.Builder)
		b.Reset()
		b.WriteString("\033[2J\033[?25h\033[?1003l\033[?1000l")
		if wnd.altScreen {
			b.WriteString("\033[?1049l")
		}
		b.Copy(wnd.f)
		wnd.builderPool.Put(b)
	}
}

func (wnd *window) restoreOut() {
	if wnd.stdout == nil {
		return
	}
	os.Stdout = wnd.stdout
	os.Stderr = wnd.stderr
}

func (wnd *window) Quit() {
	wnd.quitOnce.Do(func() { close(wnd.stopCh) })
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

// Options — настройки окна.
type Options struct {
	In          io.Reader
	Out         io.Writer
	Err         io.Writer
	Width       int
	Height      int
	SkipRawMode bool
	AltScreen   bool
}

type Option func(*Options)

func WithIO(in io.Reader, out, err io.Writer) Option {
	return func(o *Options) { o.In, o.Out, o.Err = in, out, err }
}

func WithSize(w, h int) Option {
	return func(o *Options) { o.Width, o.Height = w, h }
}

func WithSkipRawMode(v bool) Option {
	return func(o *Options) { o.SkipRawMode = v }
}

func WithAltScreen(v bool) Option {
	return func(o *Options) { o.AltScreen = v }
}

const taskBufSize = 32

func NewWindow(opts ...Option) Window {
	o := Options{
		In:        os.Stdin,
		Out:       os.Stdout,
		Err:       os.Stderr,
		AltScreen: true,
	}
	for _, opt := range opts {
		opt(&o)
	}

	wnd := &window{
		f:           o.Out,
		in:          o.In,
		width:       o.Width,
		height:      o.Height,
		skipRaw:     o.SkipRawMode,
		altScreen:   o.AltScreen,
		stopCh:      make(chan struct{}),
		keyHandlers: []KeyboardEventHandler{},
		work:        make(chan *task, taskBufSize),
		focusIndex:  -1,
		focusChange: true,
		cellBuf:     make([]cell.Cell, 0, 256),
		initCell:    cell.Cell{Char: ' '},
		builderPool: sync.Pool{
			New: func() any { return &builder.Builder{} },
		},
		cursorPos:  Pos{-1, -1},
		lastReport: time.Now(),
		bufferPool: &sync.Pool{
			New: func() any { return [][]cell.Cell(nil) },
		},
		hovered: make(map[EventHandler]struct{}),
	}

	if f, ok := o.In.(*os.File); ok {
		wnd.stdin = f
	}
	if f, ok := o.Out.(*os.File); ok {
		wnd.stdout = f
		if e, ok := o.Err.(*os.File); ok {
			wnd.stderr = e
		} else {
			wnd.stderr = f
		}
	}

	if DEBUG {
		f, err := os.Create(fmt.Sprintf("debug_log_%d", time.Now().UnixMilli()))
		if err != nil {
			log.Fatal(err)
		}
		wnd.log = f
	}
	termL.EnableANSIWindowsFile(wnd.stdout)
	if DEBUG {
		wnd.worker.Store(getGorID())
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
			wnd.LogInfo("Воркер остановлен")
			return
		case tsk := <-wnd.work:
			if tsk.msg != "" {
				wnd.LogInfo("Принята задача: '%s'", tsk.msg)
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
			wnd.stats.tasksDone.Add(1)
			if tsk.done != nil {
				close(tsk.done)
			}
			if tsk.msg != "" {
				wnd.LogInfo("Завершена задача: '%s'", tsk.msg)
			}
		}
	}
}

func (wnd *window) Width() int {
	if wnd.width > 0 {
		return wnd.width
	}
	if capture {
		if w := os.Getenv("TUI_WIDTH"); w != "" {
			if val, err := strconv.Atoi(w); err == nil && val > 0 {
				return val
			}
		}
		return 80
	}
	if wnd.stdout != nil {
		w, _ := termL.SizeFd(wnd.stdout.Fd())
		return w
	}
	return 80
}

func (wnd *window) Height() int {
	if wnd.height > 0 {
		return wnd.height
	}
	if capture {
		if h := os.Getenv("TUI_HEIGHT"); h != "" {
			if val, err := strconv.Atoi(h); err == nil && val > 0 {
				return val
			}
		}
		return 24
	}
	if wnd.stdout != nil {
		_, h := termL.SizeFd(wnd.stdout.Fd())
		return h
	}
	return 24
}

func (wnd *window) SetSize(w, h int) {
	wnd.Do(func() {
		wnd.width, wnd.height = w, h
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

func (wnd *window) handleMouseEvent(ev *input.MouseEvent) {
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

	if ev.Action == input.MouseMove {
		wnd.updateHover(under, ev.Pos)
	}

	for _, cl := range under {
		cl := cl
		wnd.doWithMessage(func() {
			cl.Send(&input.MouseEvent{
				Action: ev.Action,
				Button: ev.Button,
				Pos: input.Point{
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

func (wnd *window) updateHover(under []eventHandlerWithPos, pos input.Point) {
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
			Pos:     input.Point{X: pos.X - p.Col, Y: pos.Y - p.Line},
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
			Pos: input.Point{
				X: pos.X - cl.p.Col,
				Y: pos.Y - cl.p.Line,
			},
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

	mouse, keyboard := input.Start(wnd.in, 1)
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
			ev2 := ev
			wnd.doWithMessage(func() {
				wnd.handleMouseEvent(ev2)
			}, "mouse dispatch")
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
		title = ansi.Strip(title)
		wnd.Do(func() {
			fmt.Fprintf(wnd.f, "\033]0;%s\033\\", title)
		})
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

// SetInitCell устанавливает ячейку по умолчанию для всех пустых позиций окна.
// Обычно используется для установки фона.
func (wnd *window) SetInitCell(c cell.Cell) {
	if DEBUG && !wnd.isWorker() {
		wnd.LogFatal("SetInitCell called outside worker goroutine: data race")
	}
	wnd.initCell = c
	wnd.buf = nil
}

// SetBackground устанавливает фон пустых позиций окна.
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
