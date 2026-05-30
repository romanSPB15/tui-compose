package tui

import (
	"fmt"
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

// Яркие цвета
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

// Текст. Может быть покрашен
type Label struct {
	decoration string
	Text       string
	maxLength  int
}

func (l *Label) innerText() string {
	if l.decoration == "" {
		return l.Text
	}
	return fmt.Sprintf("%s%s\033[0m", l.decoration, l.Text)
}

// Создание объекта текста без возможности изменения.
func NewStaticLabel(txt string) *Label { return &Label{Text: txt, maxLength: len(txt)} }

// Создание объекта текста с возможностью изменения.
// maxLength это место, зарезервированное под метку в символах
func NewDynamicLabel(txt string, maxLength int) *Label {
	return &Label{Text: txt, maxLength: maxLength}
}

// Окрасить Label. Возвращает тот же объект.
func (lbl *Label) Colorize(clr Color) *Label {
	lbl.decoration += fmt.Sprintf("\033[%dm", clr)
	return lbl
}

// Сделать жирным.
func (lbl *Label) Bold() *Label {
	lbl.decoration += "\033[1m"
	return lbl
}

// Сделать курсивом.
func (lbl *Label) Italic() *Label {
	lbl.decoration += "\033[3m"
	return lbl
}

// Подчеркнуть.
func (lbl *Label) Underline() *Label {
	lbl.decoration += "\033[4m"
	return lbl
}

// Реверсировать цвет.
func (lbl *Label) Reverse() *Label {
	lbl.decoration += "\033[7m"
	return lbl
}

// Реализация tui.Component
func (lbl *Label) MaxWidth() int {
	return lbl.maxLength
}

// Реализация tui.Component
func (lbl *Label) DisplayMode() DisplayMode {
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

// Создаёт компонент, занимающий всю строку
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

// Создаёт компонент, переносящий следующие на новую строку
func NewLine() Component { return &newLine{} }

// Объект кнопки, нажимающейся от нажатия её клавиши
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
		clicked: NewStaticLabel(text).Colorize(Blue),
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
