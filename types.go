package tui

import "github.com/romanSPB15/tui-compose/v3/input"

// Widget — это интерфейс для любого TUI-виджета или контейнера. Но в контейнере эти методы не используются.
type Widget interface {
	InnerText() string // InnerText() возращает текст виджета

	MaxWidth() int // MaxWidth() возращает длину текста виджета без учёта ANSI Escape последовательностей
	MaxHeight() int
}

// Focusable это интерфейс виджетов, которые могут получить фокус.
// Переключение проиходит с помощью TAB и →.
// Отключить переключение можно через Window.DisableFocusChange().
// Добавлено в TUI 2.0.0.
type Focusable interface {
	Widget
	OnFocus()
	OnBlur()
}

// Clickable это интерфейс виджетов, которые могут быть нажатыми мышью.
// Добавлено в TUI 2.0.0.
type Clickable interface {
	Widget
	OnClick()
}

// ClickableAt это интерфейс виджетов, которые могут быть нажаты мышью.
// Добавлено в TUI 3.0.0.
type ClickableAt interface {
	Widget
	OnClickAt(x, y int)
}

// Добавлено в TUI 3.0.0.
type KeyReceiver interface {
	Focusable
	OnKeyPress(ev *input.KeyboardEvent)
}

// Disablable — интерфейс для виджетов, которые могут быть отключены.
// Добавлено в TUI 3.1.0.
type Disablable interface {
	SetDisabled(bool)
	IsDisabled() bool
}

// Window — это объект приложения.
type Window interface {
	SetContent(Widget) // SetContent() устанавливает содержимое окна.

	Redraw() // Redraw() перерисовывает все виджеты.

	Run()           // Run() — это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
	IsRunned() bool // IsRunned() возращает true, если приложение запущено. Иначе возвращает false.

	Quit()                   // Quit() — это принудительный выход из приложения.
	OnQuit() <-chan struct{} // Run() возвращает канал сигнализации о выходе.

	RegisterKeyHandler(KeyboardEventHandler)           // RegisterKeyHandler() регистрирует обработчик нажатия указанной клавиши
	RegisterClickHandler(h func(ev *input.MouseEvent)) // RegisterClickHandler() регистрирует обрабочик событий мыши

	LogInfo(message string, args ...any)  // LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
	LogFatal(message string, args ...any) // LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит

	Do(f func())
	DoAndWait(f func())

	Width() int  // Ширина окна в символах
	Height() int // Высота окна в символах

	SetTitle(title string)
	CopyToClipboard(text string)

	SetOverlay(wgt Widget)
	ShowOverlay()
	HideOverlay()

	Focus() FocusManager
}

// FocusManager — интерфейс менеджера фокуса.
// Добавлено в TUI 3.1.0.
type FocusManager interface {
	FocusedWidget() Focusable // FocusedWidget() вовзращает виджет, на котором установлен фокус.
	NextFocus()               // NextFocus() переносит фокус дальше.
	BeforeFocus()             // BeforeFocus() переносит фокус назад.
	SetFocus(Focusable) bool  // BeforeFocus() устанавливает фокус на переданный виджет.
	ClearFocus()              // ClearFocus() сбрасывает фокус.
	Disable()                 // Disable() отключает смену фокуса.
}

// Container это интерфейс контейнеров.
// Добавлено в TUI 3.0.0.
type Container interface {
	Widget
	Child() []Widget
	Pos(int) Pos
}
