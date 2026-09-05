package tui

import (
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

// Widget — это интерфейс для TUI-виджетов.
type Widget interface {
	Render(buf [][]cell.Cell)

	Width() int
	Height() int
}

type EventHandler interface {
	Widget
	Send(Event)
}

// Window — это объект приложения.
type Window interface {
	SetContent(Widget) // SetContent устанавливает содержимое окна.

	Redraw() // Redraw перерисовывает окно.

	Run()           // Run запускает TUI-приложение. Если пользователь закроет окно(Ctrl+C), то будет произведён graceful shutdown и выход из метода.
	IsRunned() bool // IsRunned возращает, запущено ли приложение.

	Quit()                   // Quit выходит из приложения.
	OnQuit() <-chan struct{} // Run возвращает канал сигнализации о выходе.

	RegisterKeyHandler(KeyboardEventHandler)         // RegisterKeyHandler регистрирует обработчик событий клавиатуры.
	RegisterClickHandler(func(ev *input.MouseEvent)) // RegisterClickHandler регистрирует обрабочик событий мыши.

	LogInfo(message string, args ...any)  // LogInfo логирует сообщение подобно fmt.Printf в файл, если включен debug режим.
	LogFatal(message string, args ...any) // LogFatal логирует сообщение подобно fmt.Printf в файл, если включен debug режим, и завершает приложение.

	Do(func())        // Do отправляет задачу в UI поток.
	DoAndWait(func()) // DoAndWait отправляет задачу в UI поток и блокируется до завершения.

	Width() int  // Width возвращает ширину окна в символах.
	Height() int // Height возвращает высоту окна в символах.

	SetTitle(string)        // SetTitle устанавливает заголовок окна терминала.
	CopyToClipboard(string) // CopyToClipboard копирует текст в буфер обмена.

	SetOverlay(Widget)
	ShowOverlay()
	HideOverlay()

	Focus() FocusManager

	SetInitCell(cell.Cell) // SetInitCell устанавливает ячейку по умолчанию для всех пустых позиций окна.
	SetBackground(Style)   // SetBackground устанавливает стиль пустых позиций окна.

	Index() // Index обновляет кеши фокуса и кликабельных виджетов. Используется при динамическом обновлении дерева виджетов.

	Commit(func()) // Commit выполняет задачу в UI потоке, и автоматически перерисовывает окно.

	SetStyleFunc(fn func(Widget)) // SetStyleFunc устанавливает функцию для стилизации виджетов.
}

// FocusManager — интерфейс менеджера фокуса.
// Добавлено в TUI 3.1.0.
type FocusManager interface {
	FocusedWidget() EventHandler // FocusedWidget вовзращает виджет, на котором установлен фокус.
	NextFocus()                  // NextFocus переносит фокус вперёд.
	BeforeFocus()                // BeforeFocus переносит фокус назад.
	SetFocus(EventHandler) bool  // BeforeFocus устанавливает фокус на переданный виджет.
	ClearFocus()                 // ClearFocus сбрасывает фокус.
	Disable()                    // Disable отключает автоматическую смену фокуса.
	Enable()                     // Enable включает автоматическую смену фокуса, если была выключена.
	FocusedIndex() int           // FocusedIndex возвращает индекс текущего фокуса, или -1 если фокус не установлен.
	SetIndex(idx int)            // SetIndex устанавливает фокус на виджет с индексом idx.
}

// Container это интерфейс контейнеров.
// Добавлено в TUI 3.0.0.
type Container interface {
	Widget
	Child() []Widget // Child возвращает детей контейнера.
	Pos(int) Pos     // Pos возвращает позицию ребёнка с индексом.
}
