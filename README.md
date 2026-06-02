# go-tui: Text User Interface на Go
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui)](https://github.com/romanSPB15/tui/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/romanSPB15/go-tui.svg)](https://pkg.go.dev/github.com/romanSPB15/go-tui)

## Обзор
Лёгкая библиотека для TUI на Go. Может:
* Создавать надписи
* Создавать и настраивать кнопки
* Обрабатывать нажатия клавиатуры
* Красить текст
## Установка
```
go get -u github.com/romanSPB15/go-tui
```
## Пример использования
```go
package main

import (
	"github.com/eiannone/keyboard"
	"github.com/romanSPB15/go-tui"
)

func main() {
	w := tui.NewWindow()
	w.AddWidgets(tui.NewStaticLabel("Привет, Go!"))
	btn := tui.NewButton("Нажми ↑", keyboard.KeyArrowUp)
	btn.OnClick = w.Quit
	w.AddWidgets(btn)
	w.Run()
}
```