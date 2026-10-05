package extra

import (
	"github.com/romanSPB15/acell"
	"github.com/romanSPB15/tui-compose/v4"
)

// ScrollView — контейнер с прокруткой. Обрезает дочерний виджет
// по своим габаритам, сдвигая его на (offsetX, offsetY).
type ScrollView struct {
	child tui.Widget

	offsetX int
	offsetY int

	width  int
	height int

	fixedW int
	fixedH int

	wnd tui.Window
}

// NewScrollView создаёт прокручиваемый контейнер.
func NewScrollView(child tui.Widget) *ScrollView {
	return &ScrollView{child: child}
}

// SetChild заменяет дочерний виджет.
func (sv *ScrollView) SetChild(w tui.Widget) {
	sv.child = w
}

// WithFixedWidth задаёт фиксированную ширину окна.
// 0 — автоматически (по содержимому или MeasureEvent).
func (sv *ScrollView) WithFixedWidth(w int) *ScrollView {
	sv.fixedW = w
	return sv
}

// WithFixedHeight задаёт фиксированную высоту окна.
// 0 — автоматически (по содержимому или MeasureEvent).
func (sv *ScrollView) WithFixedHeight(h int) *ScrollView {
	sv.fixedH = h
	return sv
}

// WithFixedSize задаёт и ширину, и высоту окна.
func (sv *ScrollView) WithFixedSize(w, h int) *ScrollView {
	sv.fixedW = w
	sv.fixedH = h
	return sv
}

// Offset возвращает текущее смещение.
func (sv *ScrollView) Offset() (x, y int) {
	return sv.offsetX, sv.offsetY
}

// SetOffset задаёт смещение с ограничением по размеру содержимого.
func (sv *ScrollView) SetOffset(x, y int) {
	sv.offsetX, sv.offsetY = x, y
	sv.clamp()
}

// ScrollBy сдвигает содержимое на (dx, dy).
func (sv *ScrollView) ScrollBy(dx, dy int) {
	sv.offsetX += dx
	sv.offsetY += dy
	sv.clamp()
}

// Width возвращает ширину окна.
// Приоритет: fixedW → width (из MeasureEvent) → intrinsic ребёнка.
func (sv *ScrollView) Width() int {
	if sv.fixedW > 0 {
		return sv.fixedW
	}
	if sv.width > 0 {
		return sv.width
	}
	if sv.child == nil {
		return 0
	}
	return sv.child.Width()
}

// Height возвращает высоту окна.
// Приоритет: fixedH → height (из MeasureEvent) → intrinsic ребёнка.
func (sv *ScrollView) Height() int {
	if sv.fixedH > 0 {
		return sv.fixedH
	}
	if sv.height > 0 {
		return sv.height
	}
	if sv.child == nil {
		return 0
	}
	return sv.child.Height()
}

func (sv *ScrollView) childW() int {
	if sv.child == nil {
		return 0
	}
	return sv.child.Width()
}

func (sv *ScrollView) childH() int {
	if sv.child == nil {
		return 0
	}
	return sv.child.Height()
}

func (sv *ScrollView) clamp() {
	maxX := sv.childW() - sv.Width()
	if maxX < 0 {
		maxX = 0
	}
	maxY := sv.childH() - sv.Height()
	if maxY < 0 {
		maxY = 0
	}
	if sv.offsetX < 0 {
		sv.offsetX = 0
	}
	if sv.offsetX > maxX {
		sv.offsetX = maxX
	}
	if sv.offsetY < 0 {
		sv.offsetY = 0
	}
	if sv.offsetY > maxY {
		sv.offsetY = maxY
	}
}

// Child возвращает дочерний виджет.
func (sv *ScrollView) Child() []tui.Widget {
	if sv.child == nil {
		return nil
	}
	return []tui.Widget{sv.child}
}

// Pos возвращает позицию ребёнка с учётом прокрутки.
func (sv *ScrollView) Pos(i int) tui.Pos {
	return tui.Pos{Line: -sv.offsetY, Col: -sv.offsetX}
}

// Render не используется — ScrollView контейнер, рендерит ядро.
func (sv *ScrollView) Render([][]acell.Cell) {}

// Send обрабатывает MeasureEvent, клавиши и мышь.
func (sv *ScrollView) Send(ev tui.Event) {
	switch e := ev.(type) {
	case *tui.WindowEvent:
		sv.wnd = e.Window
	case *tui.MeasureEvent:
		if sv.fixedW > 0 {
			sv.width = sv.fixedW
		} else {
			sv.width = e.MaxWidth
		}
		if sv.fixedH > 0 {
			sv.height = sv.fixedH
		} else {
			sv.height = e.MaxHeight
		}

		if evh, ok := sv.child.(tui.EventHandler); ok {
			evh.Send(&tui.MeasureEvent{
				MaxWidth:  sv.Width(),
				MaxHeight: 1 << 20,
			})
		}
		sv.clamp()
	case *acell.KeyboardEvent:
		sv.handleKey(e)
	case *acell.MouseEvent:
		sv.handleMouse(e)
	}
}

func (sv *ScrollView) handleKey(e *acell.KeyboardEvent) {
	switch e.Key {
	case acell.KeyArrowDown:
		sv.ScrollBy(0, 1)
	case acell.KeyArrowUp:
		sv.ScrollBy(0, -1)
	case acell.KeyArrowRight:
		sv.ScrollBy(1, 0)
	case acell.KeyArrowLeft:
		sv.ScrollBy(-1, 0)
	case acell.KeyPgDown:
		sv.ScrollBy(0, sv.Height())
	case acell.KeyPgUp:
		sv.ScrollBy(0, -sv.Height())
	case acell.KeyHome:
		sv.SetOffset(sv.offsetX, 0)
	case acell.KeyEnd:
		sv.SetOffset(sv.offsetX, sv.childH())
	default:
		return
	}
	sv.redraw()
}

func (sv *ScrollView) handleMouse(e *acell.MouseEvent) {
	switch e.Action {
	case acell.MouseWheelDown:
		sv.ScrollBy(0, 3)
	case acell.MouseWheelUp:
		sv.ScrollBy(0, -3)
	case acell.MouseWheelRight:
		sv.ScrollBy(3, 0)
	case acell.MouseWheelLeft:
		sv.ScrollBy(-3, 0)
	default:
		return
	}
	sv.redraw()
}

func (sv *ScrollView) redraw() {
	if sv.wnd == nil {
		return
	}
	sv.wnd.Index()
	sv.wnd.Redraw()
}

var _ tui.Widget = (*ScrollView)(nil)
var _ tui.Container = (*ScrollView)(nil)
var _ tui.EventHandler = (*ScrollView)(nil)
