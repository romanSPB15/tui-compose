package tui

// Widget — это интерфейс для любого TUI-виджета или контейнера. Но в контейнере эти методы не используются.
type Widget interface {
	InnerText() string // InnerText() возращает текст виджета

	MaxWidth() int // MaxWidth() возращает длину текста виджета без учёта ANSI Escape последовательностей
	MaxHeight() int
}

// Focusable это интерфейс виджетов, которые могут получить фокус.
// Переключение проиходит с помощью ← и →.
// Отключить переключение можно через Window.DisableFocusChange().
// Добавлено в TUI 2.0.0.
type Focusable interface {
	Widget
	OnFocus()
	OnBlur()
}

// Clickable это интерфейс виджетов, которые могут получить фокус и быть нажатыми.
// Добавлено в TUI 2.0.0.
type Clickable interface {
	Focusable
	OnClick()
}

// Window — это объект окна приложения.
type Window interface {
	Widgets() []Widget    // Widgets() возвращает список компонентов, добавленных в приложение.
	AddWidgets(...Widget) // AddWidgets() добавляет компонент в приложение.
	Clear()               // Clear() очищает список компонентов приложения без перерисовки.

	Redraw()          // Redraw() перерисовывает все компоненты. Важно: такая перерисовка вызывает мерцание.
	RedrawWidget(int) // RedrawWidget() перерисовывает конкретный компонент. index — это номер компонента, который нужно перерисовать.

	Run()           // Run() — это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
	IsRunned() bool // IsRunned() возращает true, если приложение запущено. Иначе возвращает false.

	Quit()                   // Quit() — это принудительный выход из приложения.
	OnQuit() <-chan struct{} // Run() возвращает канал сигнализации о выходе.

	RegisterKeyHandler(key Key, h func())        // RegisterKeyHandler() регистрирует обработчик нажатия указанной клавиши
	RegisterClickHandler(h func(ev *MouseEvent)) // RegisterClickHandler() регистрирует обрабочик событий мыши

	LogInfo(message string, args ...any)  // LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
	LogFatal(message string, args ...any) // LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит

	Do(f func())
	DoAndWait(f func())

	Width() int  // Ширина окна в символах
	Height() int // Высота окна в символах

	DisableFocusChange() // DisableFocusChange() выключает смену фокуса.

	SetContent(Widget)
}

// Container это интерфейс контейнеров.
// Добавлено в TUI 3.0.0.
type Container interface {
	Widget
	Child() []Widget
	Pos(int) Pos
	LineCount() int
}
