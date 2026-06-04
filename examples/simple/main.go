package main

import (
	"github.com/romanSPB15/tui-compose/v2"
)

func main() {
	w := tui.NewWindow()                             // Создаём приложение
	w.AddWidgets(tui.NewStaticLabel("Привет, TUI!")) // Добавляем надпись
	btn := tui.NewButton("Нажми на меня!")           // Создаём кнопку
	btn.OnClicked = w.Quit                           // Назначаем обработчик
	w.AddWidgets(btn)                                // Добавляем кнопку
	w.AddWidgets(tui.NewButton("Кнопка 2"))          // Добавляем ещё кнопку
	w.Run()                                          // Запускаем приложение
}
