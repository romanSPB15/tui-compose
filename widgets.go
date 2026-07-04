//go:build !no_widgets

package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/romanSPB15/tui-compose/v3/input"
)

// Label — это виджет текстовой метки.
type Label struct {
	ANSI  string // Приставка ANSI escape последовательности
	Text  string // Текст виджета.
	len   int
	Block bool // Отображение в блочном режиме.
}

func (l *Label) InnerText() string {
	if l.ANSI == "" {
		return l.Text + strings.Repeat(" ", l.len-len([]rune(l.Text)))
	}
	r := []rune(l.Text)
	if len(r) > l.len {
		r = r[:l.len]
	}
	l.Text = string(r)
	return fmt.Sprintf("%s%s\033[0m", l.ANSI, l.Text+strings.Repeat(" ", l.len-len([]rune(l.Text))))
}

// NewStaticLabel() создаёт виджет текста.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, len: utf8.RuneCountInString(txt)} }

// NewDynamicLabel() создаёт виджет текста с возможностью изменения содержимого в будущем.
// MaxWidth это место, зарезервированное под метку в символах.
func NewDynamicLabel(txt string, len int) *Label {
	return &Label{Text: txt, len: len}
}

// Deprecated: используйте WithStyle(tui.Fr*) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// ColorizeForeground() окрашивает текст в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForeground(clr Color) *Label {
	lbl.ANSI += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

// Deprecated: используйте WithStyle(tui.Bg*) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// ColorizeBackground() окрашивает фон текста в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackground(clr Color) *Label {
	lbl.ANSI += fmt.Sprintf("\033[%dm", clr+10)
	return lbl
}

// ColorizeForegroundRGB() окрашивает текст в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForegroundRGB(clr ColorRGB) *Label {
	lbl.ANSI += fmt.Sprintf("\033[38;2;%d;%d;%dm", clr.R, clr.G, clr.B)
	return lbl
}

// ColorizeBackgroundRGB() окрашивает фон текста в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackgroundRGB(clr ColorRGB) *Label {
	lbl.ANSI += fmt.Sprintf("\033[48;2;%d;%d;%dm", clr.R, clr.G, clr.B)
	return lbl
}

// Deprecated: используйте WithStyle(tui.Bold) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Bold() делает текст жирным.
func (lbl *Label) Bold() *Label {
	lbl.ANSI += "\033[1m"
	return lbl
}

// Deprecated: используйте WithStyle(tui.Italic) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Italic() делает текст курсивом.
func (lbl *Label) Italic() *Label {
	lbl.ANSI += "\033[3m"
	return lbl
}

// Deprecated: используйте WithStyle(tui.Underline) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Underline() подчеркивает текст.
func (lbl *Label) Underline() *Label {
	lbl.ANSI += "\033[4m"
	return lbl
}

// Deprecated: используйте WithStyle(tui.Reverse) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Reverse() реверсирует цвет текста.
func (lbl *Label) Reverse() *Label {
	lbl.ANSI += "\033[7m"
	return lbl
}

// Deprecated: используйте WithStyle(0) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Reset() убирает все декорации текста.
// Добавлено в TUI v1.5.0
func (lbl *Label) Reset() *Label {
	lbl.ANSI = ""
	return lbl
}

// WithStyle() применяет стиль к тексту.
// Добавлено в TUI v3.1.0
func (lbl *Label) WithStyle(s Style) *Label {
	lbl.ANSI = s.String()
	return lbl
}

// Deprecated: используйте WithStyle(tui.Blink) вместо этого метода.
// Будет удалён в версии v4.0.0.
//
// Blink() делает текст мигающим(работает не во всем терминалах).
// Добавлено в TUI v1.1.0
func (lbl *Label) Blink() *Label {
	lbl.ANSI += "\033[5m"
	return lbl
}

// MaxWidth() реализует интерфейс Widget
func (lbl *Label) MaxWidth() int {
	return lbl.len
}

// MaxHeight() реализует интерфейс Widget
// Добавлено в TUI v3.0.0
func (l *Label) MaxHeight() int {
	return 1
}

// SetText() устанавливает текст метки и перерисовывает окно.
// Добавлено в TUI v3.0.0
func (l *Label) SetText(new string) {
	l.Text = new
	currentWindow.Redraw()
}

// Button это виджет кнопки.
type Button struct {
	text          string
	OnClicked     func()
	style, styleF Style
	focused       bool
}

// NewButton() создаёт кнопку.
func NewButton(text string, h func()) *Button {
	btn := &Button{
		text:      text,
		OnClicked: h,
		styleF:    BgWhite | FrBlack,
	}
	return btn
}

func (btn *Button) OnFocus() {
	btn.focused = true
	currentWindow.Redraw()
}

func (btn *Button) OnBlur() {
	btn.focused = false
	currentWindow.Redraw()
}

func (btn *Button) OnClick() {
	if btn.OnClicked != nil {
		btn.OnClicked()
	}
}

func (btn *Button) InnerText() string {
	if btn.focused {
		return btn.styleF.String() + btn.text + Reset.String()
	} else {
		return btn.style.String() + btn.text + Reset.String()
	}
}

func (btn *Button) MaxWidth() int {
	return len([]rune(btn.text))
}

func (*Button) MaxHeight() int {
	return 1
}

// WithStyle устанавливает стиль кнопки когда не в фокусе.
// Добавлено в TUI v3.1.0
func (btn *Button) WithStyle(s Style) *Button {
	btn.style = s
	return btn
}

// WithFocusedStyle устанавливает стиль кнопки в фокусе.
// Добавлено в TUI v3.1.0
func (btn *Button) WithFocusedStyle(s Style) *Button {
	btn.styleF = s
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

// ColorProgress — это виджет шкалы прогресса.
// Добавлено в TUI v1.2.0
type ColorProgress struct {
	text          string
	size          int
	clrOn, clrOff Color
}

// NewColorProgress() cоздаёт виджет шкалы прогресса в виде цветных пикселей.
// len — максимальная длина в пикселях.
// on — цвет "включенных" пикселей.
// off — цвет "выключенных" пикселей.
// Добавлено в TUI v1.2.0
func NewColorProgress(len int, on, off Color) *ColorProgress {
	return &ColorProgress{
		text:   strings.Repeat(" ", len),
		size:   len,
		clrOn:  on,
		clrOff: off,
	}
}

// SetValue() устанавливает значение прогресса. Диапазон 0-1.
// Добавлено в TUI v1.2.0
func (p *ColorProgress) SetValue(f float64) {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	on := int(float64(p.size) * f)
	p.text = fmt.Sprintf("\033[%dm%s\033[%dm%s\033[0m", p.clrOn+10, strings.Repeat(" ", on), p.clrOff+10, strings.Repeat(" ", p.size-on))
	if currentWindow.IsRunned() {
		currentWindow.Redraw()
	}
}

func (p *ColorProgress) MaxWidth() int {
	return p.size
}

func (l *ColorProgress) MaxHeight() int {
	return 1
}

func (p *ColorProgress) InnerText() string {
	return p.text
}

// TextProgress — это виджет шкалы прогресса.
// Добавлено в TUI v2.0.0
type TextProgress struct {
	text      string
	size      int
	sOn, sOff rune
}

// NewTextProgress() cоздаёт виджет шкалы прогресса в виде текста.
// len — максимальная длина в пикселях.
// on — символ "включенных" пикселей.
// off — символ "выключенных" пикселей.
// Добавлено в TUI v2.0.0
func NewTextProgress(len int, on, off rune) *TextProgress {
	return &TextProgress{
		text: strings.Repeat(" ", len),
		size: len,
		sOn:  on,
		sOff: off,
	}
}

// SetValue() устанавливает значение прогресса. Диапазон 0-1.
// Добавлено в TUI v2.0.0
func (p *TextProgress) SetValue(f float64) {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	on := int(float64(p.size) * f)
	p.text = fmt.Sprintf("%s%s", strings.Repeat(string(p.sOn), on), strings.Repeat(string(p.sOff), p.size-on))
	if currentWindow.IsRunned() {
		currentWindow.Redraw()
	}
}

func (p *TextProgress) MaxWidth() int {
	return p.size
}

func (l *TextProgress) MaxHeight() int {
	return 1
}

func (p *TextProgress) InnerText() string {
	return p.text
}

// Check — виджет чекбокса.
// Вызов OnChanged происходит при изменении состояния (после переключения).
// Добавлено в TUI v1.0.0
type Check struct {
	text         string
	checkedState bool
	focused      bool
	OnChanged    func(bool)

	style  Style // обычное состояние
	styleF Style // состояние фокуса
	styleC Style // состояние "включён"
}

// NewCheck создаёт чекбокс с указанным текстом.
func NewCheck(text string) *Check {
	return &Check{
		text:   text,
		styleF: BgWhite | FrBlack, // Добавлено в TUI v3.1.0
		styleC: FrGreen,           // Добавлено в TUI v3.1.0
	}
}

// InnerText реализует интерфейс Widget.
func (c *Check) InnerText() string {
	if c.checkedState && c.focused {
		return c.styleF.String() + "[x] " + c.text + Reset.String()
	}
	if c.focused {
		return c.styleF.String() + "[ ] " + c.text + Reset.String()
	}
	if c.checkedState {
		return c.styleC.String() + "[x] " + c.text + Reset.String()
	}
	if c.style != 0 {
		return c.style.String() + "[ ] " + c.text + Reset.String()
	}
	return "[ ] " + c.text
}

// OnFocus реализует интерфейс Focusable.
// Добавлено в TUI v2.0.0
func (c *Check) OnFocus() {
	c.focused = true
	currentWindow.Redraw()
}

// OnBlur реализует интерфейс Focusable.
// Добавлено в TUI v2.0.0
func (c *Check) OnBlur() {
	c.focused = false
	currentWindow.Redraw()
}

// OnClick реализует интерфейс Clickable.
// Добавлено в TUI v3.0.0
func (c *Check) OnClick() {
	c.checkedState = !c.checkedState
	currentWindow.Redraw()
	if c.OnChanged != nil {
		c.OnChanged(c.checkedState)
	}
}

// WithStyle устанавливает стиль чекбокса (не в фокусе).
// Добавлено в TUI v3.1.0
func (c *Check) WithStyle(s Style) *Check {
	c.style = s
	return c
}

// WithFocusedStyle устанавливает стиль чекбокса в фокусе.
// Добавлено в TUI v3.1.0
func (c *Check) WithFocusedStyle(s Style) *Check {
	c.styleF = s
	return c
}

// WithCheckedStyle устанавливает стиль включенного чекбокса.
// Добавлено в TUI v3.1.0
func (c *Check) WithCheckedStyle(s Style) *Check {
	c.styleC = s
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

// MaxWidth реализует интерфейс Widget.
// Добавлено в TUI v1.0.0
func (c *Check) MaxWidth() int {
	return len([]rune("[x] " + c.text))
}

// MaxHeight реализует интерфейс Widget.
// Добавлено в TUI v3.0.0
func (c *Check) MaxHeight() int {
	return 1
}

// State возвращает текущее состояние чекбокса.
// Добавлено в TUI v1.0.0
func (c *Check) State() bool {
	return c.checkedState
}

// SetState устанавливает состояние чекбокса.
// Добавлено в TUI v1.0.0
func (c *Check) SetState(b bool) {
	c.checkedState = b
	if currentWindow.IsRunned() {
		currentWindow.Redraw()
	}
}

// WithState устанавливает состояние чекбокса и возращает его.
// Добавлено в TUI v3.1.0
func (c *Check) WithState(b bool) *Check {
	c.checkedState = b
	if currentWindow.IsRunned() {
		currentWindow.Redraw()
	}
	return c
}

// InputField — однострочное поле ввода.
// Добавлено в TUI v3.0.0
type InputField struct {
	Text      string
	CursorPos int
	width     int
	focused   bool

	style            Style // обычное состояние
	styleF           Style // состояние фокуса
	cursorStyle      Style // стиль курсора
	placeholder      string
	placeholderStyle Style

	OnChanged func(string)
	OnEnter   func(string)
}

func NewInputField(width int) *InputField {
	return &InputField{
		width:            width,
		styleF:           BgWhite | FrBlack, // стандартный стиль фокуса
		cursorStyle:      Reverse,           // инверсия по умолчанию
		placeholderStyle: FrBrightBlack,     // серый текст для плейсхолдера
	}
}

// WithStyle устанавливает стиль поля (не в фокусе).
// Добавлено в TUI v3.1.0
func (f *InputField) WithStyle(s Style) *InputField {
	f.style = s
	return f
}

// WithFocusedStyle устанавливает стиль поля в фокусе.
// Добавлено в TUI v3.1.0
func (f *InputField) WithFocusedStyle(s Style) *InputField {
	f.styleF = s
	return f
}

// WithCursorStyle устанавливает стиль курсора.
// Добавлено в TUI v3.1.0
func (f *InputField) WithCursorStyle(s Style) *InputField {
	f.cursorStyle = s
	return f
}

// WithText устанавливает текст поля.
// Добавлено в TUI v3.1.0
func (f *InputField) WithText(text string) *InputField {
	f.Text = text
	f.CursorPos = len([]rune(text))
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
	f.placeholderStyle = s
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

func (f *InputField) InnerText() string {
	style := f.style
	if f.focused {
		style = f.styleF
	}

	if f.Text == "" && f.placeholder != "" && !f.focused {
		ps := f.placeholderStyle
		return ps.String() + f.placeholder + Reset.String()
	}

	runes := []rune(f.Text)
	cursor := f.CursorPos
	cursor = max(cursor, 0)
	cursor = min(cursor, len(runes))

	var builder strings.Builder

	if !f.focused || len(runes) == 0 {
		return style.String() + f.Text + Reset.String()
	}

	// Текст до курсора (с общим стилем)
	if cursor > 0 {
		builder.WriteString(style.String() + string(runes[:cursor]))
	}

	// Курсор
	if cursor < len(runes) {
		cursorDisplay := f.cursorStyle.String() + string(runes[cursor]) + Reset.String()
		builder.WriteString(cursorDisplay)
		if cursor+1 < len(runes) {
			builder.WriteString(style.String() + string(runes[cursor+1:]) + Reset.String())
		}
	} else {
		cursorDisplay := f.cursorStyle.String() + " " + Reset.String()
		builder.WriteString(cursorDisplay)
	}

	currentLen := len([]rune(builder.String()))
	if currentLen < f.width {
		builder.WriteString(strings.Repeat(" ", f.width-currentLen))
	}

	return builder.String()
}

func (f *InputField) MaxWidth() int {
	return f.width
}

func (f *InputField) MaxHeight() int {
	return 1
}

// OnFocus реализует интерфейс Focusable.
// Добавлено в TUI v3.0.0
func (f *InputField) OnFocus() {
	f.focused = true
	currentWindow.Redraw()
}

// OnBlur реализует интерфейс Focusable.
// Добавлено в TUI v3.0.0
func (f *InputField) OnBlur() {
	f.focused = false
	currentWindow.Redraw()
}

// OnKeyPress реализует интерфейс TextInput.
// Добавлено в TUI v3.0.0
func (f *InputField) OnKeyPress(ev *input.KeyboardEvent) {
	runes := []rune(f.Text)
	switch ev.Key {
	case input.KeyDelete:
		if f.CursorPos < len(runes) {
			runes = append(runes[:f.CursorPos], runes[f.CursorPos+1:]...)
			f.Text = string(runes)
			currentWindow.Redraw()
		}
		if f.OnChanged != nil {
			f.OnChanged(f.Text)
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
			// Вставка символа
			runes = append(runes[:f.CursorPos], append([]rune{ev.Rune}, runes[f.CursorPos:]...)...)
			f.Text = string(runes)
			f.CursorPos++
			currentWindow.Redraw()
			if f.OnChanged != nil {
				f.OnChanged(f.Text)
			}
		}
	}
}

func init() {
	var _ Widget = (*Label)(nil)
	var _ Widget = (*Button)(nil)
	var _ Focusable = (*Button)(nil)
	var _ Clickable = (*Button)(nil)
	var _ Widget = (*ColorProgress)(nil)
	var _ Widget = (*TextProgress)(nil)
	var _ Widget = (*Check)(nil)
	var _ Focusable = (*Check)(nil)
	var _ Clickable = (*Check)(nil)
	var _ TextInput = (*InputField)(nil)
}
