package tui

import "github.com/eiannone/keyboard"

// Widget — это интерфейс для любого TUI-виджета.
type Widget interface {
	InnerText() string // InnerText() возращает текст виджета
	SetIndex(int)      // SetIndex() устанавливает индекс компонента для перерисовки только его.

	MaxLength() int // MaxLength() возращает длину текста виджета без учёта ANSI Escape последовательностей
	DisplayMode() DisplayMode
}

// App — это объект приложения.
type App interface {
	Widgets() []Widget    // Widgets() возвращает список компонентов, добавленных в приложение.
	AddWidgets(...Widget) // AddWidgets() добавляет компонент в приложение.
	Clear()               // Clear() очищает список компонентов приложения без перерисовки.

	Redraw()          // Redraw() перерисовывает все компоненты. Важно: такая перерисовка вызывает мерцание.
	RedrawWidget(int) // RedrawWidget() перерисовывает конкретный компонент. index — это номер компонента, который нужно перерисовать.

	Run()           // Run() — это блокирующий запуск TUI-приложения. Если пользователь закроет окно, то будет произведён graceful shutdown и выход из метода.
	IsRunned() bool // IsRunned() возращает true, если приложение запущено. Иначе возвращает false.

	Quit()                   // Quit() — это принудительный выход из приложения.
	OnQuit() <-chan struct{} // Run() возвращает канал сигнализации о выходе.

	Window() Window // Window() возвращает интерфейс окна приложения. Из него можно получить длину и ширину окна в символах.

	AddKeyHandler(key keyboard.Key, h func()) // AddKeyHandler() регистрирует обработчик нажатия указанной клавиши

	LogInfo(message string, args ...any)  // LogInfo() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug.
	LogFatal(message string, args ...any) // LogFatal() логирует указанное сообщение подобно fmt.Printf() в файл, если приложение создано как Debug. Потом в любом случае выходит

	Do(f func())
	DoAndWait(f func())
}

// Window — это объект окна приложения.
type Window interface {
	Width() int  // Ширина окна в символах
	Height() int // Высота окна в символах
}
