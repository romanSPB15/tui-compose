package tui

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/romanSPB15/tui-compose/v4/cell"
)

func newBuf(w, h int) [][]cell.Cell {
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	buf := make([][]cell.Cell, h)
	for i := range buf {
		buf[i] = make([]cell.Cell, w)
	}
	return buf
}

func renderOK(t *testing.T, w Widget) {
	t.Helper()
	width := w.Width()
	height := w.Height()
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	w.Render(newBuf(width+16, height+16))
}

func TestLabel(t *testing.T) {
	l := NewStaticLabel("hello")
	if l.Width() != 5 {
		t.Fatalf("Width: got %d, want 5", l.Width())
	}
	if l.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", l.Height())
	}
	renderOK(t, l)

	l.SetText("hi")
	if l.Text != "hi" {
		t.Fatalf("Text after SetText: got %q", l.Text)
	}
	renderOK(t, l)

	dl := NewDynamicLabel("x", 10)
	if dl.Width() != 10 {
		t.Fatalf("Dynamic Width: got %d, want 10", dl.Width())
	}
	renderOK(t, dl)

	dl.WithStyle(FrRed)
	renderOK(t, dl)

	dl.WithText("hello world")
	if dl.Text != "hello world" {
		t.Fatalf("WithText: got %q", dl.Text)
	}
	renderOK(t, dl)
}

func TestButton(t *testing.T) {
	b := NewButton("OK", func() {})
	if b.Width() != 2+2*2 {
		t.Fatalf("Width: got %d, want 6", b.Width())
	}
	if b.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", b.Height())
	}
	renderOK(t, b)

	b.WithPaddings(0, 0)
	if b.Width() != 2 {
		t.Fatalf("Width no padding: got %d, want 2", b.Width())
	}
	renderOK(t, b)

	b.WithPaddings(1, 1)
	if b.Width() != 4 {
		t.Fatalf("Width padding 1: got %d, want 4", b.Width())
	}
	if b.Height() != 3 {
		t.Fatalf("Height padding 1: got %d, want 3", b.Height())
	}
	renderOK(t, b)

	b.WithStyle(FrRed)
	renderOK(t, b)
	b.WithFocusedStyle(FrGreen)
	renderOK(t, b)
	b.WithHoverStyle(FrBlue)
	renderOK(t, b)
	b.WithDisabledStyle(FrBrightBlack)
	renderOK(t, b)

	b.WithText("quit")
	if b.GetText() != "quit" {
		t.Fatalf("GetText: got %q", b.GetText())
	}
	renderOK(t, b)

	called := false
	b.WithHandler(func() { called = true })
	_ = called
	renderOK(t, b)
}

func TestCheck(t *testing.T) {
	c := NewCheck("Enable")
	if c.Width() != 6+4 {
		t.Fatalf("Width: got %d, want 10", c.Width())
	}
	if c.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", c.Height())
	}
	renderOK(t, c)

	if c.State() {
		t.Fatal("initial state should be false")
	}
	c.SetState(true)
	if !c.State() {
		t.Fatal("state should be true")
	}
	renderOK(t, c)

	c.WithState(false).WithText("Off")
	if c.Width() != 3+4 {
		t.Fatalf("Width Off: got %d, want 7", c.Width())
	}
	renderOK(t, c)

	c.WithStyle(FrRed).WithFocusedStyle(FrGreen).WithHoverStyle(FrBlue).WithCheckedStyle(FrCyan)
	renderOK(t, c)

	c.WithOnChanged(func(b bool) {})
	renderOK(t, c)
}

func TestInputField(t *testing.T) {
	f := NewInputField(10)
	if f.Width() != 10 {
		t.Fatalf("Width: got %d, want 10", f.Width())
	}
	if f.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", f.Height())
	}
	renderOK(t, f)

	f.WithText("hello")
	if f.Text != "hello" {
		t.Fatalf("Text: got %q", f.Text)
	}
	if f.CursorPos != 5 {
		t.Fatalf("CursorPos: got %d, want 5", f.CursorPos)
	}
	renderOK(t, f)

	f.WithText("0123456789abcdef")
	if len([]rune(f.Text)) != 10 {
		t.Fatalf("Text should be truncated to 10, got %d", len([]rune(f.Text)))
	}
	renderOK(t, f)

	f.WithPlaceholder("type here")
	f.WithPlaceholderStyle(FrBrightBlack)
	renderOK(t, f)

	f.WithStyle(BgRed).WithFocusedStyle(BgBlue).WithHoverStyle(BgGreen).WithCursorStyle(FrRed)
	renderOK(t, f)

	f.WithBlinkInterval(time.Second).DisableBlink().EnableBlink()
	f.WithOverwrite(true)
	if !f.IsOverwrite() {
		t.Fatal("overwrite should be true")
	}
	f.WithOnChanged(func(s string) {}).WithOnEnter(func(s string) {})
	renderOK(t, f)
}

func TestHyperlink(t *testing.T) {
	h := NewHyperlink("click", "https://example.com")
	if h.Width() != 5 {
		t.Fatalf("Width: got %d, want 5", h.Width())
	}
	if h.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", h.Height())
	}
	renderOK(t, h)
}

func TestGauge(t *testing.T) {
	g := NewGauge(20)
	if g.Width() != 20 {
		t.Fatalf("Width: got %d, want 20", g.Width())
	}
	if g.Height() != 1 {
		t.Fatalf("Height: got %d, want 1", g.Height())
	}
	if NewGauge(1).Width() != 4 {
		t.Fatalf("min width: got %d, want 4", NewGauge(1).Width())
	}
	renderOK(t, g)

	g.WithValue(0.5)
	renderOK(t, g)
	g.WithValue(-1)
	renderOK(t, g)
	g.WithValue(2)
	renderOK(t, g)

	renderOK(t, g.ASCII())
	renderOK(t, g.Default())
	renderOK(t, g.EmptySquares())
	renderOK(t, g.Blocks())
	renderOK(t, g.BlocksFull())
	renderOK(t, g.BlocksGrid())

	g.WithOnCell(cell.Cell{Char: '#'})
	g.WithOffCell(cell.Cell{Char: '-'})
	renderOK(t, g)

	g.WithOnStyle(FrRed).WithOffStyle(FrBlue)
	renderOK(t, g)
	g.WithOnChar('+').WithOffChar('=')
	renderOK(t, g)
	g.WithLabel("loading").WithLabelStyle(FrGreen)
	renderOK(t, g)
	g.WithLabelFunc(func(f float64) string { return "x" })
	renderOK(t, g)
}

func renderFrame(t *testing.T, f *Frame) {
	t.Helper()
	buf := newBuf(f.Width()+2, f.Height()+2)
	f.border.Render(buf)
}

func TestFrame(t *testing.T) {
	content := NewStaticLabel("content")
	f := NewFrame(content)
	if f.Width() != content.Width()+1*2+2 {
		t.Fatalf("Width: got %d, want %d", f.Width(), content.Width()+1*2+2)
	}
	if f.Height() != content.Height()+2 {
		t.Fatalf("Height: got %d, want %d", f.Height(), content.Height()+2)
	}
	renderFrame(t, f)

	renderFrame(t, f.Heavy())
	renderFrame(t, f.ASCII())
	renderFrame(t, f.Rounded())
	renderFrame(t, f.Double())
	renderFrame(t, f.Default())
	renderFrame(t, f.Dashed())
	renderFrame(t, f.RoundedDashed())
	renderFrame(t, f.Bevel())
	renderFrame(t, f.BevelASCII())
	renderFrame(t, f.Custom('1', '2', '3', '4', '5', '6'))

	f.WithPaddings(1, 2)
	if f.Width() != content.Width()+2*2+2 {
		t.Fatalf("Width with h padding: got %d", f.Width())
	}
	if f.Height() != content.Height()+1*2+2 {
		t.Fatalf("Height with v padding: got %d", f.Height())
	}
	renderFrame(t, f)

	f.WithBorderStyle(FrRed)
	renderFrame(t, f)

	f.WithTitle(Title{Text: "TopLeft", Pos: TitleTopLeft, Style: FrGreen})
	f.WithTitle(Title{Text: "TopRight", Pos: TitleTopRight, Style: FrBlue})
	f.WithTitle(Title{Text: "TopCenter", Pos: TitleTopCenter, Style: FrRed})
	f.WithTitle(Title{Text: "BotLeft", Pos: TitleBottomLeft, Style: FrGreen})
	f.WithTitle(Title{Text: "BotRight", Pos: TitleBottomRight, Style: FrBlue})
	f.WithTitle(Title{Text: "BotCenter", Pos: TitleBottomCenter, Style: FrRed})
	renderFrame(t, f)

	f.WithInitCell(cell.Cell{Char: '.'})
	f.WithBackground(BgRed)
	renderFrame(t, f)

	if len(f.Child()) != 2 {
		t.Fatalf("Child len: got %d, want 2", len(f.Child()))
	}
	_ = f.Pos(0)
	_ = f.Pos(1)

	emptyF := NewFrame(nil)
	if emptyF.Width() != 1*2+2 {
		t.Fatalf("empty Width: got %d, want 4", emptyF.Width())
	}
	if emptyF.Height() != 0*2+2 {
		t.Fatalf("empty Height: got %d, want 2", emptyF.Height())
	}
	emptyBuf := newBuf(emptyF.Width()+2, emptyF.Height()+2)
	emptyF.border.Render(emptyBuf)
}

func TestImage(t *testing.T) {
	img := NewImage()
	if img.Width() != 0 {
		t.Fatalf("empty Width: got %d, want 0", img.Width())
	}
	if img.Height() != 0 {
		t.Fatalf("empty Height: got %d, want 0", img.Height())
	}

	img = make(Image, 3)
	for i := range img {
		img[i] = make([]cell.Cell, 5)
	}
	if img.Width() != 5 {
		t.Fatalf("Width: got %d, want 5", img.Width())
	}
	if img.Height() != 3 {
		t.Fatalf("Height: got %d, want 3", img.Height())
	}
	renderOK(t, img)
}

func makeTestImage(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8((x * 17) % 256),
				G: uint8((y * 29) % 256),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func TestImageLoad(t *testing.T) {
	src := makeTestImage(4, 4)
	modes := []LoadMode{
		PaletteTrueColor | HalfSymbol,
		PaletteTrueColor | OneSymbol,
		PaletteTrueColor | TwoSymbol,
		Palette16Color | HalfSymbol,
		Palette16Color | OneSymbol,
		Palette16Color | TwoSymbol,
	}
	for _, m := range modes {
		var img Image
		img = img.LoadImage(src, m)
		if img.Height() < 1 {
			t.Fatalf("mode %d: Height %d", m, img.Height())
		}
		renderOK(t, img)
	}
}

func TestDownscaleImage(t *testing.T) {
	src := makeTestImage(10, 10)
	dst := DownscaleImage(src, 5, 5)
	if dst == nil {
		t.Fatal("nil result")
	}
	if dst.Bounds().Dx() != 5 || dst.Bounds().Dy() != 5 {
		t.Fatalf("bad size: %v", dst.Bounds())
	}

	if DownscaleImage(src, 0, 5) != nil {
		t.Fatal("expected nil for zero width")
	}
	if DownscaleImage(src, 5, 0) != nil {
		t.Fatal("expected nil for zero height")
	}
	if DownscaleImage(src, -1, 5) != nil {
		t.Fatal("expected nil for negative width")
	}
}

func TestScaleToHeight(t *testing.T) {
	src := makeTestImage(20, 10)
	dst := ScaleToHeight(src, 5)
	if dst == nil {
		t.Fatal("nil result")
	}
	if dst.Bounds().Dy() != 5 {
		t.Fatalf("height: got %d, want 5", dst.Bounds().Dy())
	}
}

func TestScaleToWidth(t *testing.T) {
	src := makeTestImage(20, 10)
	dst := ScaleToWidth(src, 10)
	if dst == nil {
		t.Fatal("nil result")
	}
	if dst.Bounds().Dx() != 10 {
		t.Fatalf("width: got %d, want 10", dst.Bounds().Dx())
	}
}

func TestBrailleChar(t *testing.T) {
	if brailleChar([4][2]bool{}) != 0x2800 {
		t.Fatalf("empty: got %U", brailleChar([4][2]bool{}))
	}

	var all [4][2]bool
	for y := 0; y < 4; y++ {
		for x := 0; x < 2; x++ {
			all[y][x] = true
		}
	}
	if brailleChar(all) != 0x28FF {
		t.Fatalf("all: got %U", brailleChar(all))
	}

	var one [4][2]bool
	one[0][0] = true
	if brailleChar(one) != 0x2801 {
		t.Fatalf("one: got %U", brailleChar(one))
	}
}

func TestLoadBraille(t *testing.T) {
	var img Image
	if img.LoadBraille(nil, FrRed) != nil {
		t.Fatal("nil data should return nil")
	}
	if img.LoadBraille([][]bool{}, FrRed) != nil {
		t.Fatal("empty should return nil")
	}

	small := [][]bool{{true}}
	if img.LoadBraille(small, FrRed) != nil {
		t.Fatal("too small should return nil")
	}

	data := [][]bool{
		{true, false},
		{false, true},
		{true, false},
		{false, true},
	}
	img = img.LoadBraille(data, FrRed)
	if img.Width() != 1 || img.Height() != 1 {
		t.Fatalf("size: got %dx%d, want 1x1", img.Width(), img.Height())
	}
	renderOK(t, img)

	big := make([][]bool, 8)
	for i := range big {
		big[i] = make([]bool, 4)
		for j := range big[i] {
			big[i][j] = (i+j)%2 == 0
		}
	}
	img = img.LoadBraille(big, FrGreen)
	if img.Width() != 2 || img.Height() != 2 {
		t.Fatalf("big size: got %dx%d, want 2x2", img.Width(), img.Height())
	}
	renderOK(t, img)
}

func TestNearestANSI16Fg(t *testing.T) {
	cases := []struct {
		r, g, b uint8
		want    string
	}{
		{0, 0, 0, "30"},
		{255, 255, 255, "97"},
		{255, 0, 0, "91"},
		{0, 255, 0, "92"},
		{0, 0, 255, "94"},
		{192, 192, 192, "37"},
	}
	for _, c := range cases {
		got := nearestANSI16Fg(c.r, c.g, c.b)
		if got != c.want {
			t.Fatalf("fg(%d,%d,%d): got %q, want %q", c.r, c.g, c.b, got, c.want)
		}
	}
}

func TestNearestANSI16Bg(t *testing.T) {
	cases := []struct {
		r, g, b uint8
		want    string
	}{
		{0, 0, 0, "40"},
		{255, 255, 255, "107"},
		{255, 0, 0, "101"},
		{0, 255, 0, "102"},
		{0, 0, 255, "104"},
		{192, 192, 192, "47"},
	}
	for _, c := range cases {
		got := nearestANSI16Bg(c.r, c.g, c.b)
		if got != c.want {
			t.Fatalf("bg(%d,%d,%d): got %q, want %q", c.r, c.g, c.b, got, c.want)
		}
	}
}

func TestDisableState(t *testing.T) {
	var d DisableState
	if d.IsDisabled() {
		t.Fatal("zero value should not be disabled")
	}
	d.SetDisabled(true)
	if !d.IsDisabled() {
		t.Fatal("should be disabled")
	}
	d.SetDisabled(false)
	if d.IsDisabled() {
		t.Fatal("should not be disabled")
	}
}
