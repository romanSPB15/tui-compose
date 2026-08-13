# Changelog

## [3.0.0] — 2026-06-29

### Added
- Контейнеры `HBox` и `VBox` с автоматическим расчётом позиций (интерфейс `Container`).
- Виджет `InputField` — однострочное поле ввода с курсором.
- Поддержка мыши: SGR (1006), события кликов.
- Виджеты `Canvas` (16 цветов) и `CanvasRGB` (True Color) с двумя режимами отображения (`PixelTwoSymbol`, `PixelOneSymbol`).
- Методы `Width()` и `Height()` для `Canvas` и `CanvasRGB`.
- Методы `SetContent`, `SetOverlay`, `ShowOverlay`, `HideOverlay` для управления содержимым окна.
- Структура `Page` с методом `Open()`.
- Методы `SetTitle`, `CopyToClipboard`.
- Интерфейсы `TextInput`, `Container`, `ClickableAt`.
- Собственный парсер ANSI-последовательностей.
- `RegisterClickHandler` — глобальные обработчики мыши.
- Методы `AddWidgets`, `Clear` — заменены на `SetContent`.


### Changed
- Удалена зависимость от `github.com/eiannone/keyboard` — теперь используется собственный парсер ввода.
- Система фокуса переключения по Tab/Shift+Tab.
- Перерисовка оптимизирована (перерисовываются только изменённые строки).
- `RegisterKeyHandler` — регистрация обработчиков клавиатуры теперь с `KeyboardEvent`.

### Removed
- Зависимость от `github.com/eiannone/keyboard`, `github.com/charmbracelet/x/term` и `github.com/acarl005/stripansi`.
- `RedrawWidget(index)` из-за новой архитектуры.


### Fixed
- Исправлены ошибки Canvas\[RGB\].

### Dependencies
- Единственные внешние зависимости — `golang.org/x/term` и `golang.org/x/sys`.

## [3.0.1] — 2026-07-01

### Added
* Добавлена обработка паник в задачах UI
* Добавлено закрытие канала задач

### Fixed
* Исправлены потенциальный deadlock очереди задач

## [3.0.2] — 2026-07-02

### Added
* Добавлены обработчики OnChanged и OnEnter в TextField

## [3.0.3] — 2026-07-03

### Fixed
* Убрано мерцание при сжатии окна
* Добавлен вызов OnChanged при стирании текста в TextField

## [3.1.0] — 2026-07-13

### Added
* Стили через `Style` и `WithStyle`
* Инкрементальный рендеринг — только изменённые ячейки.
* Пакет `react` для реактивности.
* Пакет `input` для использования парсера клавиатуры отдельно.
* Пакет `term` для работы с терминалом.

### Fixed
* Баг отрисовки виджетов ANSI.


## [3.1.1] - 2026-07-13

### Added
* Выход из приложения при `Ctrl+C`

### Fixed
* Баг получения размеров окна


## [3.1.2] - 2026-07-13

### Added
* Добавлено отображение курсора при пустом тексте в `InputField`

### Removed
* Убран белый фон при фокусе по умолчанию в `InputField`

### Fixed
* Исправлен `DrawAndRender` у `Canvas` и `CanvasRGB`

## [3.1.3] - 2026-07-15

### Added
* Добавлены тесты для пакетов `cell` и `ansi`, а также для `tui.Style`

### Fixed
* Исправлены баги инкрементального рендерера.

## [3.2.0] - 2026-08-06

### Added
* `Frame` - рамка.
* low-allocation рендер.
* Метод `Commit` у `Window`, который вызывает и `Do`, и `Redraw`.
* Поддержка ввода всех печатных ASCII-символов.
* Проверка вызывающей горутины в `Redraw` и подобных методах `Window` с помощью `gorid`.
* Восстановление состояния терминала в `LogFatal`.
* Вывод стека вызовов в `LogFatal`.
* Метод `Enable` у `FocusManager`.
* Кеширование строки в `Canvas` и `CanvasRGB`.

### Fixed
* Исправлены баги рендера.
* Гонки данных в виджетах.

### Changed
* Рефакторинг комментариев в [`types.go`](https://github.com/romanSPB15/tui-compose/blob/main/types.go)

### Dependencies
- Добавлена [`github.com/romanSPB15/x/gorid`](https://github.com/romanSPB15/x/tree/master/gorid).

## [3.3.0] - 2026-08-13

### Added
* **12 новых виджетов**:
  - `Accordion` — раскрывающийся список.
  - `BarChart` — столбчатая диаграмма.
  - `BlinkLabel` — мигающая метка.
  - `PageIndicator` — индикатор страниц.
  - `PieChart` — круговая диаграмма.
  - `Sparkline` — мини-график.
  - `Spinner` — индикатор загрузки.
  - `Table` — таблица с поддержкой заголовков и стилей.
  - `Tabs` — вкладки (с расположением сверху/снизу).
  - `TextView` — многострочный текст с тегами.
  - `Tree` — древовидный список.
  - `LineChart` — полноценный линейный график.
* **Пакет `builder`** — zero-alloc билдер строк c форматированием (в 2.3 раза быстрее `fmt.Sprintf`).
* **Производительность**:
  - `Redraw` теперь использует `builder.Builder` и `sync.Pool` — снижены аллокации.

### Fixed
* Исправлена логика сброса цветов в пакете `cell`.
* Исправлены баги парсинга ввода.

### Changed
* **Пакет `gorid` встроен в `tui`** .
* Пакет `cell`:
  - Покрытие тестами поднято с ~50% до **95%**.
  - Добавлены тесты для `ParseMultiline` и `ToString`.
* Пакет `input`:
  - Покрытие тестами поднято с 0% до **60.5%**(парсинг).
* `tui.Style`:
  - Покрытие тестами поднято с ~50% до **100%**.
  - Добавлен тест для `ConvertToCellStyle`.
* Размер кода(с тестами): ~4500 -> 7000 строк за счёт новых виджетов.

### Removed
* Зависимость от `github.com/romanSPB15/x/gorid` (встроена в `tui`).

### Dependencies
- Удалена [`github.com/romanSPB15/x/gorid`](https://github.com/romanSPB15/x/tree/master/gorid).