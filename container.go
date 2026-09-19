package tui

import "github.com/romanSPB15/tui-compose/v4/cell"

// VBox — вертикальный компоновщик виджетов.
// Добавлено в TUI 3.0.0.
type VBox struct {
	children  []Widget
	positions []Pos
	gap       int
}

// NewVBox создаёт VBox с указанным содержимым.
func NewVBox(children ...Widget) *VBox {
	v := &VBox{}
	v.children = append(v.children, children...)
	return v
}

// SetGap устанавливает щель между виджетами в символах.
func (v *VBox) SetGap(gap int) {
	v.gap = gap
}

// WithGap устанавливает щель между виджетами в символах.
// Возращает тот же VBox для создания цепочки методов.
func (v *VBox) WithGap(gap int) *VBox {
	v.SetGap(gap)
	return v
}

// Add добавляет новый виджет в VBox.
func (v *VBox) Add(widgets ...Widget) {
	v.children = append(v.children, widgets...)
}

func (v *VBox) layout() {
	v.positions = make([]Pos, len(v.children))
	line := 0
	for i := range v.children {
		v.positions[i] = Pos{Line: line, Col: 0}
		if v.children[i] != nil {
			line += v.children[i].Height() + v.gap
		}
	}
}

func (v *VBox) Render([][]cell.Cell) {
}

func (v *VBox) Width() int {
	max := 0
	for _, child := range v.children {
		if child != nil {
			if child.Width() > max {
				max = child.Width()
			}
		}
	}
	return max
}

func (v *VBox) Height() int {
	total := 0
	for _, child := range v.children {
		if child != nil {
			total += child.Height() + v.gap
		}
	}
	return max(total-v.gap, 0)
}

func (v *VBox) Child() []Widget {
	v.layout()
	return v.children
}

func (v *VBox) Pos(i int) Pos {
	if v.positions == nil || len(v.positions) != len(v.children) {
		v.layout()
	}
	if i < 0 || i >= len(v.positions) {
		return Pos{}
	}

	return v.positions[i]
}

// HBox — горизонтальный компоновщик виджетов.
// Добавлено в TUI 3.0.0.
type HBox struct {
	children  []Widget
	positions []Pos
	gap       int
}

// NewHBox создаёт HBox с указанным содержимым.
func NewHBox(children ...Widget) *HBox {
	v := &HBox{gap: 1}
	v.children = append(v.children, children...)
	return v
}

// Add добавляет новый виджет в HBox.
func (v *HBox) Add(widgets ...Widget) {
	v.children = append(v.children, widgets...)
}

func (v *HBox) layout() {
	v.positions = make([]Pos, len(v.children))
	col := 0
	for i, child := range v.children {
		v.positions[i] = Pos{Line: 0, Col: col}
		if v.children[i] != nil {
			col += child.Width() + v.gap
		}
	}
}

// SetGap устанавливает щель между виджетами в символах.
// Возращает тот же HBox для создания цепочки методов.
func (v *HBox) SetGap(gap int) {
	v.gap = gap
}

// WithGap устанавливает щель между виджетами в символах.
// Возращает тот же HBox для создания цепочки методов.
func (v *HBox) WithGap(gap int) *HBox {
	v.SetGap(gap)
	return v
}

func (v *HBox) Render([][]cell.Cell) {
}

func (v *HBox) Width() int {
	if len(v.children) == 0 {
		return 0
	}
	total := 0
	for i, child := range v.children {
		if child != nil {
			total += child.Width()
			if i < len(v.children)-1 {
				total += v.gap
			}
		}
	}
	return total
}

func (v *HBox) Height() int {
	max := 0
	for _, child := range v.children {
		if child != nil {
			if h := child.Height(); h > max {
				max = h
			}
		}
	}
	return max
}

func (v *HBox) Child() []Widget {
	v.layout()
	return v.children
}

func (v *HBox) Pos(i int) Pos {
	if v.positions == nil || len(v.positions) != len(v.children) {
		v.layout()
	}
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
