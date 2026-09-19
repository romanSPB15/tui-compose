package tui

import (
	"github.com/romanSPB15/tui-compose/v4/cell"
)

// Widget — это интерфейс для TUI-виджетов.
type Widget interface {
	Render(buf [][]cell.Cell) // Render рисует в переданный буфер, начиная с координаты (0, 0). Буфер имеет размер Width × Height.

	Width() int  // Ширина в символах
	Height() int // Высота в символах
}

// EventHandler — интерфейс виджетов-обработчиков событий.
type EventHandler interface {
	Widget
	Send(Event)
}

// Window — это объект приложения. Создаётся через NewWindow,
// запускается через Run, останавливается через Quit.
// Все мутации состояния UI должны идти через Do/DoAndWait/Commit,
// потому что дерево виджетов принадлежит одной горутине.
type Window interface {
	SetContent(Widget) // SetContent устанавливает содержимое окна.

	Redraw() // Redraw перерисовывает окно.

	Run()           // Run запускает TUI-приложение. Если пользователь закроет окно(Ctrl+C), то будет произведён graceful shutdown и выход из метода.
	IsRunned() bool // IsRunned возвращает, запущено ли приложение.

	Quit()                   // Quit выходит из приложения.
	OnQuit() <-chan struct{} // OnQuit возвращает канал сигнализации о выходе.

	RegisterKeyHandler(KeyboardEventHandler) // RegisterKeyHandler регистрирует обработчик событий клавиатуры.
	RegisterClickHandler(MouseEventHandler)  // RegisterClickHandler регистрирует обработчик событий мыши.
	RegisterResizeHandler(ResizeHandler)     // RegisterResizeHandler регистрирует обработчик изменения размера окна.

	LogInfo(message string, args ...any)  // LogInfo логирует сообщение подобно fmt.Printf в файл, если включен debug режим.
	LogFatal(message string, args ...any) // LogFatal логирует сообщение подобно fmt.Printf в файл, если включен debug режим, и завершает приложение.

	Do(func())        // Do отправляет задачу в UI поток.
	DoAndWait(func()) // DoAndWait отправляет задачу в UI поток и блокируется до завершения.

	Width() int  // Width возвращает ширину окна в символах.
	Height() int // Height возвращает высоту окна в символах.

	SetTitle(string)        // SetTitle устанавливает заголовок окна терминала.
	CopyToClipboard(string) // CopyToClipboard копирует текст в буфер обмена.

	Focus() FocusManager // Focus возвращает FocusManager окна.

	SetInitCell(cell.Cell) // SetInitCell устанавливает ячейку по умолчанию для всех пустых позиций окна.
	SetBackground(Style)   // SetBackground устанавливает стиль пустых позиций окна.

	// Index обновляет кеши фокуса и кликабельных виджетов. Используется при динамическом обновлении дерева виджетов.
	// Также вызывает SetStyleFunc для всех виджетов в дереве.
	Index()

	// Commit выполняет задачу в UI потоке, и автоматически перерисовывает окно.
	//
	// go func() {
	//    w.Commit(func() { lbl.SetText("New text") })
	// }()
	Commit(func())

	// SetStyleFunc устанавливает функцию для стилизации виджетов.
	// Функция применяется при каждом Index (SetContent, SetOverlay) и
	// может менять стили, отступы и другие параметры виджетов через
	// type switch. Это позволяет задавать глобальные темы без
	// изменения кода виджетов.
	//
	//	wnd.SetStyleFunc(func(w tui.Widget) {
	//		switch x := w.(type) {
	//		case *tui.Button:
	//			x.WithStyle(tui.BgBlue).WithPaddings(2, 1)
	//		case *tui.InputField:
	//			x.WithStyle(tui.BgBlack).WithFocusedStyle(tui.BgBlue | tui.FrWhite)
	//		case *tui.Label:
	//			x.WithStyle(tui.FrCyan)
	//		}
	//	})
	//
	// Добавлено в TUI v3.5.0.
	SetStyleFunc(fn func(Widget))

	// SetAltScreenEnable включает/выключает alt-screen.
	// По умолчанию включен. Он сохраняет историю терминала.
	SetAltScreenEnable(v bool)
}

// FocusManager — интерфейс менеджера фокуса.
// Добавлено в TUI v3.1.0.
type FocusManager interface {
	FocusedWidget() EventHandler // FocusedWidget возвращает виджет, на котором установлен фокус.
	NextFocus()                  // NextFocus переносит фокус вперёд.
	BeforeFocus()                // BeforeFocus переносит фокус на предыдущий виджет.
	SetFocus(EventHandler) bool  // SetFocus устанавливает фокус на переданный виджет.
	ClearFocus()                 // ClearFocus сбрасывает фокус.
	Disable()                    // Disable отключает автоматическую смену фокуса.
	Enable()                     // Enable включает автоматическую смену фокуса обратно.
	FocusedIndex() int           // FocusedIndex возвращает индекс текущего фокуса, или -1 если фокус не установлен.
	SetIndex(idx int)            // SetIndex устанавливает фокус на виджет с индексом idx.
}

// Container это интерфейс контейнеров.
// Добавлено в TUI v3.0.0.
type Container interface {
	Widget
	Child() []Widget // Child возвращает детей контейнера.
	Pos(int) Pos     // Pos возвращает позицию ребёнка с индексом.
}
