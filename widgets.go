//go:build !no_widgets

package tui

import (
	"fmt"
	"strings"
	"time"
)

// Label — это виджет текстовой метки.
type Label struct {
	ANSI      string // Приставка ANSI escape последовательности
	Text      string // Текст виджета.
	maxLength int
	Block     bool // Отображение в блочном режиме.
}

func (l *Label) InnerText() string {
	if l.ANSI == "" {
		return l.Text
	}
	return fmt.Sprintf("%s%s\033[0m", l.ANSI, l.Text)
}

// NewStaticLabel() создаёт виджет текста.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, maxLength: len(txt)} }

// NewDynamicLabel() создаёт виджет текста с возможностью изменения содержимого в будущем.
// maxLength это место, зарезервированное под метку в символах.
func NewDynamicLabel(txt string, maxLength int) *Label {
	return &Label{Text: txt, maxLength: maxLength}
}

// ColorizeForeground() окрашивает текст в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForeground(clr Color) *Label {
	lbl.ANSI += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

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

// Bold() делает текст жирным.
func (lbl *Label) Bold() *Label {
	lbl.ANSI += "\033[1m"
	return lbl
}

// Italic() делает текст курсивом.
func (lbl *Label) Italic() *Label {
	lbl.ANSI += "\033[3m"
	return lbl
}

// Underline() подчеркивает текст.
func (lbl *Label) Underline() *Label {
	lbl.ANSI += "\033[4m"
	return lbl
}

// Reverse() реверсирует цвет текста.
func (lbl *Label) Reverse() *Label {
	lbl.ANSI += "\033[7m"
	return lbl
}

// Reset() убирает все декорации текста.
// Добавлено в TUI v1.5.0
func (lbl *Label) Reset() *Label {
	lbl.ANSI = ""
	return lbl
}

// Blink() делает текст мигающим(работает не во всем терминалах).
// Добавлено в TUI v1.1.0
func (lbl *Label) Blink() *Label {
	lbl.ANSI += "\033[7m"
	return lbl
}

// MaxLength() реализует интерфейс Widget
func (lbl *Label) MaxLength() int {
	return lbl.maxLength
}

// Button это виджет кнопки.
type Button struct {
	clicked   Widget
	selected  Widget
	base      Widget
	OnClicked func()
	Widget
	idx int
}

// NewButton() создаёт кнопку.
func NewButton(text string) *Button {
	btn := &Button{
		clicked:   NewStaticLabel(text).ColorizeForeground(Blue),
		selected:  NewStaticLabel(text).ColorizeBackground(White).ColorizeForeground(Black),
		base:      NewStaticLabel(text),
		OnClicked: func() {},
	}
	btn.Widget = btn.base
	return btn
}

func (btn *Button) OnFocus() {
	btn.Widget = btn.selected
	currentWindow.RedrawWidget(btn.idx)
}

func (btn *Button) OnBlur() {
	btn.Widget = btn.base
	currentWindow.RedrawWidget(btn.idx)
}

func (btn *Button) OnClick() {
	btn.Widget = btn.clicked
	currentWindow.RedrawWidget(btn.idx)
	currentWindow.LogInfo("%d %s", btn.idx, btn.InnerText())
	if btn.OnClicked != nil {
		btn.OnClicked()
	}
	time.Sleep(time.Millisecond * 50)
	btn.Widget = btn.base
	currentWindow.LogInfo("%d %s", btn.idx, btn.InnerText())
	currentWindow.RedrawWidget(btn.idx)
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

func (p *ColorProgress) MaxLength() int {
	return p.size
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

func (p *TextProgress) MaxLength() int {
	return p.size
}

func (p *TextProgress) InnerText() string {
	return p.text
}

// Check — виджет чекбокса.
// Вызов OnChanged происходит при изменении состояния (после переключения).
type Check struct {
	text                    string
	checkedState            bool
	focused                 bool
	OnChanged               func()
	base, checked, selected *Label
}

func NewCheck(text string) *Check {
	c := &Check{
		text:      text,
		OnChanged: func() {},
	}
	c.updateWidgets()
	return c
}

func (c *Check) updateWidgets() {
	unchecked := "[ ] " + c.text
	checked := "[x] " + c.text

	c.base = NewStaticLabel(unchecked)
	c.checked = NewStaticLabel(checked).ColorizeForeground(Blue)
	c.selected = NewStaticLabel(checked).ColorizeBackground(White).ColorizeForeground(Black)
	if !c.checkedState {
		c.selected = NewStaticLabel(unchecked).ColorizeBackground(White).ColorizeForeground(Black)
	}
}

func (c *Check) InnerText() string {
	if c.focused {
		return c.selected.InnerText()
	}
	if c.checkedState {
		return c.checked.InnerText()
	}
	return c.base.InnerText()
}

func (c *Check) OnFocus() {
	c.focused = true
	currentWindow.Redraw()
}

func (c *Check) OnBlur() {
	c.focused = false
	currentWindow.Redraw()
}

func (c *Check) OnClick() {
	c.checkedState = !c.checkedState
	c.updateWidgets()
	currentWindow.Redraw()
	if c.OnChanged != nil {
		c.OnChanged()
	}
}

func (c *Check) MaxLength() int {
	return len("[x] " + c.text)
}

// State() возвращает значение чекбокса.
func (c *Check) State() bool {
	return c.checkedState
}

// SetState() устанавливает значение чекбокса.
func (c *Check) SetState(b bool) {
	c.checkedState = b
	c.updateWidgets()
	if currentWindow.IsRunned() {
		currentWindow.Redraw()
	}
}
