package extra

import (
	"math"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

// Slider — горизонтальный слайдер.
// Поддерживает клик-прыжок, перетаскивание бегунка мышью и
// управление стрелками/Home/End/PgUp/PgDn.
// Добавлено в TUI v4.0.0.
type Slider struct {
	min, max float64
	value    float64
	step     float64

	width int

	trackRune, fillRune, thumbRune rune

	trackStyle, fillStyle, thumbStyle cell.Style
	focusStyle, hoverStyle            cell.Style
	labelStyle                        cell.Style

	focused    bool
	hovered    bool
	dragging   bool
	dragButton int

	labelFunc func(float64) string
	OnChanged func(float64)

	wnd tui.Window
}

// NewSlider создаёт слайдер шириной width (минимум 3).
// Диапазон по умолчанию [0, 1], значение 0, шаг непрерывный.
// Добавлено в TUI v4.0.0.
func NewSlider(width int) *Slider {
	if width < 3 {
		width = 3
	}
	return &Slider{
		min:        0,
		max:        1,
		width:      width,
		trackRune:  '─',
		fillRune:   '━',
		thumbRune:  '●',
		trackStyle: cell.Style{Fg: "90"},
		fillStyle:  cell.Style{Fg: "36"},
		thumbStyle: cell.Style{Fg: "36"},
		focusStyle: cell.Style{Fg: "33", Bg: "236"},
		hoverStyle: cell.Style{Fg: "96"},
		dragButton: -1,
	}
}

// WithRange задаёт диапазон [min, max]. Порядок аргументов не важен.
func (s *Slider) WithRange(min, max float64) *Slider {
	if min > max {
		min, max = max, min
	}
	s.min, s.max = min, max
	s.value = s.clamp(s.value)
	return s
}

// WithValue устанавливает значение (зажимается в [min, max]).
func (s *Slider) WithValue(v float64) *Slider {
	s.value = s.clamp(v)
	return s
}

// WithStep задаёт шаг; 0 — непрерывный.
func (s *Slider) WithStep(step float64) *Slider {
	if step < 0 {
		step = 0
	}
	s.step = step
	return s
}

// WithOnChanged задаёт обработчик изменения значения.
func (s *Slider) WithOnChanged(fn func(float64)) *Slider {
	s.OnChanged = fn
	return s
}

// WithLabelFunc задаёт функцию подписи по значению.
func (s *Slider) WithLabelFunc(fn func(float64) string) *Slider {
	s.labelFunc = fn
	return s
}

func (s *Slider) WithTrackRune(r rune) *Slider { s.trackRune = r; return s }
func (s *Slider) WithFillRune(r rune) *Slider  { s.fillRune = r; return s }
func (s *Slider) WithThumbRune(r rune) *Slider { s.thumbRune = r; return s }

func (s *Slider) WithTrackStyle(st tui.Style) *Slider {
	s.trackStyle = tui.ConvertToCellStyle(st)
	return s
}

func (s *Slider) WithFillStyle(st tui.Style) *Slider {
	s.fillStyle = tui.ConvertToCellStyle(st)
	return s
}

func (s *Slider) WithThumbStyle(st tui.Style) *Slider {
	s.thumbStyle = tui.ConvertToCellStyle(st)
	return s
}

func (s *Slider) WithFocusedStyle(st tui.Style) *Slider {
	s.focusStyle = tui.ConvertToCellStyle(st)
	return s
}

func (s *Slider) WithHoverStyle(st tui.Style) *Slider {
	s.hoverStyle = tui.ConvertToCellStyle(st)
	return s
}

func (s *Slider) WithLabelStyle(st tui.Style) *Slider {
	s.labelStyle = tui.ConvertToCellStyle(st)
	return s
}

// ASCII — пресет для использования ASCII.
func (s *Slider) ASCII() *Slider {
	s.trackRune = '-'
	s.fillRune = '='
	s.thumbRune = 'O'
	return s
}

func (s *Slider) Value() float64 { return s.value }

func (s *Slider) Percent() float64 {
	if s.max == s.min {
		return 0
	}
	return (s.value - s.min) / (s.max - s.min)
}

func (s *Slider) SetValue(v float64) { s.setValue(v) }

func (s *Slider) Width() int { return s.width }

func (s *Slider) Height() int { return 1 }

func (s *Slider) Render(buf [][]cell.Cell) {
	if len(buf) == 0 || len(buf[0]) == 0 {
		return
	}
	w := s.width
	if w > len(buf[0]) {
		w = len(buf[0])
	}
	if w < 1 {
		return
	}

	pos := s.thumbPos()
	thumbStyle := s.thumbStyle
	switch {
	case s.focused && s.focusStyle != (cell.Style{}):
		thumbStyle = s.focusStyle
	case s.hovered && s.hoverStyle != (cell.Style{}):
		thumbStyle = s.hoverStyle
	}

	for x := 0; x < w; x++ {
		switch {
		case x < pos:
			buf[0][x] = cell.Cell{Char: s.fillRune, Style: s.fillStyle}
		case x == pos:
			buf[0][x] = cell.Cell{Char: s.thumbRune, Style: thumbStyle}
		default:
			buf[0][x] = cell.Cell{Char: s.trackRune, Style: s.trackStyle}
		}
	}

	if s.labelFunc != nil {
		runes := []rune(s.labelFunc(s.value))
		start := w/2 - len(runes)/2
		for i, r := range runes {
			x := start + i
			if x < 0 || x >= w {
				continue
			}
			buf[0][x] = cell.Cell{Char: r, Style: s.labelStyle}
		}
	}
}

func (s *Slider) Send(ev tui.Event) {
	switch e := ev.(type) {
	case *tui.WindowEvent:
		s.wnd = e.Window
	case *tui.CheckFocusableEvent:
		e.Result = true
	case *tui.FocusEvent:
		s.focused = e.Focused
		if !e.Focused {
			s.dragging = false
			s.dragButton = -1
		}
	case *tui.MouseHoverEvent:
		s.hovered = e.Entered
		if s.wnd != nil {
			s.wnd.Redraw()
		}
	case *input.KeyboardEvent:
		s.handleKey(e)
	case *input.MouseEvent:
		s.handleMouse(e)
	}
}

func (s *Slider) handleMouse(e *input.MouseEvent) {
	switch e.Action {
	case input.MousePress:
		if e.Button != 0 {
			return
		}
		s.dragging = true
		s.dragButton = e.Button
		s.jumpTo(e.Pos.X)
		if s.wnd != nil {
			s.wnd.Focus().SetFocus(s)
		}

	case input.MouseMove:
		if !s.dragging {
			return
		}
		s.jumpTo(e.Pos.X)

	case input.MouseRelease:
		if !s.dragging {
			return
		}
		if s.dragButton != -1 && e.Button != s.dragButton {
			return
		}
		s.jumpTo(e.Pos.X)
		s.dragging = false
		s.dragButton = -1
	}
}

func (s *Slider) handleKey(e *input.KeyboardEvent) {
	step := s.step
	if step <= 0 {
		step = (s.max - s.min) / 100
	}
	switch e.Key {
	case input.KeyArrowLeft:
		s.setValue(s.value - step)
	case input.KeyArrowRight:
		s.setValue(s.value + step)
	case input.KeyHome:
		s.setValue(s.min)
	case input.KeyEnd:
		s.setValue(s.max)
	case input.KeyPgUp:
		s.setValue(s.value + 10*step)
	case input.KeyPgDown:
		s.setValue(s.value - 10*step)
	default:
		return
	}
	if s.wnd != nil {
		s.wnd.Redraw()
	}
}

func (s *Slider) jumpTo(x int) {
	if s.width < 2 {
		return
	}
	p := float64(x) / float64(s.width-1)
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	if s.setValue(s.min+p*(s.max-s.min)) && s.wnd != nil {
		s.wnd.Redraw()
	}
}

func (s *Slider) thumbPos() int {
	w := s.width
	if w < 1 {
		return 0
	}
	pos := int(math.Round(s.Percent() * float64(w-1)))
	if pos < 0 {
		pos = 0
	}
	if pos > w-1 {
		pos = w - 1
	}
	return pos
}

func (s *Slider) setValue(v float64) bool {
	v = s.clamp(v)
	if s.step > 0 {
		v = s.min + math.Round((v-s.min)/s.step)*s.step
		v = s.clamp(v)
	}
	if v == s.value {
		return false
	}
	s.value = v
	if s.OnChanged != nil {
		s.OnChanged(v)
	}
	return true
}

func (s *Slider) clamp(v float64) float64 {
	if v < s.min {
		return s.min
	}
	if v > s.max {
		return s.max
	}
	return v
}

func init() {
	var _ tui.Widget = (*Slider)(nil)
	var _ tui.EventHandler = (*Slider)(nil)
}
