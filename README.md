# TUI Compose
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui-compose)](https://github.com/romanSPB15/tui-compose/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/romanSPB15/tui-compose/v2.svg)](https://pkg.go.dev/github.com/romanSPB15/tui-compose/v2)

![Скриншот](demo.gif)

## Обзор
Фреймворк для разработки TUI-интерфейсов на Go.
Преимущества:
* Простой API
* Небольшой размер исполняемого файла ≈ 3МБ
* Поддержка кастомной графики через Canvas
* Поддержка конкурентности на уровне фреймворка
* Поддержка Windows & Linux

## Установка
```
go get -u github.com/romanSPB15/tui-compose/v2@v2.0.1
```
## Пример использования
```go
package main

import (
	"github.com/romanSPB15/tui-compose/v2"
)

func main() {
	wnd := tui.NewWindow()                             // Создаём приложение
	wnd.AddWidgets(tui.NewStaticLabel("Привет, TUI!")) // Добавляем надпись
	btn := tui.NewButton("Нажми на меня!")           // Создаём кнопку
	btn.OnClicked = wnd.Quit                           // Назначаем обработчик
	wnd.AddWidgets(btn)                                // Добавляем кнопку
	wnd.Run()                                          // Запускаем приложение
}
```
## Обновление v2.0.0!
[Release Notes](https://github.com/romanSPB15/tui-compose/blob/main/Release-Notes.md)
|
[Документация](https://pkg.go.dev/github.com/romanSPB15/tui-compose/v2)
|
[Примеры](https://github.com/romanSPB15/tui-compose/tree/main/examples)
