// Библитека TUI для Go
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Цвет
type Color int

// Обычные цвета
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Яркие цвета(работают не во всем терминалах)
const (
	BrightBlack Color = iota + 90
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

type DisplayMode int

const (
	DisplayInline  DisplayMode = iota // В одну строку
	DisplayBlock                      // На отдельной строке
	DisplayNewLine                    // Перенос строки
)

// Текст. Может быть декорирован
type Label struct {
	decoration string
	Text       string
	maxLength  int
	Block      bool
}

func (l *Label) innerText() string {
	if l.decoration == "" {
		return l.Text
	}
	return fmt.Sprintf("%s%s\033[0m", l.decoration, l.Text)
}

// Создание объекта текста без возможности изменения.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, maxLength: len(txt)} }

// Создание объекта текста с возможностью изменения содержимого в будущем.
// maxLength это место, зарезервированное под метку в символах.
func NewDynamicLabel(txt string, maxLength int) *Label {
	return &Label{Text: txt, maxLength: maxLength}
}

// Окрасить текст в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForeground(clr Color) *Label {
	lbl.decoration += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

// Окрасить фон текста в один из стандартных цветов.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackground(clr Color) *Label {
	lbl.decoration += fmt.Sprintf("\033[%dm", clr+10)
	return lbl
}

// Окрасить текст в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeForegroundRGB(r, g, b uint8) *Label {
	lbl.decoration += fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	return lbl
}

// Окрасить фон текста в RGB.
// Добавлено в TUI v1.1.0
func (lbl *Label) ColorizeBackgroundRGB(r, g, b uint8) *Label {
	lbl.decoration += fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	return lbl
}

// Сделать текст жирным.
func (lbl *Label) Bold() *Label {
	lbl.decoration += "\033[1m"
	return lbl
}

// Сделать текст курсивом.
func (lbl *Label) Italic() *Label {
	lbl.decoration += "\033[3m"
	return lbl
}

// Подчеркнуть текст.
func (lbl *Label) Underline() *Label {
	lbl.decoration += "\033[4m"
	return lbl
}

// Реверсировать цвет текста.
func (lbl *Label) Reverse() *Label {
	lbl.decoration += "\033[7m"
	return lbl
}

// Сделать текст мигающим(работает не во всем терминалах).
// Добавлено в TUI v1.1.0
func (lbl *Label) Blink() *Label {
	lbl.decoration += "\033[7m"
	return lbl
}

// Реализация tui.Component.
func (lbl *Label) MaxWidth() int {
	return lbl.maxLength
}

// Реализация tui.Component.
func (lbl *Label) DisplayMode() DisplayMode {
	if lbl.Block {
		return DisplayBlock
	}
	return DisplayInline
}

// Реализация tui.Component
func (l *Label) setIndex(int) {}

/////////////////////////////////////////////////////////////////////////////////

type spaser struct{}

func (l *spaser) innerText() string {
	return ""
}

func (l *spaser) DisplayMode() DisplayMode {
	return DisplayBlock
}

func (l *spaser) MaxWidth() int {
	return 0
}

func (l *spaser) setIndex(int) {}

// Создаёт пустой компонент, занимающий всю строку для визуального разделения.
func Spaser() Component { return &spaser{} }

/////////////////////////////////////////////////////////////////////////////////

type newLine struct{}

func (l *newLine) innerText() string {
	return ""
}

func (l *newLine) DisplayMode() DisplayMode {
	return DisplayNewLine
}

func (l *newLine) MaxWidth() int {
	return 0
}

func (l *newLine) setIndex(int) {}

// Создаёт компонент, занимающий весь остаток его строки, что переносит следующие на новую строку.
func NewLine() Component { return &newLine{} }

// Объект кнопки, нажимающейся от нажатия её клавиши. Обработчик в OnClick.
type Button struct {
	clicked Component
	base    Component
	OnClick func()
	Component
	idx int
}

// Создаёт tui.Button
func NewButton(text string, key keyboard.Key) *Button {
	btn := &Button{
		clicked: NewStaticLabel(text).ColorizeForeground(Blue),
		base:    NewStaticLabel(text),
		OnClick: func() {},
	}
	btn.Component = btn.base
	currentApp.AddKeyHandler(key, func() {
		if btn.OnClick != nil {
			btn.OnClick()
		}
		btn.Component = btn.clicked
		currentApp.RedrawComponent(btn.idx)
		currentApp.LogInfo("%d %s", btn.idx, btn.innerText())
		time.Sleep(time.Millisecond * 50)
		btn.Component = btn.base
		currentApp.LogInfo("%d %s", btn.idx, btn.innerText())
		currentApp.RedrawComponent(btn.idx)
	})
	return btn
}

func (btn *Button) setIndex(idx int) {
	btn.idx = idx
	btn.base.setIndex(idx)
	btn.clicked.setIndex(idx)
}

type ColorProgress struct {
	base          Label
	size          int
	clrOn, clrOff Color
	idx           int
}

func (p *ColorProgress) SetValue(f float64) {
	on := int(float64(p.size) * f)
	p.base.Text = fmt.Sprintf("\033[%dm%s\033[%dm%s\033[0m", p.clrOn+10, strings.Repeat(" ", on), p.clrOff+10, strings.Repeat(" ", p.size-on))
	currentApp.RedrawComponent(p.idx)
}

func (p *ColorProgress) setIndex(idx int) {
	p.idx = idx
	p.base.setIndex(idx)
}

func NewColorProgress(len int, on, off Color) *ColorProgress {
	return &ColorProgress{
		base:   *NewDynamicLabel(strings.Repeat(" ", len), len),
		size:   len,
		clrOn:  on,
		clrOff: off,
	}
}
