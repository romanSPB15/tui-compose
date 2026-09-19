//go:build !no_widgets

package tui

import (
	"fmt"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v4/builder"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
	"github.com/romanSPB15/tui-compose/v4/term"
)

// DisableState хранит состояние disabled.
type DisableState struct {
	disabled bool
}

// SetDisabled устанавливает состояние disabled.
func (d *DisableState) SetDisabled(v bool) {
	d.disabled = v
}

// IsDisabled возвращает true, если виджет отключён.
func (d *DisableState) IsDisabled() bool {
	return d.disabled
}

// Label — текстовая метка.
type Label struct {
	style cell.Style
	Text  string
	len   int
}

// Render рисует текст в буфер.
func (l *Label) Render(buf [][]cell.Cell) {
	runes := []rune(l.Text)
	for i := range utf8.RuneCountInString(l.Text) {
		buf[0][i] = cell.Cell{Char: runes[i], Style: l.style}
	}
}

// NewStaticLabel создаёт метку с шириной по тексту.
func NewStaticLabel(txt string) *Label {
	return &Label{Text: txt, len: utf8.RuneCountInString(txt)}
}

// NewDynamicLabel создаёт метку с зарезервированной шириной len.
func NewDynamicLabel(txt string, len int) *Label {
	return &Label{Text: txt, len: len}
}

// WithStyle применяет стиль.
func (lbl *Label) WithStyle(s Style) *Label {
	lbl.style = ConvertToCellStyle(s)
	return lbl
}

// ColorizeBackgroundRGB устанавливает цвет фона в RGB.
func (lbl *Label) ColorizeBackgroundRGB(clr ColorRGB) *Label {
	lbl.style.Bg = fmt.Sprintf("48;2;%d;%d;%d", clr.R, clr.G, clr.B)
	return lbl
}

// ColorizeForegroundRGB устанавливает цвет текста в RGB.
func (lbl *Label) ColorizeForegroundRGB(clr ColorRGB) *Label {
	lbl.style.Fg = fmt.Sprintf("38;2;%d;%d;%d", clr.R, clr.G, clr.B)
	return lbl
}

// Width возвращает ширину метки.
func (lbl *Label) Width() int {
	return lbl.len
}

// Height возвращает высоту метки.
func (l *Label) Height() int {
	return 1
}

// SetText устанавливает текст.
func (l *Label) SetText(new string) {
	l.Text = new
}

// WithText устанавливает текст и возвращает метку.
func (l *Label) WithText(new string) *Label {
	l.Text = new
	return l
}

// Button — виджет кнопки.
type Button struct {
	text                  string
	OnClicked             func()
	style, styleF, styleH cell.Style
	styleD                cell.Style
	focused               bool
	hovered               bool
	paddingH, paddingV    int
	DisableState
	wnd Window
}

// NewButton создаёт кнопку с текстом и обработчиком нажатия.
func NewButton(text string, h func()) *Button {
	whiteBg := cell.Style{Fg: "30", Bg: "47"}
	return &Button{
		text:      text,
		OnClicked: h,
		style:     whiteBg,
		styleF:    whiteBg,
		styleH:    whiteBg,
		styleD:    cell.Style{Fg: "90"},
		paddingH:  2,
	}
}

// Render рисует кнопку в буфер.
func (btn *Button) Render(buf [][]cell.Cell) {
	var s cell.Style
	switch {
	case btn.IsDisabled():
		s = btn.styleD
	case btn.focused:
		s = btn.styleF
	case btn.hovered && btn.styleH != (cell.Style{}):
		s = btn.styleH
	default:
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
	textStartX := (w - len(textRunes)) / 2
	textY := (h - 1) / 2

	for i, r := range textRunes {
		buf[textY][textStartX+i] = cell.Cell{Char: r, Style: s}
	}
}

// Width возвращает ширину кнопки.
func (btn *Button) Width() int {
	return utf8.RuneCountInString(btn.text) + 2*btn.paddingH
}

// Height возвращает высоту кнопки.
func (btn *Button) Height() int {
	return 1 + 2*btn.paddingV
}

// GetText возвращает текст кнопки.
func (btn *Button) GetText() string {
	return btn.text
}

// WithPaddings задаёт внутренние отступы.
func (btn *Button) WithPaddings(h, v int) *Button {
	btn.paddingH = h
	btn.paddingV = v
	return btn
}

// WithStyle задаёт стиль кнопки без фокуса.
func (btn *Button) WithStyle(s Style) *Button {
	btn.style = ConvertToCellStyle(s)
	return btn
}

// WithFocusedStyle задаёт стиль кнопки в фокусе.
func (btn *Button) WithFocusedStyle(s Style) *Button {
	btn.styleF = ConvertToCellStyle(s)
	return btn
}

// WithHoverStyle задаёт стиль кнопки при наведении.
func (btn *Button) WithHoverStyle(s Style) *Button {
	btn.styleH = ConvertToCellStyle(s)
	return btn
}

// WithDisabledStyle задаёт стиль отключённой кнопки.
func (btn *Button) WithDisabledStyle(s Style) *Button {
	btn.styleD = ConvertToCellStyle(s)
	return btn
}

// WithText устанавливает текст кнопки.
func (btn *Button) WithText(text string) *Button {
	btn.text = text
	return btn
}

// WithHandler устанавливает обработчик нажатия.
func (btn *Button) WithHandler(h func()) *Button {
	btn.OnClicked = h
	return btn
}

// Send обрабатывает событие.
func (btn *Button) Send(ev Event) {
	switch ev := ev.(type) {
	case *WindowEvent:
		btn.wnd = ev.Window
	case *CheckFocusableEvent:
		ev.Result = !btn.IsDisabled()
	case *FocusEvent:
		btn.focused = ev.Focused
	case *MouseHoverEvent:
		btn.hovered = ev.Entered
		if btn.wnd != nil {
			btn.wnd.Redraw()
		}
	case *input.MouseEvent:
		if ev.Action == input.MousePress {
			if btn.OnClicked != nil {
				btn.OnClicked()
			}
			if !btn.focused {
				btn.Send(FocusEvent{true})
			}
		}
	case *input.KeyboardEvent:
		if (ev.Key == input.KeyEnter || ev.Key == input.KeySpace) && btn.OnClicked != nil {
			btn.OnClicked()
		}
	}
}

// Check — виджет чекбокса.
type Check struct {
	text         string
	checkedState bool
	focused      bool
	hovered      bool
	OnChanged    func(bool)

	style, styleF, styleC, styleH cell.Style
	wnd                           Window
}

// NewCheck создаёт чекбокс с указанным текстом.
func NewCheck(text string) *Check {
	return &Check{
		text:   text,
		styleF: cell.Style{Fg: "30", Bg: "47"},
		styleC: cell.Style{Fg: "32"},
		styleH: cell.Style{Fg: "30", Bg: "47"},
	}
}

// Render рисует чекбокс в буфер.
func (c *Check) Render(buf [][]cell.Cell) {
	var s cell.Style
	switch {
	case c.focused:
		s = c.styleF
	case c.hovered && c.styleH != (cell.Style{}):
		s = c.styleH
	case c.checkedState:
		s = c.styleC
	default:
		s = c.style
	}

	buf[0][0] = cell.Cell{Char: '[', Style: s}
	if c.checkedState {
		buf[0][1] = cell.Cell{Char: 'x', Style: s}
	} else {
		buf[0][1] = cell.Cell{Char: ' ', Style: s}
	}
	buf[0][2] = cell.Cell{Char: ']', Style: s}
	buf[0][3] = cell.Cell{Char: ' ', Style: s}

	runes := []rune(c.text)
	for i := range utf8.RuneCountInString(c.text) {
		buf[0][i+4] = cell.Cell{Char: runes[i], Style: s}
	}
}

// Width возвращает ширину чекбокса.
func (c *Check) Width() int {
	return utf8.RuneCountInString(c.text) + 4
}

// Height возвращает высоту чекбокса.
func (c *Check) Height() int {
	return 1
}

// State возвращает текущее состояние.
func (c *Check) State() bool {
	return c.checkedState
}

// SetState устанавливает состояние.
func (c *Check) SetState(b bool) {
	c.checkedState = b
}

// WithState устанавливает состояние и возвращает чекбокс.
func (c *Check) WithState(b bool) *Check {
	c.checkedState = b
	return c
}

// WithStyle задаёт стиль чекбокса.
func (c *Check) WithStyle(s Style) *Check {
	c.style = ConvertToCellStyle(s)
	return c
}

// WithFocusedStyle задаёт стиль чекбокса в фокусе.
func (c *Check) WithFocusedStyle(s Style) *Check {
	c.styleF = ConvertToCellStyle(s)
	return c
}

// WithHoverStyle задаёт стиль чекбокса при наведении.
func (c *Check) WithHoverStyle(s Style) *Check {
	c.styleH = ConvertToCellStyle(s)
	return c
}

// WithCheckedStyle задаёт стиль включённого чекбокса.
func (c *Check) WithCheckedStyle(s Style) *Check {
	c.styleC = ConvertToCellStyle(s)
	return c
}

// WithText устанавливает текст чекбокса.
func (c *Check) WithText(text string) *Check {
	c.text = text
	return c
}

// WithOnChanged задаёт обработчик изменения состояния.
func (c *Check) WithOnChanged(h func(bool)) *Check {
	c.OnChanged = h
	return c
}

// Send обрабатывает событие.
func (c *Check) Send(ev Event) {
	switch ev := ev.(type) {
	case *WindowEvent:
		c.wnd = ev.Window
	case *CheckFocusableEvent:
		ev.Result = true
	case *FocusEvent:
		c.focused = ev.Focused
	case *MouseHoverEvent:
		c.hovered = ev.Entered
		if c.wnd != nil {
			c.wnd.Redraw()
		}
	case *input.MouseEvent:
		if ev.Action != input.MousePress {
			return
		}
		c.checkedState = !c.checkedState
		if c.wnd != nil {
			c.wnd.Redraw()
		}
		if c.OnChanged != nil {
			c.OnChanged(c.checkedState)
		}
	case *input.KeyboardEvent:
		if ev.Key == input.KeyEnter || ev.Key == input.KeySpace {
			c.checkedState = !c.checkedState
			if c.wnd != nil {
				c.wnd.Redraw()
			}
			if c.OnChanged != nil {
				c.OnChanged(c.checkedState)
			}
		}
	}
}

// InputField — однострочное поле ввода.
type InputField struct {
	Text      string
	CursorPos int
	width     int

	focused   bool
	hovered   bool
	overwrite bool

	style, styleF, styleH cell.Style
	cursorStyle           cell.Style
	placeholder           string
	placeholderStyle      cell.Style

	OnChanged func(string)
	OnEnter   func(string)

	wnd Window

	blinkEnabled  bool
	blinkInterval time.Duration
	cursorVisible bool
	blinkOnce     sync.Once
}

// NewInputField создаёт поле ввода заданной ширины.
func NewInputField(width int) *InputField {
	return &InputField{
		width:            width,
		style:            cell.Style{Bg: "44"},
		cursorStyle:      cell.Style{Bg: "47", Fg: "34"},
		placeholderStyle: cell.Style{Fg: "90"},
		blinkEnabled:     true,
		blinkInterval:    500 * time.Millisecond,
		cursorVisible:    true,
	}
}

// WithStyle задаёт стиль поля без фокуса.
func (f *InputField) WithStyle(s Style) *InputField {
	f.style = ConvertToCellStyle(s)
	return f
}

// WithFocusedStyle задаёт стиль поля в фокусе.
func (f *InputField) WithFocusedStyle(s Style) *InputField {
	f.styleF = ConvertToCellStyle(s)
	return f
}

// WithHoverStyle задаёт стиль поля при наведении.
func (f *InputField) WithHoverStyle(s Style) *InputField {
	f.styleH = ConvertToCellStyle(s)
	return f
}

// WithCursorStyle задаёт стиль курсора.
func (f *InputField) WithCursorStyle(s Style) *InputField {
	f.cursorStyle = ConvertToCellStyle(s)
	return f
}

// WithText устанавливает текст поля.
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

// WithPlaceholder задаёт текст-подсказку.
func (f *InputField) WithPlaceholder(text string) *InputField {
	f.placeholder = text
	return f
}

// WithPlaceholderStyle задаёт стиль подсказки.
func (f *InputField) WithPlaceholderStyle(s Style) *InputField {
	f.placeholderStyle = ConvertToCellStyle(s)
	return f
}

// WithOnChanged задаёт обработчик изменения текста.
func (f *InputField) WithOnChanged(h func(string)) *InputField {
	f.OnChanged = h
	return f
}

// WithOnEnter задаёт обработчик нажатия Enter.
func (f *InputField) WithOnEnter(h func(string)) *InputField {
	f.OnEnter = h
	return f
}

// WithBlinkInterval задаёт интервал мигания курсора.
func (f *InputField) WithBlinkInterval(d time.Duration) *InputField {
	f.blinkInterval = d
	return f
}

// DisableBlink отключает мигание и оставляет курсор постоянно видимым.
func (f *InputField) DisableBlink() *InputField {
	f.blinkEnabled = false
	f.cursorVisible = true
	return f
}

// EnableBlink включает мигание курсора.
func (f *InputField) EnableBlink() *InputField {
	f.blinkEnabled = true
	return f
}

// WithOverwrite включает или выключает режим затирания.
func (f *InputField) WithOverwrite(v bool) *InputField {
	f.overwrite = v
	return f
}

// IsOverwrite возвращает true, если включён режим затирания.
func (f *InputField) IsOverwrite() bool {
	return f.overwrite
}

// Render рисует поле ввода в буфер.
func (f *InputField) Render(buf [][]cell.Cell) {
	fieldStyle := f.style
	switch {
	case f.focused:
		fieldStyle = f.styleF
	case f.hovered && f.styleH != (cell.Style{}):
		fieldStyle = f.styleH
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

	if f.focused && f.cursorVisible {
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

// Width возвращает ширину поля.
func (f *InputField) Width() int {
	return f.width
}

// Height возвращает высоту поля.
func (f *InputField) Height() int {
	return 1
}

// Send обрабатывает событие.
func (f *InputField) Send(ev Event) {
	switch ev := ev.(type) {
	case *WindowEvent:
		f.wnd = ev.Window
	case *CheckFocusableEvent:
		ev.Result = true
	case *FocusEvent:
		f.focused = ev.Focused
		if ev.Focused {
			f.cursorVisible = true
			f.startBlink()
		}
	case *MouseHoverEvent:
		f.hovered = ev.Entered
		if f.wnd != nil {
			f.wnd.Redraw()
		}
	case *input.KeyboardEvent:
		runes := []rune(f.Text)
		switch ev.Key {
		case input.KeyDelete:
			if f.CursorPos < len(runes) {
				runes = append(runes[:f.CursorPos], runes[f.CursorPos+1:]...)
				f.Text = string(runes)
				f.cursorVisible = true
				if f.wnd != nil {
					f.wnd.Redraw()
				}
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
			f.cursorVisible = true
			if f.wnd != nil {
				f.wnd.Redraw()
			}
			if f.OnChanged != nil {
				f.OnChanged(f.Text)
			}
		case input.KeyArrowRight:
			if f.CursorPos < len(runes) {
				f.CursorPos++
				f.cursorVisible = true
				if f.wnd != nil {
					f.wnd.Redraw()
				}
			}
		case input.KeyArrowLeft:
			if f.CursorPos > 0 {
				f.CursorPos--
				f.cursorVisible = true
				if f.wnd != nil {
					f.wnd.Redraw()
				}
			}
		case input.KeyEnter:
			if f.OnEnter != nil {
				f.OnEnter(f.Text)
			}
		case input.KeyInsert:
			f.overwrite = !f.overwrite
			f.cursorVisible = true
			if f.wnd != nil {
				f.wnd.Redraw()
			}
		default:
			if ev.Rune != 0 {
				if f.width > 0 && utf8.RuneCountInString(f.Text) >= f.width && f.CursorPos >= len(runes) {
					return
				}
				if f.overwrite && f.CursorPos < len(runes) {
					// затираем символ под курсором
					runes[f.CursorPos] = ev.Rune
					f.Text = string(runes)
					f.CursorPos++
				} else {
					// вставка со сдвигом
					if f.width > 0 && utf8.RuneCountInString(f.Text) >= f.width {
						return
					}
					runes = append(runes[:f.CursorPos], append([]rune{ev.Rune}, runes[f.CursorPos:]...)...)
					f.Text = string(runes)
					f.CursorPos++
				}
				f.cursorVisible = true
				if f.wnd != nil {
					f.wnd.Redraw()
				}
				if f.OnChanged != nil {
					f.OnChanged(f.Text)
				}
			}
		}
	case *input.MouseEvent:
		if ev.Action == input.MousePress {
			runes := []rune(f.Text)
			pos := ev.Pos.X
			if pos > len(runes) {
				pos = len(runes)
			}
			if pos < 0 {
				pos = 0
			}
			f.CursorPos = pos
			f.focused = true
			f.cursorVisible = true
			f.startBlink()
			if f.wnd != nil {
				f.wnd.Focus().SetFocus(f)
				f.wnd.Redraw()
			}
		}
	}
}

// startBlink запускает цикл мигания курсора один раз за жизнь виджета.
func (f *InputField) startBlink() {
	if f.wnd == nil {
		return
	}
	f.blinkOnce.Do(func() {
		go f.blinkLoop()
	})
}

// blinkLoop переключает видимость курсора с заданным интервалом,
// пока виджет в фокусе и мигание включено.
func (f *InputField) blinkLoop() {
	ticker := time.NewTicker(f.blinkInterval)
	defer ticker.Stop()

	quit := f.wnd.OnQuit()
	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			f.wnd.Do(func() {
				if f.focused && f.blinkEnabled {
					f.cursorVisible = !f.cursorVisible
					f.wnd.Redraw()
				}
			})
		}
	}
}

func init() {
	var _ Widget = (*Label)(nil)
	var _ EventHandler = (*Button)(nil)
	var _ EventHandler = (*Check)(nil)
	var _ EventHandler = (*InputField)(nil)
}

// NewHyperlink создаёт кликабельный текст, открывающий URL.
func NewHyperlink(text, url string) *Button {
	return NewButton(text, func() {
		term.OpenURL(url)
	}).WithPaddings(0, 0)
}

// Gauge — горизонтальная шкала прогресса.
//
// Значение задаётся в [0, 1] через WithValue. Поверх шкалы выводится
// подпись: по умолчанию процент, либо результат LabelFunc.
type Gauge struct {
	value           float64
	size            int
	cellOn, cellOff cell.Cell
	LabelFunc       func(float64) string
	labelStyle      cell.Style
}

// NewGauge создаёт шкалу заданной ширины (минимум 4).
func NewGauge(size int) *Gauge {
	if size < 4 {
		size = 4
	}
	return &Gauge{
		size:    size,
		cellOn:  cell.Cell{Char: '∎', Style: ConvertToCellStyle(FrBlue)},
		cellOff: cell.Cell{Char: '∎', Style: ConvertToCellStyle(FrBrightBlack)},
	}
}

// WithValue устанавливает значение шкалы, приводимое к [0, 1].
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

// Width возвращает ширину шкалы.
func (p *Gauge) Width() int {
	return p.size
}

// Height возвращает высоту шкалы.
func (p *Gauge) Height() int {
	return 1
}

// Render рисует шкалу в буфер.
func (p *Gauge) Render(cells [][]cell.Cell) {
	var i int
	for ; i < int(math.Round(p.value*float64(p.size))); i++ {
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

	label := b.String()
	i = p.size/2 - utf8.RuneCountInString(label)/2
	for j, r := range []rune(label) {
		x := i + j
		if x < 0 || x >= p.size {
			continue
		}
		cells[0][x] = cell.Cell{Char: r, Style: p.labelStyle}
	}
}

// WithOnCell задаёт ячейку заполненной части.
func (p *Gauge) WithOnCell(c cell.Cell) *Gauge {
	p.cellOn = c
	return p
}

// WithOffCell задаёт ячейку пустой части.
func (p *Gauge) WithOffCell(c cell.Cell) *Gauge {
	p.cellOff = c
	return p
}

// WithOnStyle задаёт стиль заполненной части, сохраняя символ.
func (p *Gauge) WithOnStyle(c Style) *Gauge {
	p.cellOn.Style = ConvertToCellStyle(c)
	return p
}

// WithOffStyle задаёт стиль пустой части, сохраняя символ.
func (p *Gauge) WithOffStyle(c Style) *Gauge {
	p.cellOff.Style = ConvertToCellStyle(c)
	return p
}

// WithOnChar задаёт символ заполненной части, сохраняя стиль.
func (p *Gauge) WithOnChar(r rune) *Gauge {
	p.cellOn.Char = r
	return p
}

// WithOffChar задаёт символ пустой части, сохраняя стиль.
func (p *Gauge) WithOffChar(r rune) *Gauge {
	p.cellOff.Char = r
	return p
}

// ASCII переключает на ######----.
func (p *Gauge) ASCII() *Gauge {
	p.WithOnChar('#')
	p.WithOffChar('-')
	return p
}

// Default переключает на ∎∎∎∎∎∎∎∎∎∎.
func (p *Gauge) Default() *Gauge {
	p.WithOnChar('∎')
	p.WithOffChar('∎')
	return p
}

// EmptySquares переключает на ■■■□□□.
func (p *Gauge) EmptySquares() *Gauge {
	p.WithOnChar('■')
	p.WithOffChar('□')
	return p
}

// Blocks переключает на ██████░░░░.
func (p *Gauge) Blocks() *Gauge {
	p.WithOnChar('█')
	p.WithOffChar('░')
	return p
}

// BlocksFull переключает на ██████████.
func (p *Gauge) BlocksFull() *Gauge {
	p.WithOnChar('█')
	p.WithOffChar('█')
	return p
}

// BlocksGrid переключает на ▓▓▓▓▓▓▒▒▒▒.
func (p *Gauge) BlocksGrid() *Gauge {
	p.WithOnChar('▓')
	p.WithOffChar('▒')
	return p
}

// WithLabelFunc задаёт функцию подписи по значению.
func (p *Gauge) WithLabelFunc(fn func(float64) string) *Gauge {
	p.LabelFunc = fn
	return p
}

// WithLabel задаёт постоянную подпись.
func (p *Gauge) WithLabel(lbl string) *Gauge {
	p.LabelFunc = func(float64) string { return lbl }
	return p
}

// WithLabelStyle задаёт стиль подписи.
func (p *Gauge) WithLabelStyle(s Style) *Gauge {
	p.labelStyle = ConvertToCellStyle(s)
	return p
}
