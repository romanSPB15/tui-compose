package tui

import "strings"

func linesCount(widget Widget) int {
	text := widget.InnerText()
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return len(strings.Split(text, "\n"))
}

// Добавлено в TUI 3.0.0.
type VBox struct {
	children  []Widget
	positions []Pos
}

func NewVBox(children ...Widget) *VBox {
	v := &VBox{}
	v.children = append(v.children, children...)
	v.layout()
	return v
}

func (v *VBox) Add(widgets ...Widget) {
	v.children = append(v.children, widgets...)
	v.layout()
}

func (v *VBox) layout() {
	v.positions = make([]Pos, len(v.children))
	line := 0
	for i := range v.children {
		v.positions[i] = Pos{Line: line, Col: 0}
		line += linesCount(v.children[i])
	}
}

func (v *VBox) InnerText() string { return "" }

func (v *VBox) MaxWidth() int {
	max := 0
	for _, child := range v.children {
		if child.MaxWidth() > max {
			max = child.MaxWidth()
		}
	}
	return max
}
func (v *VBox) MaxHeight() int {
	total := 0
	for _, child := range v.children {
		total += child.MaxHeight()
	}
	return total
}

func (v *VBox) Child() []Widget {
	return v.children
}

func (v *VBox) Pos(i int) Pos {
	if i < 0 || i >= len(v.positions) {
		return Pos{}
	}
	return v.positions[i]
}

// Добавлено в TUI 3.0.0.
type HBox struct {
	children  []Widget
	positions []Pos
	gap       int
}

func NewHBox(children ...Widget) *HBox {
	v := &HBox{gap: 1}
	v.children = append(v.children, children...)
	v.layout()
	return v
}

func (v *HBox) Add(widgets ...Widget) {
	v.children = append(v.children, widgets...)
	v.layout()
}

func (v *HBox) layout() {
	v.positions = make([]Pos, len(v.children))
	col := 0
	for i, child := range v.children {
		v.positions[i] = Pos{Line: 0, Col: col}
		col += child.MaxWidth() + v.gap
	}
}

func (v *HBox) SetGap(gap int) {
	v.gap = gap
	v.layout()
}

func (v *HBox) WithGap(gap int) *HBox {
	v.SetGap(gap)
	return v
}

func (v *HBox) InnerText() string { return "" }

func (v *HBox) MaxWidth() int {
	if len(v.children) == 0 {
		return 0
	}
	total := 0
	for i, child := range v.children {
		total += child.MaxWidth()
		if i < len(v.children)-1 {
			total += v.gap
		}
	}
	return total
}

func (v *HBox) MaxHeight() int {
	max := 0
	for _, child := range v.children {
		if h := child.MaxHeight(); h > max {
			max = h
		}
	}
	return max
}

func (v *HBox) Child() []Widget {
	return v.children
}

func (v *HBox) Pos(i int) Pos {
	if i < 0 || i >= len(v.positions) {
		return Pos{}
	}
	return v.positions[i]
}

func init() {
	var _ Widget = (*VBox)(nil)
	var _ Container = (*VBox)(nil)
	var _ Widget = (*HBox)(nil)
	var _ Container = (*HBox)(nil)
}
