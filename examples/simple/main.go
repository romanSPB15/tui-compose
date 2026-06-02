package main

import (
	"github.com/romanSPB15/go-tui"
)

func main() {
	w := tui.NewWindow()                             // Создаём приложение
	w.AddWidgets(tui.NewStaticLabel("Привет, TUI!")) // Добавляем надпись
	btn := tui.NewButton("Нажми ↑")                  // Создаём кнопку
	btn.OnClicked = w.Quit                           // Назначаем обработчик
	w.AddWidgets(btn)                                // Добавляем кнопку
	w.Run()                                          // Запускаем приложение
}
