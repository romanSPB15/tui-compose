![demo/demo.gif](demo/demo.gif)
# TUI Compose
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest release](https://img.shields.io/github/v/release/romanSPB15/tui-compose)](https://github.com/romanSPB15/tui-compose/releases)
[![Test](https://github.com/romanSPB15/tui-compose/actions/workflows/test.yaml/badge.svg)](https://github.com/romanSPB15/tui-compose/actions/workflows/test.yaml)
[![Documentation](https://img.shields.io/badge/GitHub-Wiki-red?logo=github)](https://github.com/romanSPB15/tui-compose/wiki)
![Lightweight](https://img.shields.io/badge/Lightweight-~8500_lines-brightgreen)
[![Examples](https://img.shields.io/badge/Examples-21-blue?logo=github)](https://github.com/romanSPB15/tui-compose/tree/main/examples)

**Лёгкий путь для создания приложений в терминале на Go.**

* 🍬 21 готовый виджет — полноценный набор.
    - **Базовые:**
        * Label, Button, Hyperlink, Gauge, InputField, Check, Frame(рамка)
    - **Изображения:**
        * Canvas(16 цветов), CanvasRGB(True Color)
    - **`extra`**:
        * LineChart, BarChart, PieChart, Sparkline, Tree, Table, Tabs, Accordion, Spinner, BlinkLabel, PageIndicator, TextView.
* 👓 Diff-рендер, low-allocation — без мерцаний и аллокаций
* ✨ Автоматический фокус по Tab/Shift+Tab
* 💎 Минимальный размер: 8500 строк кода, только `x/sys` + `x/term`
* 🎁 Полная поддержка Windows — без WSL
* 🛠 Детекция data race при вызове методов окна — `-tags debug`
* 🎨 Удобная кастомизация через битовые маски Style — быстро и удобно
* 🔧 Кастомные виджеты — 3-6 методов(зависит от функцонала)
* 🛒 Контейнеры — не нужно считать координаты вручную
* 🚀 react-надстройка — пишите реактивные приложения
* 🏹 Поддержка мыши - клики

---

<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get -u github.com/romanSPB15/tui-compose/v3</pre></h3>
<p align="center">
<table>
<tbody>
</tbody>
</table>
</p>

## Быстрый старт
```go
package main

import (
	"strconv"

	"github.com/romanSPB15/tui-compose/v3"
)

func main() {
	wnd := tui.NewWindow()
	wnd.SetTitle("Моё приложение")

	label := tui.NewStaticLabel("Привет, TUI!").WithStyle(tui.FrCyan)

	btnQuit := tui.NewButton("Выход", func() {
		wnd.Quit()
	})

	v := 0

	btnAdd := tui.NewButton("+", func() {
		v++
		label.SetText(strconv.Itoa(v))
		wnd.Redraw()
	}).WithStyle(tui.BgRed)

	btnSub := tui.NewButton("-", func() {
		if v > 0 {
			v--
			label.SetText(strconv.Itoa(v))
			wnd.Redraw()
		}
	}).WithStyle(tui.BgBlue)

	box := tui.NewVBox(label, tui.NewHBox(btnAdd, btnSub), btnQuit)
	wnd.SetContent(box)

	wnd.Run()
}
```

*Подробная документация доступна в [Wiki](https://github.com/romanSPB15/tui-compose/wiki).*<br>
*Другие примеры доступны в [examples](https://github.com/romanSPB15/tui-compose/tree/main/examples).*


## Сравнение с `gizak/termui`
| Критерий                     | tui-compose(v3.3.0)                       | termui(v3.1.0)                          |
|------------------------------|-------------------------------------------|-----------------------------------------|
| Размер                       | ✅ ~8500 строк                            | ⚠ 4325 строк + 5393(termbox-go)        |
| Примеры                      | 21                                        | 16                                      |
| Возраст                      | ~3 месяца                                 | 12 лет                                  |
| Набор виджетов               | ✅ Огромный(21) — от кнопок до графиков   | ⚠ Средний(12)                          |
| Рендеринг                    | ✅ low-allocation, diff-рендер            | ❌ Полная перерисовка                  |
| Кастомизация                 | ✅ Полная(стили, выбор символов)          | ⚠ Частичная                            |
| Зависимости                  | ✅ 2(только `x/sys` + `x/term`)           | ❌ 5+ — `termbox-go`, `x/sys` и другие |
| Тесты                        | ✅ Покрытие низкоуровневых пакетов        | ❌ Нет                                 |
| Интерактивность «из коробки» | ✅ Есть                                   | ❌ Нет(реализовывать вручную)          |
| Композиция                   | ✅ VBox + HBox                            | ⚠ Grid + абсолютные                    |

## Лицензия
[**MIT**](https://github.com/romanSPB15/tui-compose/blob/main/LICENSE)

<div align="center">
  <h3>⭐ Нравится проект? Поддержи его звездой!</h3>
  <p>
    <a href="https://github.com/romanSPB15/tui-compose">
      <img src="https://img.shields.io/github/stars/romanSPB15/tui-compose?style=for-the-badge" alt="GitHub stars">
    </a>
  </p>
</div>
