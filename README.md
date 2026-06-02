# go-tui: Text User Interface на Go
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui)](https://github.com/romanSPB15/tui/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/romanSPB15/go-tui.svg)](https://pkg.go.dev/github.com/romanSPB15/go-tui)

## Обзор
Фреймворк для разработки TUI-интерфейсов на Go.
Преимущества:
* Простой API
* Небольшой размер ≈ 1МБ
* Поддержка кастомной графики через Canvas
* Поддержка конкурентности на уровне фреймворка
* Поддержка Windows & Linux
## Установка
```
go get -u github.com/romanSPB15/go-tui
```
## Пример использования
```go
package main

import (
	"github.com/romanSPB15/go-tui"
)

func main() {
	w := tui.NewWindow()                             // Создаём приложение
	w.AddWidgets(tui.NewStaticLabel("Привет, TUI!")) // Добавляем надпись
	btn := tui.NewButton("Нажми на меня!")           // Создаём кнопку
	btn.OnClicked = w.Quit                           // Назначаем обработчик
	w.AddWidgets(btn)                                // Добавляем кнопку
	w.Run()                                          // Запускаем приложение
}
```
## Обновление v2.0.0!
[Release Notes](https://github.com/romanSPB15/go-tui/blob/main/Release-Notes.md)
|
[Документация](https://pkg.go.dev/github.com/romanSPB15/go-tui@v2.0.0)
