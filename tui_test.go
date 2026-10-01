package tui

import (
	"encoding/json"
	"io"
	"strconv"
	"testing"

	"github.com/romanSPB15/acell"
)

func assertBuffer(t *testing.T, i int, expected, actual [][]acell.Cell) {
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

func cells(chars string, styles ...acell.Style) []acell.Cell {
	runes := []rune(chars)

	res := make([]acell.Cell, len(runes))

	var currentStyle acell.Style
	if len(styles) > 0 {
		currentStyle = styles[0]
	}

	sIdx := 1

	for i, ch := range runes {
		res[i] = acell.Cell{Char: ch, Style: currentStyle}
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

func addToWindowSize(c [][]acell.Cell, w, h int, initCell acell.Cell) [][]acell.Cell {
	for i := range c {
		for len(c[i]) < w {
			c[i] = append(c[i], initCell)
		}
	}
	if len(c) < h {
		emptyRow := make([]acell.Cell, w)
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

func (w *widget) Render(buf [][]acell.Cell) {
	c := acell.ParseMultiline(w.text)
	for y := range len(c) {
		for x := range len(c[0]) {
			if y < len(buf) && x < len(buf[0]) {
				buf[y][x] = c[y][x]
			}
		}
	}
}

func (w *widget) Width() int {
	return w.width
}

func (w *widget) Height() int {
	return w.height
}

func TestRender(t *testing.T) {
	width, height := 80, 24

	fakeT := acell.NewWithTerm(&fakeRawTerminal{width: width, height: height})
	fakeT.Buf = acell.NewBuf(width, height)

	wnd := NewWindow(WithTerminal(fakeT)).(*window)
	wnd.capture = true
	wnd.width = width
	wnd.height = height

	tt := []struct {
		Content  Widget
		Expected [][]acell.Cell
	}{
		{
			Content: NewVBox(NewStaticLabel("Hello"), NewButton("", nil)),
			Expected: [][]acell.Cell{
				cells("Hello"),
				cells("    ", acell.Style{Fg: "30", Bg: "47"}),
			},
		},

		{
			Content: NewStaticLabel("12345World").WithStyle(FrRed | BgBlue | Bold),
			Expected: [][]acell.Cell{
				cells("12345World", acell.Style{
					Fg:   "31",
					Bg:   "44",
					Args: acell.Bold,
				}),
			},
		},
		{
			Content:  nil,
			Expected: [][]acell.Cell{},
		},
		{
			Content:  NewVBox(nil, nil, nil),
			Expected: [][]acell.Cell{},
		},
		{
			Content: &widget{
				text:   "123\nhello",
				width:  10,
				height: 1,
			},
			Expected: [][]acell.Cell{
				cells("123"),
			},
		},
	}

	if wnd.Focus() == nil {
		t.Fatal("invalid Window.Focus()")
	}

	for i, tv := range tt {
		wnd.SetContent(tv.Content)
		wnd.render()
		buf := wnd.t.Buf
		assertBuffer(t, i, addToWindowSize(tv.Expected, width, height, wnd.initCell), buf)
	}
}

func TestRedraw(t *testing.T) {
	const (
		width  = 80
		height = 24
	)

	t.Setenv("TUI_WIDTH", strconv.Itoa(width))
	t.Setenv("TUI_HEIGHT", strconv.Itoa(height))

	tt := []struct {
		Content  Widget
		Expected [][]acell.Cell
	}{
		{
			Content: NewVBox(NewStaticLabel("Hello"), NewButton("", nil)),
			Expected: [][]acell.Cell{
				cells("Hello"),
				cells("    ", acell.Style{Fg: "30", Bg: "47"}),
			},
		},
		{
			Content: NewStaticLabel("12345World").WithStyle(FrRed | BgBlue | Bold),
			Expected: [][]acell.Cell{
				cells("12345World", acell.Style{
					Fg:   "31",
					Bg:   "44",
					Args: acell.Bold,
				}),
			},
		},
		{
			Content:  nil,
			Expected: [][]acell.Cell{},
		},
		{
			Content:  NewVBox(nil, nil, nil),
			Expected: [][]acell.Cell{},
		},
		{
			Content: &widget{
				text:   "123\nhello",
				width:  10,
				height: 1,
			},
			Expected: [][]acell.Cell{
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
			Expected: [][]acell.Cell{
				cells("123 hello"),
			},
		},
	}

	for i, tv := range tt {
		fakeTerm := &fakeRawTerminal{width: width, height: height}
		wnd := NewWindow(WithTerminal(acell.NewWithTerm(fakeTerm))).(*window)
		wnd.capture = true
		wnd.width = width
		wnd.height = height
		wnd.runned = true
		wnd.worker.Store(getGorID()) // чтобы DEBUG-проверки не срабатывали

		wnd.SetContent(tv.Content)
		wnd.Redraw()

		var buf2 [][]acell.Cell
		err := json.NewDecoder(&fakeTerm.out).Decode(&buf2)
		if err != nil && err != io.EOF {
			t.Fatalf("#%d: decode error: %v", i, err)
		}

		assertBuffer(t, i, addToWindowSize(tv.Expected, width, height, wnd.initCell), buf2)
	}
}

func TestSize(t *testing.T) {
	const (
		width  = 40
		height = 10
	)

	fakeT := acell.NewWithTerm(&fakeRawTerminal{width: width, height: height})
	fakeT.Buf = acell.NewBuf(width, height)

	wnd := NewWindow(WithTerminal(fakeT)).(*window)
	wnd.capture = true
	wnd.width = width
	wnd.height = height

	if w := wnd.Width(); w != width {
		t.Fatalf("invalid width: expected %d, got %d", width, w)
	}
	if h := wnd.Height(); h != height {
		t.Fatalf("invalid height: expected %d, got %d", height, h)
	}
}
