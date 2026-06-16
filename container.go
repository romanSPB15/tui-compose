package tui

import "strings"

func linesCount(widget Widget) int {
	text := widget.InnerText()
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return len(strings.Split(text, "\n"))
}

type VBox struct {
	children  []Widget
	positions []Pos
	idx       int
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

func (v *VBox) MaxLength() int {
	max := 0
	for _, child := range v.children {
		if child.MaxLength() > max {
			max = child.MaxLength()
		}
	}
	return max
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
