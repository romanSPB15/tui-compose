package react

import (
	"github.com/romanSPB15/tui-compose/v3"
	"github.com/romanSPB15/tui-compose/v3/input"
)

// App — реактивная обёртка над Window.
// T — тип модели, хранящей всё состояние приложения.
type App[T any] struct {
	wnd    tui.Window
	model  T
	render func(*App[T], T) tui.Widget
}

// New создаёт новое реактивное приложение.
// render — функция, которая по модели возвращает Widget.
// Она будет вызываться при каждом изменении модели.
func New[T any](initial T, render func(*App[T], T) tui.Widget) *App[T] {
	a := &App[T]{
		wnd:    tui.NewWindow(),
		model:  initial,
		render: render,
	}
	a.wnd.SetContent(render(a, initial))
	return a
}

// Mutate безопасно изменяет модель и автоматически перерисовывает UI.
// f — функция, которая получает указатель на текущую модель и изменяет её.
func (a *App[T]) Mutate(f func(*T)) {
	a.wnd.DoAndWait(func() {
		f(&a.model)
		a.wnd.SetContent(a.render(a, a.model))
		a.wnd.Redraw()
	})
}

// Run запускает приложение.
func (a *App[T]) Run() {
	a.wnd.Run()
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
func (a *App[T]) RegisterKeyHandler(h func(*input.KeyboardEvent)) {
	a.wnd.RegisterKeyHandler(h)
}

// SetTitle устанавливает заголовок окна.
func (a *App[T]) SetTitle(title string) {
	a.wnd.SetTitle(title)
}

// Width возвращает ширину окна.
func (a *App[T]) Width() int {
	return a.wnd.Width()
}

// Height возвращает высоту окна.
func (a *App[T]) Height() int {
	return a.wnd.Height()
}

// SetOverlay устанавливает оверлей.
func (a *App[T]) SetOverlay(w tui.Widget) {
	a.wnd.SetOverlay(w)
}

// ShowOverlay показывает оверлей.
func (a *App[T]) ShowOverlay() {
	a.wnd.ShowOverlay()
}

// HideOverlay скрывает оверлей.
func (a *App[T]) HideOverlay() {
	a.wnd.HideOverlay()
}

func (a *App[T]) Focus() tui.FocusManager {
	return a.wnd.Focus()
}
