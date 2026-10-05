package tui

import "github.com/romanSPB15/acell"

type flexItem struct {
	widget Widget
	weight int
}

func FlexItem(w Widget, weight int) flexItem {
	return flexItem{widget: w, weight: weight}
}

type Flex struct {
	items      []flexItem
	horizontal bool
	gap        int
	positions  []Pos
	width      int
	height     int
}

func NewFlex(items ...flexItem) *Flex {
	return &Flex{items: items}
}

func (f *Flex) Horizontal() *Flex {
	f.horizontal = true
	return f
}

func (f *Flex) Vertical() *Flex {
	f.horizontal = false
	return f
}

func (f *Flex) WithGap(g int) *Flex {
	f.gap = g
	return f
}

func (f *Flex) Width() int  { return f.width }
func (f *Flex) Height() int { return f.height }

func (f *Flex) Render([][]acell.Cell) {}

func (f *Flex) Child() []Widget {
	ws := make([]Widget, len(f.items))
	for i, it := range f.items {
		ws[i] = it.widget
	}
	return ws
}

func (f *Flex) Pos(i int) Pos {
	if i < 0 || i >= len(f.positions) {
		return Pos{}
	}
	return f.positions[i]
}

func (f *Flex) Send(ev Event) {
	e, ok := ev.(*MeasureEvent)
	if !ok {
		return
	}
	f.positions = make([]Pos, len(f.items))
	if f.horizontal {
		f.measureH(e.MaxWidth, e.MaxHeight)
	} else {
		f.measureV(e.MaxWidth, e.MaxHeight)
	}
}

func (f *Flex) measureV(maxW, maxH int) {
	f.width = maxW
	f.height = maxH

	n := len(f.items)
	if n == 0 {
		return
	}

	totalGap := f.gap * (n - 1)
	avail := max(maxH-totalGap, 0)

	intrinsic := 0
	totalWeight := 0
	for _, it := range f.items {
		if it.weight == 0 {
			if evh, ok := it.widget.(EventHandler); ok {
				evh.Send(&MeasureEvent{MaxWidth: maxW, MaxHeight: avail})
			}
			intrinsic += it.widget.Height()
		} else {
			totalWeight += it.weight
		}
	}

	remaining := max(avail-intrinsic, 0)

	y := 0
	for i, it := range f.items {
		f.positions[i] = Pos{Line: y, Col: 0}

		h := it.widget.Height()
		if it.weight > 0 && totalWeight > 0 {
			h = remaining * it.weight / totalWeight
			if evh, ok := it.widget.(EventHandler); ok {
				evh.Send(&MeasureEvent{MaxWidth: maxW, MaxHeight: h})
			}
			h = max(min(it.widget.Height(), remaining), 0)
		}
		y += h + f.gap
	}
}

func (f *Flex) measureH(maxW, maxH int) {
	f.width = maxW
	f.height = maxH

	n := len(f.items)
	if n == 0 {
		return
	}

	totalGap := f.gap * (n - 1)
	avail := max(maxW-totalGap, 0)

	intrinsic := 0
	totalWeight := 0
	for _, it := range f.items {
		if it.weight == 0 {
			if evh, ok := it.widget.(EventHandler); ok {
				evh.Send(&MeasureEvent{MaxWidth: avail, MaxHeight: maxH})
			}
			intrinsic += it.widget.Width()
		} else {
			totalWeight += it.weight
		}
	}

	remaining := max(avail-intrinsic, 0)

	x := 0
	for i, it := range f.items {
		f.positions[i] = Pos{Line: 0, Col: x}

		w := it.widget.Width()
		if it.weight > 0 && totalWeight > 0 {
			w = remaining * it.weight / totalWeight
			if evh, ok := it.widget.(EventHandler); ok {
				evh.Send(&MeasureEvent{MaxWidth: w, MaxHeight: maxH})
			}
			w = max(min(it.widget.Width(), remaining), 0)
		}
		x += w + f.gap
	}
}

var _ Widget = (*Flex)(nil)
var _ Container = (*Flex)(nil)
var _ EventHandler = (*Flex)(nil)
