package tea

import (
	"strings"
	"time"

	"github.com/romanSPB15/tui-compose/v4"
	"github.com/romanSPB15/tui-compose/v4/cell"
	"github.com/romanSPB15/tui-compose/v4/input"
)

// Msg — сообщение, которое обрабатывается моделью.
// Может быть любого типа: событие клавиатуры, мыши, таймер, пользовательское.
// Добавлено в TUI v3.4.0.
type Msg any

// Cmd — команда, которая выполняется асинхронно и возвращает Msg.
// Возвращённый Msg будет доставлен в Update.
// Добавлено в TUI v3.4.0.
type Cmd func() Msg

// quitMsg — служебное сообщение о завершении программы.
type quitMsg struct{}

// batchMsg — служебное сообщение для одновременного запуска нескольких команд.
type batchMsg struct{ cmds []Cmd }

// Quit возвращает Msg, который завершает программу.
// Добавлено в TUI v3.4.0.
func Quit() Msg { return quitMsg{} }

// QuitCmd возвращает Cmd, который при выполнении завершает программу.
// Используется в Update для возврата команды выхода.
// Добавлено в TUI v3.4.0.
func QuitCmd() Cmd {
	return func() Msg { return quitMsg{} }
}

// Batch объединяет несколько команд в одну.
// Все команды выполняются параллельно, их Msg доставляются в Update.
// Добавлено в TUI v3.4.0.
func Batch(cmds ...Cmd) Cmd {
	return func() Msg {
		return batchMsg{cmds: cmds}
	}
}

// After создаёт команду, которая выполнится через указанный интервал.
// fn получает время срабатывания и возвращает Msg.
// Добавлено в TUI v3.4.0.
func After(d time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		<-time.After(d)
		return fn(time.Now())
	}
}

// Model — интерфейс модели приложения в стиле Elm.
// Init возвращает начальную команду (может быть nil).
// Update обрабатывает сообщение и возвращает новую модель и следующую команду.
// View строит виджет по текущему состоянию модели.
// Добавлено в TUI v3.4.0.
type Model interface {
	Init() Cmd
	Update(msg Msg) (Model, Cmd)
	View() tui.Widget
}

// Program — объект приложения, связывающий модель с окном и циклом событий.
// Добавлено в TUI v3.4.0.
type Program struct {
	window tui.Window
	model  Model
	msgCh  chan Msg
}

// currentProgram — последняя созданная программа.
// Используется пакетными функциями Send, Focus и ClearFocusCmd.
var currentProgram *Program

// NewProgram создаёт программу с указанной начальной моделью.
// Делает программу текущей для пакетных функций.
// Добавлено в TUI v3.4.0.
func NewProgram(initialModel Model) *Program {
	pr := &Program{
		model: initialModel,
		msgCh: make(chan Msg, 64),
	}
	currentProgram = pr
	return pr
}

// Run создаёт окно, регистрирует обработчики ввода, запускает цикл событий
// и блокируется до завершения программы.
// Добавлено в TUI v3.4.0.
func (p *Program) Run() {
	p.window = tui.NewWindow()
	p.window.SetContent(p.model.View())

	p.window.RegisterClickHandler(func(ev *input.MouseEvent) {
		p.Send(*ev)
	})

	p.window.RegisterKeyHandler(func(ev *input.KeyboardEvent) {
		p.Send(*ev)
	})

	go p.eventLoop()
	p.window.Run()
}

// Send отправляет сообщение в цикл событий.
// Не блокируется, если буфер канала не заполнен.
// Добавлено в TUI v3.4.0.
func (p *Program) Send(msg Msg) {
	p.msgCh <- msg
}

// eventLoop — основной цикл обработки сообщений.
// Выполняет Init, обрабатывает batch-команды, вызывает Update,
// запускает асинхронные команды и обновляет View после каждого сообщения.
func (p *Program) eventLoop() {
	if cmd := p.model.Init(); cmd != nil {
		p.msgCh <- cmd()
	}

	for msg := range p.msgCh {
		if batch, ok := msg.(batchMsg); ok {
			for _, cmd := range batch.cmds {
				go func(c Cmd) {
					p.msgCh <- c()
				}(cmd)
			}
			continue
		}

		newModel, cmd := p.model.Update(msg)
		p.model = newModel

		if cmd != nil {
			go func(c Cmd) {
				p.msgCh <- c()
			}(cmd)
		}

		select {
		case <-p.window.OnQuit():
			close(p.msgCh)
			return
		default:
			p.window.Commit(func() {
				idx := p.window.Focus().FocusedIndex()

				// SetContent сбрасывает фокус — сохраняем и восстанавливаем его
				p.window.SetContent(p.model.View())

				p.window.Focus().SetIndex(idx)
			})

			if _, ok := msg.(quitMsg); ok {
				p.window.Quit()
				close(p.msgCh)
				return
			}
		}
	}
}

// stringView — простой виджет для отображения многострочного текста.
// Используется в примерах, чтобы быстро отрендерить строку без создания Label.
type stringView struct {
	wnd  tui.Window
	text string
}

func (sv *stringView) Render(buf [][]cell.Cell) {
	p := cell.ParseMultiline(sv.text)
	copy(buf, p)
}

func (sv *stringView) Width() int {
	max := 0
	for _, line := range strings.Split(sv.text, "\n") {
		if l := len([]rune(line)); l > max {
			max = l
		}
	}
	return max
}

func (sv *stringView) Height() int {
	return strings.Count(sv.text, "\n") + 1
}

// NewStringView создаёт виджет для отображения многострочной строки.
// Строка парсится как ANSI-текст (поддерживаются escape-последовательности).
// Добавлено в TUI v3.4.0.
func NewStringView(s string) tui.Widget {
	return &stringView{
		text: s,
	}
}

// Send отправляет сообщение в текущую программу.
// Если программа не создана, ничего не делает.
// Добавлено в TUI v3.4.0.
func Send(msg Msg) {
	if currentProgram != nil {
		currentProgram.Send(msg)
	}
}

// ClearFocus сбрасывает фокус в окне программы.
// Добавлено в TUI v3.4.0.
func (p *Program) ClearFocus() {
	if p.window != nil {
		p.window.Commit(func() {
			p.window.Focus().ClearFocus()
		})
	}
}

// ClearFocusCmd возвращает команду, которая сбрасывает фокус
// в текущей программе при выполнении.
// Добавлено в TUI v3.4.0.
func ClearFocusCmd() Cmd {
	return func() Msg {
		if currentProgram != nil {
			currentProgram.ClearFocus()
		}
		return nil
	}
}

// Focus возвращает менеджер фокуса текущей программы.
// Если программа не создана, возвращает nil.
// Добавлено в TUI v3.4.0.
func Focus() tui.FocusManager {
	if currentProgram == nil {
		return nil
	}
	return currentProgram.window.Focus()
}
