package react

import (
	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

// App — реактивная обёртка над Window.
// T — тип модели, хранящей всё состояние приложения.
type App[T any] struct {
	wnd    tui.Window
	model  T
	render func(*App[T], T) tui.Widget
}

// New создаёт новое реактивное приложение.
func New[T any](initial T, render func(*App[T], T) tui.Widget) *App[T] {
	a := &App[T]{
		wnd:    tui.NewWindow(),
		model:  initial,
		render: render,
	}
	a.wnd.SetContent(render(a, initial))
	return a
}

// Get возвращает копию текущей модели.
func (a *App[T]) Get() T {
	return a.model
}

// Mutate безопасно изменяет модель и автоматически перерисовывает UI.
func (a *App[T]) Mutate(f func(*T)) {
	a.wnd.Commit(func() {
		idx := a.wnd.Focus().FocusedIndex()
		f(&a.model)
		a.wnd.SetContent(a.render(a, a.model))
		if idx >= 0 {
			a.wnd.Focus().SetIndex(idx)
		}
	})
}

// Commit выполняет функцию в UI-потоке и перерисовывает через react-render.
// В отличие от Window.Commit, здесь перерисовка идёт через render-функцию.
func (a *App[T]) Commit(f func()) {
	a.Mutate(func(_ *T) {
		f()
	})
}

// Run запускает приложение.
func (a *App[T]) Run() {
	a.wnd.Run()
}

// Redraw перерисовывает окно.
func (a *App[T]) Redraw() {
	a.wnd.Redraw()
}

// IsRunned возвращает true, если приложение запущено.
func (a *App[T]) IsRunned() bool {
	return a.wnd.IsRunned()
}

// Quit завершает приложение.
func (a *App[T]) Quit() {
	a.wnd.Quit()
}

// OnQuit возвращает канал для ожидания завершения.
func (a *App[T]) OnQuit() <-chan struct{} {
	return a.wnd.OnQuit()
}

// RegisterKeyHandler регистрирует обработчик клавиатуры.
func (a *App[T]) RegisterKeyHandler(h tui.KeyboardEventHandler) {
	a.wnd.RegisterKeyHandler(h)
}

// RegisterClickHandler регистрирует глобальный обработчик мыши.
func (a *App[T]) RegisterClickHandler(h func(ev *input.MouseEvent)) {
	a.wnd.RegisterClickHandler(h)
}

// LogInfo логирует сообщение в debug-файл.
func (a *App[T]) LogInfo(message string, args ...any) {
	a.wnd.LogInfo(message, args...)
}

// LogFatal логирует сообщение и завершает приложение.
func (a *App[T]) LogFatal(message string, args ...any) {
	a.wnd.LogFatal(message, args...)
}

// Do отправляет задачу в UI-поток.
func (a *App[T]) Do(f func()) {
	a.wnd.Do(f)
}

// DoAndWait отправляет задачу в UI-поток и блокируется до завершения.
func (a *App[T]) DoAndWait(f func()) {
	a.wnd.DoAndWait(f)
}

// Width возвращает ширину окна.
func (a *App[T]) Width() int {
	return a.wnd.Width()
}

// Height возвращает высоту окна.
func (a *App[T]) Height() int {
	return a.wnd.Height()
}

// SetTitle устанавливает заголовок окна.
func (a *App[T]) SetTitle(title string) {
	a.wnd.SetTitle(title)
}

// CopyToClipboard копирует текст в буфер обмена.
func (a *App[T]) CopyToClipboard(text string) {
	a.wnd.CopyToClipboard(text)
}

// Focus возвращает менеджер фокуса.
func (a *App[T]) Focus() tui.FocusManager {
	return a.wnd.Focus()
}

// SetInitCell устанавливает ячейку по умолчанию.
func (a *App[T]) SetInitCell(c cell.Cell) {
	a.wnd.SetInitCell(c)
}

// SetBackground устанавливает фон пустых позиций.
func (a *App[T]) SetBackground(s tui.Style) {
	a.wnd.SetBackground(s)
}

// Index обновляет кеши фокуса и кликабельных виджетов.
func (a *App[T]) Index() {
	a.wnd.Index()
}

// SetStyleFunc устанавливает глобальную функцию стилизации.
func (a *App[T]) SetStyleFunc(fn func(tui.Widget)) {
	a.wnd.SetStyleFunc(fn)
}
