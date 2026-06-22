![demo.gif](demo.gif)

# TUI Compose
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui-compose)](https://github.com/romanSPB15/tui-compose/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/romanSPB15/tui-compose/v2.svg)](https://pkg.go.dev/github.com/romanSPB15/tui-compose/v2)

Легковесный фреймворк для удобной разработки TUI-интерфейсов на Go с простым API, готовыми виджетами и встроенной поддержкой асинхронного обновления UI из любой горутины.

## Быстрый старт
```go
package main

import "github.com/romanSPB15/tui-compose/v3"

func main() {
    wnd := tui.NewWindow()
    wnd.SetTitle("Моё приложение")

    label := tui.NewStaticLabel("Привет, TUI!").ColorizeForeground(tui.Blue)

    btn := tui.NewButton("Выход", func() {
        wnd.Quit()
    })

    box := tui.NewVBox(label, btn)
    wnd.SetContent(box)

    wnd.Run()
}
```

## Готовые компоненты
--------------------------------------------------------------------------------
| `Label`         | Текстовая метка                                            |
--------------------------------------------------------------------------------
| `Button`        | Кнопка с обработчиком нажатия                              |
--------------------------------------------------------------------------------
| `Check`         | Чекбокс                                                    |
--------------------------------------------------------------------------------
| `TextEntry`     | Текстовое поле ввода                                       |
--------------------------------------------------------------------------------
| `ColorProgress` | Прогресс бар из цветных блоков                             |
--------------------------------------------------------------------------------
| `TextProgress`  | Прогресс бар из любых символов                             |
--------------------------------------------------------------------------------
| `Canvas`        | Холст с 16-цветной псевдографикой                          |
--------------------------------------------------------------------------------
| `CanvasRGB`     | Холст с RGB-псевдографикой(требуется терминал с True Color)|
--------------------------------------------------------------------------------



## Установка
```
go get -u github.com/romanSPB15/tui-compose/v3
```