package extra

import (
	"testing"

	"github.com/romanSPB15/tui-compose/v4"
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

func renderOK(t *testing.T, w tui.Widget) {
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

func TestSpinner(t *testing.T) {
	cases := []struct {
		typ  int
		want int
	}{
		{SpinnerUnderscore, 1},
		{SpinnerBraille, 1},
		{SpinnerDots, 3},
		{SpinnerLine, 1},
		{SpinnerBrailleReverse, 1},
	}
	for _, c := range cases {
		s := NewSpinner(c.typ)
		if s.Height() != 1 {
			t.Fatalf("spinner %d Height: got %d, want 1", c.typ, s.Height())
		}
		if s.Width() != c.want {
			t.Fatalf("spinner %d Width: got %d, want %d", c.typ, s.Width(), c.want)
		}
		renderOK(t, s)
	}
}

func TestTable(t *testing.T) {
	data := [][]TableCell{
		{{Text: "abc", Align: AlignLeft}, {Text: "de", Align: AlignRight}},
		{{Text: "f", Align: AlignCenter}, {Text: "ghij", Align: AlignLeft}},
	}
	tbl := NewTable(data)
	if tbl.Width() != 1+(3+3)+(4+3) {
		t.Fatalf("Width: got %d", tbl.Width())
	}
	if tbl.Height() != 3 {
		t.Fatalf("Height: got %d", tbl.Height())
	}
	renderOK(t, tbl)
	renderOK(t, NewTable(data).WithHorSeparator(NoHorSeparator))
	renderOK(t, NewTable(data).WithHorSeparator(BetweenHorSeparator))
	renderOK(t, NewTable(data).WithHorSeparator(EverywhereHorSeparator))
	renderOK(t, tbl.ASCII())
	renderOK(t, tbl.Rounded())
}

func TestTableEmpty(t *testing.T) {
	tbl := NewTable(nil)
	tbl.Render(newBuf(1, 1))
	tbl2 := NewTable([][]TableCell{})
	tbl2.Render(newBuf(1, 1))
}

func TestTabs(t *testing.T) {
	tv1 := NewTextView(2).WithLines([]string{"abc", "de"})
	tv2 := NewTextView(1).WithLines([]string{"x"})
	tabs := NewTabs([]Tab{
		{Title: "A", Content: tv1},
		{Title: "BB", Content: tv2},
	})
	if tabs.Height() != tv1.Height()+1 {
		t.Fatalf("Height: got %d, want %d", tabs.Height(), tv1.Height()+1)
	}
	wantW := max(tv1.Width(), 4)
	if tabs.Width() != wantW {
		t.Fatalf("Width: got %d, want %d", tabs.Width(), wantW)
	}
	tabs.Render(newBuf(tabs.Width()+16, tabs.Height()+16))
	renderOK(t, &tabs.topPanel)
	if len(tabs.Child()) == 0 {
		t.Fatal("Tabs should return children")
	}
}

func TestTextArea(t *testing.T) {
	ta := NewTextArea(10, 3)
	if ta.Width() != 10 {
		t.Fatalf("Width: got %d", ta.Width())
	}
	if ta.Height() != 3 {
		t.Fatalf("Height: got %d", ta.Height())
	}
	renderOK(t, ta)

	ta.WithText("a\nb\nc\nd\ne\nf\ng\nh\ni\nj")
	ta.WithLineNumbers(true)
	if ta.Width() != 2+1+10 {
		t.Fatalf("Width with line numbers: got %d, want %d", ta.Width(), 2+1+10)
	}
	renderOK(t, ta)
}

func TestTextView(t *testing.T) {
	tv := NewTextView(3).WithLines([]string{
		"[fg-red]hi[-]",
		"longer line",
		"[bold]x[-]",
	})
	if tv.Height() != 3 {
		t.Fatalf("Height: got %d", tv.Height())
	}
	if tv.Width() != 11 {
		t.Fatalf("Width: got %d, want 11", tv.Width())
	}
	renderOK(t, tv)

	tv.WithFixedWidth(7)
	if tv.Width() != 7 {
		t.Fatalf("Width fixed: got %d", tv.Width())
	}
	renderOK(t, tv)
}

func TestTree(t *testing.T) {
	tr := NewTree([]TreeNode{
		{Label: "root"},
		{Label: "child", Depth: 1},
		{Label: "deeper", Depth: 2},
	})
	if tr.Height() != 3 {
		t.Fatalf("Height: got %d", tr.Height())
	}
	if tr.Width() != 12 {
		t.Fatalf("Width: got %d, want 12", tr.Width())
	}
	renderOK(t, tr)
	renderOK(t, tr.ASCII())
	renderOK(t, tr.Rounded())
	renderOK(t, tr.Heavy())
}

func TestAccordion(t *testing.T) {
	content := NewTextView(1).WithLines([]string{"hello world"})
	acc := NewAccordion("Title", content)
	if acc.Width() < 1 {
		t.Fatalf("Width: got %d", acc.Width())
	}
	if acc.Height() < 1 {
		t.Fatalf("Height: got %d", acc.Height())
	}
	if len(acc.Child()) == 0 {
		t.Fatal("Accordion should return children")
	}
	acc.Render(newBuf(acc.Width()+16, acc.Height()+16))
}

func TestBarChart(t *testing.T) {
	bc := NewBarChart().WithValues([]int{1, 2, 3}).WithBarWidth(2).WithSpace(1).WithDataHeight(5)
	if bc.Width() != 3*(2+1)-1 {
		t.Fatalf("Width: got %d, want %d", bc.Width(), 3*(2+1)-1)
	}
	if bc.Height() != 5 {
		t.Fatalf("Height: got %d", bc.Height())
	}
	renderOK(t, bc)
	renderOK(t, bc.AutoScale())
	renderOK(t, bc.WithBarStyle(func(i, v int) tui.Style { return tui.FrRed }))
	renderOK(t, bc.WithTextStyle(func(i, v int) tui.Style { return tui.FrGreen }))
}

func TestBlinkLabel(t *testing.T) {
	bl := NewBlinkLabel(5).WithText("hello")
	if bl.Width() < 1 {
		t.Fatalf("Width: got %d", bl.Width())
	}
	if bl.Height() < 1 {
		t.Fatalf("Height: got %d", bl.Height())
	}
	renderOK(t, bl)
}

func TestLineChart(t *testing.T) {
	lc := NewLineChart().WithData([]Series{
		{Values: []int{1, 2, 3, 4, 5}},
	}).WithPointDistance(4).WithDataHeight(10)
	if lc.Height() != 10 {
		t.Fatalf("Height: got %d", lc.Height())
	}
	if lc.Width() != (5-1)*4+1 {
		t.Fatalf("Width: got %d, want %d", lc.Width(), (5-1)*4+1)
	}
	renderOK(t, lc)

	lc.WithYLabels([]int{1, 3, 5})
	if lc.Width() != (5-1)*4+1+4 {
		t.Fatalf("Width with YLabels: got %d", lc.Width())
	}
	if lc.Height() != 11 {
		t.Fatalf("Height with YLabels: got %d", lc.Height())
	}
	renderOK(t, lc)

	lc.WithXLabels([]string{"a", "b", "c", "d", "e"}).WithDisplayPoints(true)
	renderOK(t, lc)

	renderOK(t, lc.WithASCIIAxis())
	renderOK(t, lc.WithRoundedAxis())
	renderOK(t, lc.WithDefaultAxis())
	renderOK(t, lc.AutoScale())
}

func TestPageIndicator(t *testing.T) {
	pi := NewPageIndicator(7).SetCurrent(3)
	if pi.Width() != 7 {
		t.Fatalf("Width: got %d", pi.Width())
	}
	if pi.Height() != 1 {
		t.Fatalf("Height: got %d", pi.Height())
	}
	renderOK(t, pi)
	renderOK(t, pi.WithActive('X').WithInactive('o').WithStyle(tui.FrRed))
}

func TestPieChart(t *testing.T) {
	pc := NewPieChart([]PieData{
		{Label: "A", Value: 30, Color: tui.BgRed},
		{Label: "B", Value: 70, Color: tui.BgBlue},
	}).WithRadius(10)
	if pc.Height() != 10 {
		t.Fatalf("Height: got %d", pc.Height())
	}
	if pc.Width() != 2*10+2 {
		t.Fatalf("Width: got %d", pc.Width())
	}
	renderOK(t, pc)
	renderOK(t, pc.WithShowLegend(true))
	renderOK(t, pc.WithShowPercent(false))
	renderOK(t, pc.WithValueStyle(tui.FrRed).WithLegendStyle(tui.FrGreen))
}

func TestSlider(t *testing.T) {
	s := NewSlider(20)
	if s.Width() != 20 {
		t.Fatalf("Width: got %d", s.Width())
	}
	if s.Height() != 1 {
		t.Fatalf("Height: got %d", s.Height())
	}
	if NewSlider(1).Width() != 3 {
		t.Fatalf("min width: got %d", NewSlider(1).Width())
	}
	renderOK(t, s)
	renderOK(t, s.WithRange(0, 100).WithValue(50).WithStep(5))
	renderOK(t, s.WithLabelFunc(func(f float64) string { return "x" }))
	renderOK(t, s.ASCII())
}

func TestSparkline(t *testing.T) {
	sp := NewSparkline().WithValues([]int{1, 2, 3, 4}).WithHeight(3)
	if sp.Width() != 4 {
		t.Fatalf("Width: got %d", sp.Width())
	}
	if sp.Height() != 3 {
		t.Fatalf("Height: got %d", sp.Height())
	}
	renderOK(t, sp)
	renderOK(t, sp.AutoScale())
	renderOK(t, sp.ASCII())
	renderOK(t, sp.Unicode())
	renderOK(t, sp.WithBarStyle(func(i, v int) tui.Style { return tui.FrRed }))
}
