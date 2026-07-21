![demo/demo.gif](demo/demo.gif)

# TUI Compose
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui-compose)](https://github.com/romanSPB15/tui-compose/releases)
[![Test](https://github.com/romanSPB15/tui-compose/actions/workflows/test.yaml/badge.svg)](https://github.com/romanSPB15/tui-compose/actions/workflows/test.yaml)
[![Documentation](https://img.shields.io/badge/GitHub-Wiki-red?logo=github)](https://github.com/romanSPB15/tui-compose/wiki)
![Lightweight](https://img.shields.io/badge/Lightweight-<4000_lines-brightgreen)
![Paradigm](https://img.shields.io/badge/Paradigm-Reactive%20%2B%20Imperative-blue?logo=go)
[![Examples](https://img.shields.io/badge/2-Examples-red)](https://github.com/romanSPB15/tui-compose/tree/main/examples)

Легковесный фреймворк для удобной разработки TUI-интерфейсов на Go с поддержкой мыши, готовыми виджетами и автоматической системой фокуса без зависимостей. Идеально подходит для дашбордов, консольных утилит и интерактивных CLI-приложений.

## Быстрый старт
```go
package main

import "github.com/romanSPB15/tui-compose/v3"

func main() {
    wnd := tui.NewWindow()
    wnd.SetTitle("Моё приложение")

    label := tui.NewStaticLabel("Привет, TUI!").ColorizeForeground(tui.Cyan)

    btn := tui.NewButton("Выход", func() {
        wnd.Quit()
    })

    box := tui.NewVBox(label, btn)
    wnd.SetContent(box)

    wnd.Run()
}
```

## Готовые виджеты

| Виджет          | Описание                                                    |
|-----------------|-------------------------------------------------------------|
| `Label`         | Текстовая метка                                             |
| `Button`        | Кнопка с обработчиком нажатия                               |
| `Check`         | Чекбокс                                                     |
| `InputField`    | Текстовое поле ввода                                        |
| `ColorProgress` | Прогресс бар из цветных блоков                              |
| `TextProgress`  | Прогресс бар из любых символов                              |
| `Canvas`        | Холст с 16-цветными пикселями                               |
| `CanvasRGB`     | Холст с RGB-пикселями                                       |



## Установка
```
go get -u github.com/romanSPB15/tui-compose/v3
```


#
