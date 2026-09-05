//go:build !no_widgets

package tui

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v4/builder"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
	"github.com/romanSPB15/tui-compose/v4/term"
)

// DisableState хранит состояние disabled и предоставляет методы.
type DisableState struct {
	disabled bool
}

func (d *DisableState) SetDisabled(v bool) {
	d.disabled = v
}

func (d *DisableState) IsDisabled() bool {
	return d.disabled
}

// Label — это виджет текстовой метки.
type Label struct {
	style cell.Style
	Text  string // Текст виджета.
	len   int
}

func (l *Label) Render(buf [][]cell.Cell) {
	runes := []rune(l.Text)
	for i := range utf8.RuneCountInString(l.Text) {
		buf[0][i] = cell.Cell{Char: runes[i], Style: l.style}
	}
}

// NewStaticLabel() создаёт виджет текста.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, len: utf8.RuneCountInString(txt)} }

// NewDynamicLabel() создаёт виджет текста с возможностью изменения содержимого в будущем.
// Width это место, зарезервированное под метку в символах.
func NewDynamicLabel(txt string, len int) *Label {
	return &Label{Text: txt, len: len}
}

// WithStyle применяет стиль к тексту.
// Добавлено в TUI v3.1.0.
func (lbl *Label) WithStyle(s Style) *Label {
	lbl.style = ConvertToCellStyle(s)
	return lbl
}

// ColorizeBackgroundRGB устанавливает цвет фона текста в RGB.
// Добавлено в TUI v1.1.0.
func (lbl *Label) ColorizeBackgroundRGB(clr ColorRGB) *Label {
	lbl.style.Bg = fmt.Sprintf("48;2;%d;%d;%d", clr.R, clr.G, clr.B)
	return lbl
}

// ColorizeForegroundRGB устанавливает цвет текста в RGB.
// Добавлено в TUI v1.1.0.
func (lbl *Label) ColorizeForegroundRGB(clr ColorRGB) *Label {
	lbl.style.Fg = fmt.Sprintf("38;2;%d;%d;%d", clr.R, clr.G, clr.B)
	return lbl
}

// Width() реализует интерфейс Widget
// Добавлено в TUI v4.0.0
func (lbl *Label) Width() int {
	return lbl.len
}

// Height() реализует интерфейс Widget
// Добавлено в TUI v4.0.0
func (l *Label) Height() int {
	return 1
}

// SetText() устанавливает текст метки.
// Добавлено в TUI v3.0.0
func (l *Label) SetText(new string) {
	l.Text = new
}

// WithText() устанавливает текст метки.
// Добавлено в TUI v4.0.0
func (l *Label) WithText(new string) *Label {
	l.Text = new
	return l
}

// Button это виджет кнопки.
type Button struct {
	text                  string
	OnClicked             func()
	style, styleF, styleD cell.Style
	focused               bool
	paddingH, paddingV    int
	DisableState
}

// NewButton() создаёт кнопку.
func NewButton(text string, h func()) *Button {
	btn := &Button{
		text:      text,
		OnClicked: h,
		styleF:    cell.Style{Fg: "30", Bg: "47"},
		styleD:    cell.Style{Fg: "90"},
		paddingH:  2,
		paddingV:  0,
	}
	return btn
}

func (btn *Button) Render(buf [][]cell.Cell) {
	var s cell.Style
	if btn.IsDisabled() {
		s = btn.styleD
	} else if btn.focused {
		s = btn.styleF
	} else {
		s = btn.style
	}

	w := btn.Width()
	h := btn.Height()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			buf[y][x] = cell.Cell{Char: ' ', Style: s}
		}
	}

	textRunes := []rune(btn.text)
	textLen := len(textRunes)
	textStartX := (w - textLen) / 2
	textY := (h - 1) / 2

	for i, r := range textRunes {
		buf[textY][textStartX+i] = cell.Cell{Char: r, Style: s}
	}
}

func (btn *Button) Width() int {
	return utf8.RuneCountInString(btn.text) + 2*btn.paddingH
}

func (btn *Button) Height() int {
	return 1 + 2*btn.paddingV
}

// WithPaddings устанавливает внутренние отступы.
// Добавлено в TUI v3.3.1
func (btn *Button) WithPaddings(h, v int) *Button {
	btn.paddingH = h
	btn.paddingV = v
	return btn
}

// WithStyle устанавливает стиль кнопки когда не в фокусе.
// Добавлено в TUI v3.1.0
func (btn *Button) WithStyle(s Style) *Button {
	btn.style = ConvertToCellStyle(s)
	return btn
}

// WithFocusedStyle устанавливает стиль кнопки в фокусе.
// Добавлено в TUI v3.1.0
func (btn *Button) WithFocusedStyle(s Style) *Button {
	btn.styleF = ConvertToCellStyle(s)
	return btn
}

// WithDisabledStyle устанавливает стиль кнопки в фокусе.
// Добавлено в TUI v3.1.0
func (btn *Button) WithDisabledStyle(s Style) *Button {
	btn.styleD = ConvertToCellStyle(s)
	return btn
}

// WithText устанавливает текст кнопки.
// Добавлено в TUI v3.1.0
func (btn *Button) WithText(text string) *Button {
	btn.text = text
	return btn
}

// WithHandler устанавливает обработчик нажатия.
// Добавлено в TUI v3.1.0
func (btn *Button) WithHandler(h func()) *Button {
	btn.OnClicked = h
	return btn
}

func (btn *Button) Send(ev Event) {
	switch ev := ev.(type) {
	case *CheckFocusableEvent:
		ev.Result = !btn.IsDisabled()
	case *FocusEvent:
		btn.focused = ev.Focused
	case *input.MouseEvent:
		if btn.OnClicked != nil {
			btn.OnClicked()
		}
		if !btn.focused {
			btn.Send(FocusEvent{true})
		}
	case *input.KeyboardEvent:
		if (ev.Key == input.KeyEnter || ev.Key == input.KeySpace) && btn.OnClicked != nil {
			btn.OnClicked()
		}
	}
}

// Check — виджет чекбокса.
// Вызов OnChanged происходит при изменении состояния (после переключения).
// Добавлено в TUI v1.0.0
type Check struct {
	text         string
	checkedState bool
	focused      bool
	OnChanged    func(bool)

	style  cell.Style // обычное состояние
	styleF cell.Style // состояние фокуса
	styleC cell.Style // состояние "включён"
}

func NewCheck(text string) *Check {
	return &Check{
		text:   text,
		styleF: cell.Style{Fg: "30", Bg: "47"},
		styleC: cell.Style{Fg: "32"},
	}
}

// Render реализует интерфейс Widget.
func (c *Check) Render(buf [][]cell.Cell) {
	s := c.style
	if c.focused {
		s = c.styleF
	} else if c.checkedState {
		s = c.styleC
	}
	buf[0][0] = cell.Cell{Char: '[', Style: s}
	if c.checkedState {
		buf[0][1] = cell.Cell{Char: 'x', Style: s}
	}
	buf[0][2] = cell.Cell{Char: ']', Style: s}

	runes := []rune(c.text)
	for i := range utf8.RuneCountInString(c.text) {
		buf[0][i+4] = cell.Cell{Char: runes[i], Style: s}
	}
}

func (c *Check) Width() int {
	return utf8.RuneCountInString(c.text) + 4
}

func (c *Check) Height() int {
	return 1
}

// State возвращает текущее состояние чекбокса.
func (c *Check) State() bool {
	return c.checkedState
}

// SetState устанавливает состояние чекбокса.
func (c *Check) SetState(b bool) {
	c.checkedState = b
}

// WithState устанавливает состояние чекбокса и возращает его.
// Добавлено в TUI v3.1.0
func (c *Check) WithState(b bool) *Check {
	c.checkedState = b
	return c
}

// WithStyle устанавливает стиль чекбокса (не в фокусе).
// Добавлено в TUI v3.1.0
func (c *Check) WithStyle(s Style) *Check {
	c.style = ConvertToCellStyle(s)
	return c
}

// WithFocusedStyle устанавливает стиль чекбокса в фокусе.
// Добавлено в TUI v3.1.0
func (c *Check) WithFocusedStyle(s Style) *Check {
	c.styleF = ConvertToCellStyle(s)
	return c
}

// WithCheckedStyle устанавливает стиль включенного чекбокса.
// Добавлено в TUI v3.1.0
func (c *Check) WithCheckedStyle(s Style) *Check {
	c.styleC = ConvertToCellStyle(s)
	return c
}

// WithText устанавливает текст чекбокса.
// Добавлено в TUI v3.1.0
func (c *Check) WithText(text string) *Check {
	c.text = text
	return c
}

// WithOnChanged устанавливает обработчик изменения состояния.
// Добавлено в TUI v3.1.0
func (c *Check) WithOnChanged(h func(bool)) *Check {
	c.OnChanged = h
	return c
}

func (c *Check) Send(ev Event) {
	switch ev := ev.(type) {
	case *CheckFocusableEvent:
		ev.Result = true
	case *FocusEvent:
		c.focused = ev.Focused
	case *input.MouseEvent:
		c.checkedState = !c.checkedState
		if !c.focused {
			c.Send(FocusEvent{true})
		}
		currentWindow.Redraw()
		if c.OnChanged != nil {
			c.OnChanged(c.checkedState)
		}

	case *input.KeyboardEvent:
		if ev.Key == input.KeyEnter || ev.Key == input.KeySpace {
			c.checkedState = !c.checkedState
			currentWindow.Redraw()
			if c.OnChanged != nil {
				c.OnChanged(c.checkedState)
			}
		}
	}
}

// InputField — однострочное поле ввода.
// Добавлено в TUI v3.0.0
type InputField struct {
	Text      string
	CursorPos int
	width     int
	focused   bool

	style            cell.Style // обычное состояние
	styleF           cell.Style // состояние фокуса
	cursorStyle      cell.Style // стиль курсора
	placeholder      string
	placeholderStyle cell.Style

	OnChanged func(string)
	OnEnter   func(string)
}

func NewInputField(width int) *InputField {
	return &InputField{
		width:            width,
		style:            cell.Style{Bg: "44"},
		cursorStyle:      cell.Style{Bg: "47", Fg: "34"},
		placeholderStyle: cell.Style{Fg: "90"},
	}
}

// WithStyle устанавливает стиль поля (не в фокусе).
// Добавлено в TUI v3.1.0
func (f *InputField) WithStyle(s Style) *InputField {
	f.style = ConvertToCellStyle(s)
	return f
}

// WithFocusedStyle устанавливает стиль поля в фокусе.
// Добавлено в TUI v3.1.0
func (f *InputField) WithFocusedStyle(s Style) *InputField {
	f.styleF = ConvertToCellStyle(s)
	return f
}

// WithCursorStyle устанавливает стиль курсора.
// Добавлено в TUI v3.1.0
func (f *InputField) WithCursorStyle(s Style) *InputField {
	f.cursorStyle = ConvertToCellStyle(s)
	return f
}

// WithText устанавливает текст поля.
// Добавлено в TUI v3.1.0
func (f *InputField) WithText(text string) *InputField {
	l := utf8.RuneCountInString(text)
	if l > f.width {
		l = f.width
		text = string([]rune(text)[:f.width])
	}
	f.Text = text
	f.CursorPos = l
	return f
}

// WithPlaceholder устанавливает текст-подсказку.
// Добавлено в TUI v3.1.0
func (f *InputField) WithPlaceholder(text string) *InputField {
	f.placeholder = text
	return f
}

// WithPlaceholderStyle устанавливает стиль плейсхолдера.
// Добавлено в TUI v3.1.0
func (f *InputField) WithPlaceholderStyle(s Style) *InputField {
	f.placeholderStyle = ConvertToCellStyle(s)
	return f
}

// WithOnChanged устанавливает обработчик изменения текста.
// Добавлено в TUI v3.1.0
func (f *InputField) WithOnChanged(h func(string)) *InputField {
	f.OnChanged = h
	return f
}

// WithOnEnter устанавливает обработчик нажатия Enter.
// Добавлено в TUI v3.1.0
func (f *InputField) WithOnEnter(h func(string)) *InputField {
	f.OnEnter = h
	return f
}

func (f *InputField) Render(buf [][]cell.Cell) {
	fieldStyle := f.style
	if f.focused {
		fieldStyle = f.styleF
	}

	var displayText string
	var textStyle cell.Style
	if !f.focused && f.Text == "" && f.placeholder != "" {
		displayText = f.placeholder
		textStyle = f.placeholderStyle.Merge(f.style)
	} else {
		displayText = f.Text
		textStyle = fieldStyle
	}

	runes := []rune(displayText)
	if len(runes) > f.width {
		runes = runes[:f.width]
	}
	displayText = string(runes)

	for i := 0; i < f.width; i++ {
		var ch rune
		var st cell.Style
		if i < len(runes) {
			ch = runes[i]
			st = textStyle
		} else {
			ch = ' '
			st = fieldStyle
		}
		buf[0][i] = cell.Cell{Char: ch, Style: st}
	}

	if f.focused {
		cursorPos := f.CursorPos

		if cursorPos > len(runes) {
			cursorPos = len(runes)
		}

		if cursorPos < f.width {
			var cursorChar rune
			if cursorPos < len(runes) {
				cursorChar = runes[cursorPos]
			} else {
				cursorChar = ' '
			}
			buf[0][cursorPos] = cell.Cell{Char: cursorChar, Style: f.cursorStyle}
		}
	}
}

func (f *InputField) Width() int {
	return f.width
}

func (f *InputField) Height() int {
	return 1
}

func (f *InputField) Send(ev Event) {
	switch ev := ev.(type) {
	case *CheckFocusableEvent:
		ev.Result = true
	case *FocusEvent:
		f.focused = ev.Focused
	case *input.KeyboardEvent:
		runes := []rune(f.Text)
		switch ev.Key {
		case input.KeyDelete:
			if f.CursorPos < len(runes) {
				runes = append(runes[:f.CursorPos], runes[f.CursorPos+1:]...)
				f.Text = string(runes)
				currentWindow.Redraw()
				if f.OnChanged != nil {
					f.OnChanged(f.Text)
				}
			}
		case input.KeyBackspace:
			if f.CursorPos <= 0 {
				return
			}
			runes = append(runes[:f.CursorPos-1], runes[f.CursorPos:]...)
			f.Text = string(runes)
			f.CursorPos--
			currentWindow.Redraw()
			if f.OnChanged != nil {
				f.OnChanged(f.Text)
			}
		case input.KeyArrowRight:
			if f.CursorPos < len(runes) {
				f.CursorPos++
				currentWindow.Redraw()
			}
		case input.KeyArrowLeft:
			if f.CursorPos > 0 {
				f.CursorPos--
				currentWindow.Redraw()
			}
		case input.KeyEnter:
			if f.OnEnter != nil {
				f.OnEnter(f.Text)
			}
		default:
			if ev.Rune != 0 {
				if f.width > 0 && utf8.RuneCountInString(f.Text) >= f.width {
					return
				}
				runes = append(runes[:f.CursorPos], append([]rune{ev.Rune}, runes[f.CursorPos:]...)...)
				f.Text = string(runes)
				f.CursorPos++
				currentWindow.Redraw()
				if f.OnChanged != nil {
					f.OnChanged(f.Text)
				}
			}
		}
	case *input.MouseEvent:
		if !f.focused {
			f.Send(FocusEvent{true})
		}
	}
}

func init() {
	var _ Widget = (*Label)(nil)
	var _ EventHandler = (*Button)(nil)
	var _ EventHandler = (*Check)(nil)
	var _ EventHandler = (*InputField)(nil)
}

// NewHyperlink создаёт гиперссылку.
func NewHyperlink(text string, url string) *Button {
	return NewButton(text, func() {
		term.OpenURL(url)
	}).WithPaddings(0, 0)
}

// Gauge — это виджет шкалы прогресса.
// Добавлено в TUI v3.4.0
type Gauge struct {
	value           float64
	size            int
	cellOn, cellOff cell.Cell
	LabelFunc       func(float64) string
	labelStyle      cell.Style
}

// NewGauge() cоздаёт Gauge.
func NewGauge(size int) *Gauge {
	if size < 4 {
		size = 4
	}
	return &Gauge{
		size:    size,
		value:   0.0,
		cellOn:  cell.Cell{Char: '∎', Style: ConvertToCellStyle(FrBlue)},
		cellOff: cell.Cell{Char: '∎', Style: ConvertToCellStyle(FrBrightBlack)},
	}
}

// WithValue() устанавливает значение Gauge в диапазоне от 0 до 1.
func (p *Gauge) WithValue(f float64) *Gauge {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	p.value = f
	return p
}

func (p *Gauge) Width() int {
	return p.size
}

func (l *Gauge) Height() int {
	return 1
}

func (p *Gauge) Render(cells [][]cell.Cell) {
	var i int
	for ; i < int(math.Round(float64(p.value)*float64(p.size))); i++ {
		cells[0][i] = p.cellOn
	}
	for ; i < p.size; i++ {
		cells[0][i] = p.cellOff
	}
	var b builder.Builder
	if p.LabelFunc == nil {
		b.WriteFormat("%d%%", int(math.Round(p.value*100)))
	} else {
		b.WriteString(p.LabelFunc(p.value))
	}

	i = p.size/2 - len([]rune(b.String()))/2
	for j, r := range b.String() {
		x := i + j
		if x < 0 || x >= p.size {
			continue
		}
		cells[0][x] = cell.Cell{Char: r, Style: p.labelStyle}
	}
}

func (p *Gauge) WithOnCell(c cell.Cell) *Gauge {
	p.cellOn = c
	return p
}

func (p *Gauge) WithOffCell(c cell.Cell) *Gauge {
	p.cellOff = c
	return p
}

func (p *Gauge) WithOnStyle(c Style) *Gauge {
	p.cellOn.Style = ConvertToCellStyle(c)
	return p
}

func (p *Gauge) WithOffStyle(c Style) *Gauge {
	p.cellOff.Style = ConvertToCellStyle(c)
	return p
}

func (p *Gauge) WithOnChar(r rune) *Gauge {
	p.cellOn.Char = r
	return p
}

func (p *Gauge) WithOffChar(r rune) *Gauge {
	p.cellOff.Char = r
	return p
}

// ASCII устанавливает ASCII-стиль. ######----
func (p *Gauge) ASCII() *Gauge {
	p.WithOnChar('#')
	p.WithOffChar('-')
	return p
}

// Default устанавливает стиль по умолчанию. ∎∎∎∎∎∎∎∎∎∎
func (p *Gauge) Default() *Gauge {
	p.WithOnChar('∎')
	p.WithOffChar('∎')
	return p
}

// EmptySquares устанавливает стиль ∎∎∎∎∎∎□□□□
func (p *Gauge) EmptySquares() *Gauge {
	p.WithOnChar('■')
	p.WithOffChar('□')
	return p
}

// Blocks устанавливает стиль ██████░░░░
func (p *Gauge) Blocks() *Gauge {
	p.WithOnChar('█')
	p.WithOffChar('░')
	return p
}

// BlocksFull устанавливает стиль ██████████
func (p *Gauge) BlocksFull() *Gauge {
	p.WithOnChar('█')
	p.WithOffChar('█')
	return p
}

// BlocksGrid устанавливает стиль ▓▓▓▓▓▓▒▒▒▒
func (p *Gauge) BlocksGrid() *Gauge {
	p.WithOnChar('▓')
	p.WithOffChar('▒')
	return p
}

// WithLabelFunc устанавливает функцию метки.
func (p *Gauge) WithLabelFunc(fn func(float64) string) *Gauge {
	p.LabelFunc = fn
	return p
}

// WithLabel устанавливает текст метки.
func (p *Gauge) WithLabel(lbl string) *Gauge {
	p.LabelFunc = func(f float64) string {
		return lbl
	}
	return p
}

// WithLabel устанавливает стиль метки.
func (p *Gauge) WithLabelStyle(s Style) *Gauge {
	p.labelStyle = ConvertToCellStyle(s)
	return p
}
