package tui

import (
	"encoding/json"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/romanSPB15/tui-compose/v4/cell"
)

func assertBuffer(t *testing.T, i int, expected, actual [][]cell.Cell) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Fatalf("#%d: height mismatch: expected %d, got %d", i, len(expected), len(actual))
	}
	for y := 0; y < len(expected); y++ {
		if len(expected[y]) != len(actual[y]) {
			t.Fatalf("#%d: width mismatch at row %d: expected %d, got %d", i, y, len(expected[y]), len(actual[y]))
		}
		for x := 0; x < len(expected[y]); x++ {
			if expected[y][x] != actual[y][x] {
				t.Errorf("#%d: cell mismatch at (%d,%d): expected %v, got %v", i, x, y, expected[y][x], actual[y][x])
			}
		}
	}
}

func cells(chars string, styles ...cell.Style) []cell.Cell {
	runes := []rune(chars)

	res := make([]cell.Cell, len(runes))

	var currentStyle cell.Style
	if len(styles) > 0 {
		currentStyle = styles[0]
	}

	sIdx := 1

	for i, ch := range runes {
		res[i] = cell.Cell{Char: ch, Style: currentStyle}
		if sIdx < len(styles) {
			currentStyle = styles[sIdx]
			sIdx++
		}
	}
	return res
}

const (
	width  = 40
	height = 10
)

func addToWindowSize(c [][]cell.Cell, w, h int, initCell cell.Cell) [][]cell.Cell {
	for i := range c {
		for len(c[i]) < w {
			c[i] = append(c[i], initCell)
		}
	}
	if len(c) < h {
		emptyRow := make([]cell.Cell, w)
		for i := range emptyRow {
			emptyRow[i] = initCell
		}
		for len(c) < h {
			c = append(c, emptyRow)
		}
	}
	return c
}

type widget struct {
	width, height int
	text          string
}

func (w *widget) Render(buf [][]cell.Cell) {
	c := cell.ParseMultiline(w.text)
	copy(buf, c)
}

func (w *widget) Width() int {
	return w.width
}

func (w *widget) Height() int {
	return w.height
}

func TestRender(t *testing.T) {
	capture = true
	t.Cleanup(func() { capture = false })

	t.Setenv("TUI_WIDTH", strconv.Itoa(width))
	t.Setenv("TUI_HEIGHT", strconv.Itoa(height))

	wnd := NewWindow().(*window)

	tt := []struct {
		Content  Widget
		Expected [][]cell.Cell
	}{
		{
			Content: NewVBox(NewStaticLabel("Hello"), NewButton("", nil)),
			Expected: [][]cell.Cell{
				cells("Hello"),
			},
		},

		{
			Content: NewStaticLabel("12345World").WithStyle(FrRed | BgBlue | Bold),
			Expected: [][]cell.Cell{
				cells("12345World", cell.Style{
					Fg:   "31",
					Bg:   "44",
					Args: cell.Bold,
				}),
			},
		},
		{
			Content:  nil,
			Expected: [][]cell.Cell{},
		},
		{
			Content:  NewVBox(nil, nil, nil),
			Expected: [][]cell.Cell{},
		},
		{
			Content: &widget{
				text:   "123\nhello",
				width:  10,
				height: 1,
			},
			Expected: [][]cell.Cell{
				cells("123"),
			},
		},
	}

	if CurrentWindow() != wnd {
		t.Fatal("invalid CurrentWindow()")
	}

	if wnd.Focus() == nil {
		t.Fatal("invalid Window.Focus()")
	}

	wnd.SetTitle("123")

	for i, tv := range tt {
		wnd.SetContent(tv.Content)
		buf := wnd.render()

		assertBuffer(t, i, addToWindowSize(tv.Expected, width, height, wnd.initCell), buf)
	}
}
func TestRedraw(t *testing.T) {
	capture = true
	t.Cleanup(func() { capture = false })

	t.Setenv("TUI_WIDTH", strconv.Itoa(width))
	t.Setenv("TUI_HEIGHT", strconv.Itoa(height))

	tt := []struct {
		Content  Widget
		Expected [][]cell.Cell
	}{
		{
			Content: NewVBox(NewStaticLabel("Hello"), NewButton("", nil)),
			Expected: [][]cell.Cell{
				cells("Hello"),
			},
		},
		{
			Content: NewStaticLabel("12345World").WithStyle(FrRed | BgBlue | Bold),
			Expected: [][]cell.Cell{
				cells("12345World", cell.Style{
					Fg:   "31",
					Bg:   "44",
					Args: cell.Bold,
				}),
			},
		},
		{
			Content:  nil,
			Expected: [][]cell.Cell{},
		},
		{
			Content:  NewVBox(nil, nil, nil),
			Expected: [][]cell.Cell{},
		},
		{
			Content: &widget{
				text:   "123\nhello",
				width:  10,
				height: 1,
			},
			Expected: [][]cell.Cell{
				cells("123"),
			},
		},
		{
			Content: NewHBox(&widget{
				text:   "123\nhello",
				width:  3,
				height: 1,
			}, &widget{
				text:   "hello\n123",
				width:  10,
				height: 1,
			}),
			Expected: [][]cell.Cell{
				cells("123 hello"),
			},
		},
	}

	for i, tv := range tt {
		wnd := NewWindow().(*window)
		wnd.runned = true
		pr, pw := io.Pipe()
		wnd.f = pw

		done := make(chan struct{})
		go func() {
			wnd.SetContent(tv.Content)
			wnd.Redraw()
			pw.Close()
			close(done)
		}()

		var buf2 [][]cell.Cell
		err := json.NewDecoder(pr).Decode(&buf2)
		if err != nil && err != io.EOF {
			t.Fatalf("#%d: decode error: %v", i, err)
		}

		<-done

		assertBuffer(t, i, addToWindowSize(tv.Expected, width, height, wnd.initCell), buf2)
	}
}

func TestSize(t *testing.T) {
	capture = true
	t.Cleanup(func() { capture = false })

	t.Setenv("TUI_WIDTH", strconv.Itoa(width))
	t.Setenv("TUI_HEIGHT", strconv.Itoa(height))

	wnd := NewWindow()

	w := wnd.Width()
	if w != width {
		t.Fatalf("invalid width: expected %d, but got %d", width, w)
	}

	h := wnd.Height()
	if h != height {
		t.Fatalf("invalid height: expected %d, but got %d", height, h)
	}
}

func TestRun(t *testing.T) {
	capture = true
	t.Cleanup(func() { capture = false })

	t.Setenv("TUI_WIDTH", strconv.Itoa(width))
	t.Setenv("TUI_HEIGHT", strconv.Itoa(height))

	wnd := NewWindow().(*window)
	pr, pw := io.Pipe()
	wnd.f = pw

	content := NewVBox(NewStaticLabel("Hello"), NewButton("", nil))
	wnd.SetContent(content)

	var buf2 [][]cell.Cell
	doneReading := make(chan struct{})
	go func() {
		err := json.NewDecoder(pr).Decode(&buf2)
		if err != nil && err != io.EOF {
			t.Logf("decode error: %v", err)
		}
		close(doneReading)
	}()

	doneRun := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered panic in Run: %v", r)
				t.Fail()
			}
		}()
		wnd.Run()
		close(doneRun)
	}()

	time.Sleep(100 * time.Millisecond)

	wnd.Quit()

	select {
	case <-doneRun:
		// OK
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not quit in time")
	}

	pw.Close()
	<-doneReading

	expected := addToWindowSize([][]cell.Cell{
		cells("Hello"),
	}, width, height, wnd.initCell)

	assertBuffer(t, 0, expected, buf2)
}
