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
* Инкрементальный рендеринг — только изменённые ячееки.
* Пакет `react` для реактивности.
* Пакет `keyboard` для использования парсера клавиатуры отдельно.
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